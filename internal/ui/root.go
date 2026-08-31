package ui

import (
	"emissioncalculator/internal/calculation"
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
	modeBriquettes = "briketts"
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
	resultBasis           *canvas.Text
	configProvider        func() calculation.Config
	state                 resultState
}

func NewRootWindow(app fyne.App) fyne.Window {
	app.Settings().SetTheme(emissionTheme{Theme: theme.LightTheme()})
	window := app.NewWindow("Emissionsrechner")
	window.Resize(fyne.NewSize(1080, 980))
	settings := newSettingsStore(app.Preferences())
	oilView := buildReferenceViewWithConfig(window, modeOil, settings.Config)
	briquettesView := buildReferenceViewWithConfig(window, modeBriquettes, settings.Config)
	briquettesView.content.Hide()
	navigation := buildSharedNavigation(window, settings, oilView, briquettesView)
	views := container.NewStack(oilView.content, briquettesView.content)
	window.SetContent(container.NewBorder(navigation, nil, nil, nil, views))
	return window
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
