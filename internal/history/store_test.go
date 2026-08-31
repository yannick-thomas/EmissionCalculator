package history

import (
	"emissioncalculator/internal/models"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStorePersistsProjectsAndCalculations(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	store := NewStore(path)
	project, err := store.EnsureDefaultProject()
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.CreateProject("Kunde Nord")
	if err != nil {
		t.Fatal(err)
	}
	record := models.CalculationRecord{SchemaVersion: 1, Valid: true, FuelType: "oil", Quantity: 100, Unit: "L"}
	if _, err := store.SaveCalculation(second.ID, record); err != nil {
		t.Fatal(err)
	}

	reloaded := NewStore(path)
	projects, err := reloaded.Projects()
	if err != nil || len(projects) != 2 {
		t.Fatalf("unexpected projects: %+v, %v", projects, err)
	}
	entries, err := reloaded.Entries(second.ID)
	if err != nil || len(entries) != 1 || entries[0].Record.Quantity != 100 {
		t.Fatalf("unexpected entries: %+v, %v", entries, err)
	}
	if project.ID == second.ID {
		t.Fatal("expected distinct project IDs")
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatalf("history file permissions are not private: %v, %v", info, err)
	}
}

func TestStoreRejectsUnknownSchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	if err := os.WriteFile(path, []byte(`{"schema_version":99}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewStore(path).Projects(); err == nil {
		t.Fatal("expected unknown schema to be rejected")
	}
}

func TestEntriesNewestFirst(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "history.json"))
	current := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	store.now = func() time.Time { current = current.Add(time.Minute); return current }
	project, err := store.EnsureDefaultProject()
	if err != nil {
		t.Fatal(err)
	}
	for _, quantity := range []float64{1, 2} {
		_, err := store.SaveCalculation(project.ID, models.CalculationRecord{Valid: true, Quantity: quantity})
		if err != nil {
			t.Fatal(err)
		}
	}
	entries, err := store.Entries(project.ID)
	if err != nil || entries[0].Record.Quantity != 2 {
		t.Fatalf("expected newest entry first: %+v, %v", entries, err)
	}
}
