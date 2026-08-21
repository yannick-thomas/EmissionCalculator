package ui

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestReferenceViewCalculatesBothModesAndPrints(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	window := app.NewWindow("test")

	originalExport := exportLabel
	originalOpen := openExportedFile
	defer func() {
		exportLabel = originalExport
		openExportedFile = originalOpen
	}()

	printed := false
	exportLabel = func(emissions, cost, energy, co2 string) (string, error) {
		printed = emissions == "26,76 kg CO2" && cost != "" && energy != "" && co2 == "0,2664"
		return "/tmp/emissions.pdf", nil
	}
	openExportedFile = func(string) error { return nil }

	oilView := buildReferenceView(window, "oil")
	oilView.quantityEntry.SetText("10.0")
	test.Tap(oilView.calculateButton)
	if oilView.resultValue.Text != "26,76 kg CO2" {
		t.Fatalf("unexpected oil result: %s", oilView.resultValue.Text)
	}
	test.Tap(oilView.printButton)
	if !printed {
		t.Fatal("expected the print button to export the calculated oil result")
	}

	briquettesView := buildReferenceView(window, "briketts")
	briquettesView.quantityEntry.SetText("1,5")
	briquettesView.quantityEntry.OnSubmitted(briquettesView.quantityEntry.Text)
	if briquettesView.resultValue.Text != "2827,20 kg CO2" {
		t.Fatalf("unexpected briquettes result after enter: %s", briquettesView.resultValue.Text)
	}
}
