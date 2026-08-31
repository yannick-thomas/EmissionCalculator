package history

import (
	"crypto/rand"
	"emissioncalculator/internal/models"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const SchemaVersion = 1

type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Entry struct {
	ID        string                   `json:"id"`
	ProjectID string                   `json:"project_id"`
	SavedAt   time.Time                `json:"saved_at"`
	Record    models.CalculationRecord `json:"record"`
}

type Database struct {
	SchemaVersion int       `json:"schema_version"`
	Projects      []Project `json:"projects"`
	Entries       []Entry   `json:"entries"`
}

// Store persists projects and calculations in one private, atomically replaced JSON file.
type Store struct {
	path string
	mu   sync.Mutex
	now  func() time.Time
}

func NewStore(path string) *Store {
	return &Store{path: path, now: time.Now}
}

func NewDefaultStore() (*Store, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("Konfigurationsverzeichnis ermitteln: %w", err)
	}
	return NewStore(filepath.Join(configDir, "EmissionCalculator", "history.json")), nil
}

func (store *Store) EnsureDefaultProject() (Project, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	db, err := store.loadLocked()
	if err != nil {
		return Project{}, err
	}
	if len(db.Projects) > 0 {
		return db.Projects[0], nil
	}
	project, err := newProject("Standard", store.now())
	if err != nil {
		return Project{}, err
	}
	db.Projects = append(db.Projects, project)
	if err := store.saveLocked(db); err != nil {
		return Project{}, err
	}
	return project, nil
}

func (store *Store) CreateProject(name string) (Project, error) {
	name = strings.TrimSpace(name)
	if name == "" || len([]rune(name)) > 100 {
		return Project{}, fmt.Errorf("Projektname muss 1 bis 100 Zeichen lang sein")
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	db, err := store.loadLocked()
	if err != nil {
		return Project{}, err
	}
	for _, project := range db.Projects {
		if strings.EqualFold(project.Name, name) {
			return Project{}, fmt.Errorf("Projekt %q existiert bereits", name)
		}
	}
	project, err := newProject(name, store.now())
	if err != nil {
		return Project{}, err
	}
	db.Projects = append(db.Projects, project)
	if err := store.saveLocked(db); err != nil {
		return Project{}, err
	}
	return project, nil
}

func (store *Store) SaveCalculation(projectID string, record models.CalculationRecord) (Entry, error) {
	if !record.Valid {
		return Entry{}, fmt.Errorf("nur gültige Berechnungen können gespeichert werden")
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	db, err := store.loadLocked()
	if err != nil {
		return Entry{}, err
	}
	projectIndex := -1
	for index := range db.Projects {
		if db.Projects[index].ID == projectID {
			projectIndex = index
			break
		}
	}
	if projectIndex == -1 {
		return Entry{}, fmt.Errorf("Projekt %q wurde nicht gefunden", projectID)
	}
	id, err := randomID()
	if err != nil {
		return Entry{}, err
	}
	now := store.now().UTC()
	entry := Entry{ID: id, ProjectID: projectID, SavedAt: now, Record: record}
	db.Entries = append(db.Entries, entry)
	db.Projects[projectIndex].UpdatedAt = now
	if err := store.saveLocked(db); err != nil {
		return Entry{}, err
	}
	return entry, nil
}

func (store *Store) Projects() ([]Project, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	db, err := store.loadLocked()
	if err != nil {
		return nil, err
	}
	projects := append([]Project(nil), db.Projects...)
	sort.SliceStable(projects, func(i, j int) bool { return projects[i].UpdatedAt.After(projects[j].UpdatedAt) })
	return projects, nil
}

func (store *Store) Entries(projectID string) ([]Entry, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	db, err := store.loadLocked()
	if err != nil {
		return nil, err
	}
	var entries []Entry
	for _, entry := range db.Entries {
		if entry.ProjectID == projectID {
			entries = append(entries, entry)
		}
	}
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].SavedAt.After(entries[j].SavedAt) })
	return entries, nil
}

func (store *Store) loadLocked() (Database, error) {
	data, err := os.ReadFile(store.path)
	if os.IsNotExist(err) {
		return Database{SchemaVersion: SchemaVersion}, nil
	}
	if err != nil {
		return Database{}, fmt.Errorf("Historie lesen: %w", err)
	}
	var db Database
	if err := json.Unmarshal(data, &db); err != nil {
		return Database{}, fmt.Errorf("Historie lesen: %w", err)
	}
	return migrate(db)
}

func migrate(db Database) (Database, error) {
	switch db.SchemaVersion {
	case 0:
		db.SchemaVersion = SchemaVersion
	case SchemaVersion:
	default:
		return Database{}, fmt.Errorf("Historie verwendet die nicht unterstützte Schema-Version %d", db.SchemaVersion)
	}
	return db, nil
}

func (store *Store) saveLocked(db Database) error {
	db.SchemaVersion = SchemaVersion
	dir := filepath.Dir(store.path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("Historieverzeichnis erstellen: %w", err)
	}
	temporary, err := os.CreateTemp(dir, ".history-*.tmp")
	if err != nil {
		return fmt.Errorf("Historie vorbereiten: %w", err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if err := temporary.Chmod(0600); err != nil {
		temporary.Close()
		return fmt.Errorf("Historiedatei absichern: %w", err)
	}
	encoder := json.NewEncoder(temporary)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(db); err != nil {
		temporary.Close()
		return fmt.Errorf("Historie schreiben: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return fmt.Errorf("Historie synchronisieren: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("Historie schließen: %w", err)
	}
	if err := os.Rename(temporaryPath, store.path); err != nil {
		return fmt.Errorf("Historie übernehmen: %w", err)
	}
	return nil
}

func newProject(name string, now time.Time) (Project, error) {
	id, err := randomID()
	if err != nil {
		return Project{}, err
	}
	now = now.UTC()
	return Project{ID: id, Name: name, CreatedAt: now, UpdatedAt: now}, nil
}

func randomID() (string, error) {
	data := make([]byte, 12)
	if _, err := rand.Read(data); err != nil {
		return "", fmt.Errorf("ID erzeugen: %w", err)
	}
	return hex.EncodeToString(data), nil
}
