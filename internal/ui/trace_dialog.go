package ui

import (
	"net/url"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func showTraceDialog(view *referenceView) {
	record := view.result
	steps := container.NewVBox()
	for index, step := range record.Trace.Steps {
		steps.Add(container.NewVBox(
			canvasText(strconv.Itoa(index+1)+". "+step.Title, 14, appPalette.textPrimary, true),
			canvasText(step.Expression, 12, appPalette.textSecondary, false),
			canvasText("= "+formatFloat(step.Result, traceDecimals(step.Result))+" "+step.Unit, 13, appPalette.textPrimary, true),
			verticalGap(10),
		))
	}
	metadata := container.NewVBox(
		canvasText("Faktorenpaket: "+record.Trace.FactorPackID, 12, appPalette.textSecondary, false),
		canvasText("Faktorstand: "+strconv.Itoa(record.Factor.SourceYear)+" · gültig seit "+record.Factor.ValidFrom.Format("02.01.2006"), 12, appPalette.textSecondary, false),
		canvasText("Preisstand: "+strconv.Itoa(record.Price.ReferenceYear)+" · Berechnungsjahr: "+strconv.Itoa(record.CalculationYear), 12, appPalette.textSecondary, false),
	)
	addSourceLink(metadata, "Quelle Emissionsfaktor", record.Factor.SourceURL)
	addSourceLink(metadata, "Quelle CO₂-Preis", record.Price.SourceURL)
	stepsScroll := container.NewVScroll(steps)
	stepsScroll.SetMinSize(fyne.NewSize(620, 340))
	content := container.NewVBox(
		canvasText("Jeder Wert wird aus dem gespeicherten Faktor-Snapshot neu aufgebaut.", 13, appPalette.textSecondary, false),
		verticalGap(12),
		stepsScroll,
		verticalGap(12),
		metadata,
	)
	showResponsiveDialog("So wurde gerechnet", "Schließen", content, view.window, fyne.NewSize(760, 680))
}

func addSourceLink(target *fyne.Container, label, rawURL string) {
	if rawURL == "" {
		return
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	target.Add(widget.NewHyperlink(label, parsed))
}

func traceDecimals(value float64) int {
	if value >= 1000 {
		return 2
	}
	return 4
}
