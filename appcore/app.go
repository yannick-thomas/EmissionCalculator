package appcore

import (
	"emissioncalculator/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func Run() {
	fyneApp := app.NewWithID("emissioncalculator")
	window := ui.NewRootWindow(fyneApp)
	window.ShowAndRun()
}

func NewApp() fyne.App {
	return app.NewWithID("emissioncalculator")
}
