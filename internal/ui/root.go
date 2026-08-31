package ui

import (
	"emissioncalculator/internal/calculation"
	"emissioncalculator/internal/history"
	"emissioncalculator/internal/models"
	"emissioncalculator/internal/pdf"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

var exportLabel = pdf.ExportLabel
var openExportedFile = openFile

const (
	modeOil        = "oil"
	modeBriquettes = "briquettes"
	modeNaturalGas = "natural_gas"
	modeLPG        = "lpg"
)

type resultState int

const (
	resultStateReady resultState = iota
	resultStateInvalid
	resultStateCalculated
	resultStateStale
	resultStatePDFCreated
)

type referenceView struct {
	content               fyne.CanvasObject
	scroll                *container.Scroll
	window                fyne.Window
	quantityEntry         *focusEntry
	quantityControl       *quantityControl
	calculateButton       *actionButton
	printButton           *circleIconButton
	saveButton            *circleIconButton
	scenarioButton        *circleIconButton
	traceButton           *circleIconButton
	status                *canvas.Text
	headerStatus          *canvas.Text
	headerStatusDot       *canvas.Circle
	result                models.CalculationRecord
	mode                  string
	resultValue           *canvas.Text
	resultUnit            *canvas.Text
	resultValueRow        *fyne.Container
	resultHint            *canvas.Text
	resultBackground      *canvas.Image
	resultBadge           *canvas.Text
	resultBadgeBackground *canvas.Rectangle
	costValue             *canvas.Text
	energyValue           *canvas.Text
	co2Value              *canvas.Text
	resultBasis           *calculationBasis
	configProvider        func() calculation.Config
	saveHistory           func(models.CalculationRecord) error
	state                 resultState
}

func NewRootWindow(app fyne.App) fyne.Window {
	app.Settings().SetTheme(emissionTheme{Theme: theme.LightTheme()})
	window := app.NewWindow("Emissionsrechner")
	window.Resize(fyne.NewSize(1240, 900))
	window.CenterOnScreen()
	settings := newSettingsStore(app.Preferences())
	historyController := newHistoryController()
	views := make([]*referenceView, 0, len(calculation.Catalog))
	objects := make([]fyne.CanvasObject, 0, len(calculation.Catalog))
	for index, descriptor := range calculation.Catalog {
		view := buildReferenceViewWithConfig(window, string(descriptor.Fuel), settings.Config)
		view.saveHistory = historyController.Save
		if index > 0 {
			view.content.Hide()
		}
		views = append(views, view)
		objects = append(objects, view.content)
	}
	navigation := buildSharedNavigation(window, settings, historyController, views)
	window.SetContent(container.NewBorder(navigation, nil, nil, nil, container.NewStack(objects...)))
	return window
}

func newHistoryController() *historyController {
	store, err := history.NewDefaultStore()
	if err != nil {
		return &historyController{initErr: err}
	}
	controller, err := newHistoryControllerWithStore(store)
	if err != nil {
		return &historyController{store: store, initErr: err}
	}
	return controller
}

func openFile(path string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", path).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", path).Start()
	default:
		return exec.Command("xdg-open", path).Start()
	}
}
