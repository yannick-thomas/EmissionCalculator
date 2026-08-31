package ui

import (
	"emissioncalculator/internal/history"
	"emissioncalculator/internal/models"
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type historyController struct {
	store           *history.Store
	activeProjectID string
	initErr         error
}

func newHistoryControllerWithStore(store *history.Store) (*historyController, error) {
	project, err := store.EnsureDefaultProject()
	if err != nil {
		return nil, err
	}
	return &historyController{store: store, activeProjectID: project.ID}, nil
}

func (controller *historyController) Save(record models.CalculationRecord) error {
	if controller.initErr != nil {
		return controller.initErr
	}
	_, err := controller.store.SaveCalculation(controller.activeProjectID, record)
	return err
}

func showHistoryDialog(window fyne.Window, controller *historyController) {
	if controller.initErr != nil {
		dialog.ShowError(controller.initErr, window)
		return
	}
	projects, err := controller.store.Projects()
	if err != nil {
		dialog.ShowError(err, window)
		return
	}
	projectByName := make(map[string]history.Project, len(projects))
	names := make([]string, 0, len(projects))
	selectedName := ""
	for _, project := range projects {
		projectByName[project.Name] = project
		names = append(names, project.Name)
		if project.ID == controller.activeProjectID {
			selectedName = project.Name
		}
	}
	sort.Strings(names)
	entriesBox := container.NewVBox()
	refreshEntries := func() {
		entriesBox.RemoveAll()
		entries, loadErr := controller.store.Entries(controller.activeProjectID)
		if loadErr != nil {
			entriesBox.Add(widget.NewLabel("Historie konnte nicht geladen werden: " + loadErr.Error()))
			return
		}
		if len(entries) == 0 {
			entriesBox.Add(widget.NewLabel("Noch keine Berechnungen in diesem Projekt."))
			return
		}
		for _, entry := range entries {
			record := entry.Record
			row := fmt.Sprintf("%s · %s · %s %s · %s kg CO₂ · %s €", entry.SavedAt.Local().Format("02.01.2006 15:04"), titleForMode(record.FuelType), formatQuantityDisplay(record.Quantity), record.Unit, formatFloat(float64(record.Emissions), 2), formatFloat(record.EmissionCost, 2))
			entriesBox.Add(widget.NewLabel(row))
		}
	}
	projectSelect := widget.NewSelect(names, func(name string) {
		if project, ok := projectByName[name]; ok {
			controller.activeProjectID = project.ID
			refreshEntries()
		}
	})
	projectSelect.Selected = selectedName
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Neues Projekt")
	newProjectButton := widget.NewButton("Projekt anlegen", func() {
		project, createErr := controller.store.CreateProject(nameEntry.Text)
		if createErr != nil {
			dialog.ShowError(createErr, window)
			return
		}
		projectByName[project.Name] = project
		projectSelect.Options = append(projectSelect.Options, project.Name)
		sort.Strings(projectSelect.Options)
		controller.activeProjectID = project.ID
		projectSelect.SetSelected(project.Name)
		nameEntry.SetText("")
		projectSelect.Refresh()
	})
	refreshEntries()
	entriesScroll := container.NewVScroll(entriesBox)
	entriesScroll.SetMinSize(fyne.NewSize(700, 340))
	content := container.NewVBox(
		canvasText("Aktives Projekt", 13, appPalette.textPrimary, true),
		projectSelect,
		container.NewBorder(nil, nil, nil, newProjectButton, nameEntry),
		verticalGap(12),
		canvasText("Lokale Berechnungshistorie", 13, appPalette.textPrimary, true),
		entriesScroll,
	)
	showResponsiveDialog("Projekte und Historie", "Schließen", content, window, fyne.NewSize(860, 680))
}
