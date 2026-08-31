package ui

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"emissioncalculator/internal/models"

	"fyne.io/fyne/v2"
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
	exportLabel = func(r models.CalculationRecord) (string, error) {
		printed = r.Valid &&
			math.Abs(float64(r.Emissions)-26.76) < 0.01 &&
			math.Abs(r.CO2PerKWh-0.2664) < 0.001
		return "/tmp/emissions.pdf", nil
	}
	openExportedFile = func(string) error { return nil }

	oilView := buildReferenceView(window, modeOil)
	oilView.quantityEntry.SetText("10.0")
	if oilView.result.Valid {
		t.Fatal("expected input changes not to calculate automatically")
	}
	test.Tap(oilView.calculateButton)
	if oilView.resultValue.Text != "26,76" || oilView.resultUnit.Text != "kg CO₂" {
		t.Fatalf("unexpected oil result: %s", oilView.resultValue.Text)
	}
	if oilView.headerStatus.Text != "Berechnet" || !oilView.result.Valid {
		t.Fatal("expected the successful calculation state")
	}
	test.Tap(oilView.printButton)
	if !printed {
		t.Fatal("expected the print button to export the calculated oil result")
	}
	if oilView.headerStatus.Text != "PDF erstellt" {
		t.Fatalf("unexpected PDF status: %s", oilView.headerStatus.Text)
	}

	briquettesView := buildReferenceView(window, modeBriquettes)
	briquettesView.quantityEntry.SetText("1,5")
	briquettesView.quantityEntry.OnSubmitted(briquettesView.quantityEntry.Text)
	if briquettesView.resultValue.Text != "2.827,20" || briquettesView.resultUnit.Text != "kg CO₂" {
		t.Fatalf("unexpected briquettes result after enter: %s", briquettesView.resultValue.Text)
	}
}

func TestReferenceViewShowsAndClearsValidationState(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	window := app.NewWindow("test")
	view := buildReferenceView(window, modeOil)

	view.quantityEntry.SetText("ungültig")
	test.Tap(view.calculateButton)
	if !view.quantityControl.invalid || view.headerStatus.Text != "Eingabe prüfen" {
		t.Fatal("expected invalid input to update field and header state")
	}
	if view.result.Valid || view.printButton.Disabled() == false {
		t.Fatal("expected invalid input to clear the result and disable PDF export")
	}

	view.quantityEntry.SetText("12,5")
	if view.quantityControl.invalid || view.headerStatus.Text != "Bereit" || view.result.Valid {
		t.Fatal("expected valid input to wait for an explicit calculation")
	}
	test.Tap(view.calculateButton)
	view.quantityEntry.SetText("13")
	if view.headerStatus.Text != "Eingabe geändert – neu berechnen" || !view.printButton.Disabled() || !view.saveButton.Disabled() {
		t.Fatal("expected changed input to mark the visible result as stale and disable export")
	}
}

func TestCustomControlsSupportKeyboardActivation(t *testing.T) {
	activated := 0
	action := newActionButton("Berechnen", func() { activated++ })
	action.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	navigation := newFuelNavigationButton("oil", "Heizöl", false, func() { activated++ })
	navigation.TypedKey(&fyne.KeyEvent{Name: fyne.KeySpace})
	print := newCircleIconButton(printIconResource(), func() { activated++ })
	print.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
	print.Disable()
	print.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
	if activated != 3 {
		t.Fatalf("expected three keyboard activations, got %d", activated)
	}
}

func TestFormatQuantityDisplay(t *testing.T) {
	tests := map[float64]string{
		12500:   "12.500",
		12.5:    "12,5",
		1234.75: "1.234,75",
	}
	for value, expected := range tests {
		if actual := formatQuantityDisplay(value); actual != expected {
			t.Fatalf("formatQuantityDisplay(%v) = %q, expected %q", value, actual, expected)
		}
	}
}

func TestResultTextSizeAdaptsToLongValues(t *testing.T) {
	if resultTextSize("33.036,05") != 72 {
		t.Fatal("expected regular result size")
	}
	if resultTextSize("1.234.567,89") >= 72 {
		t.Fatal("expected long result to use a smaller size")
	}
}

func TestGermanResultFormattingMatchesReference(t *testing.T) {
	if actual := formatGermanNumberString("33036,05", 2); actual != "33.036,05" {
		t.Fatalf("unexpected emission format: %s", actual)
	}
	if actual := formatMeasurement("124009,295 kWh", 2); actual != "124.009,30 kWh" {
		t.Fatalf("unexpected energy format: %s", actual)
	}
}

func TestReferenceViewFitsTargetWindow(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	window := app.NewWindow("test")
	view := buildReferenceView(window, modeOil)
	if view.scroll == nil {
		t.Fatal("expected the long reference layout to be vertically scrollable")
	}
	window.SetContent(view.content)
	window.Resize(fyne.NewSize(1080, 650))

	minimum := view.content.MinSize()
	if minimum.Width > 1080 || minimum.Height > 650 {
		t.Fatalf("content minimum %v does not fit the target window", minimum)
	}
	captured := window.Canvas().Capture()
	if captured.Bounds().Dx() != 1080 || captured.Bounds().Dy() != 650 {
		t.Fatalf("unexpected rendered size: %v", captured.Bounds())
	}
}

func TestReferenceViewRendersAtTargetWindowSizes(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	window := app.NewWindow("test")
	view := buildReferenceView(window, modeOil)
	window.SetContent(view.content)

	for _, size := range []fyne.Size{
		fyne.NewSize(1080, 650),
		fyne.NewSize(1080, 980),
		fyne.NewSize(1440, 980),
	} {
		window.Resize(size)
		captured := window.Canvas().Capture()
		if captured.Bounds().Dx() != int(size.Width) || captured.Bounds().Dy() != int(size.Height) {
			t.Fatalf("unexpected rendered size for %v: %v", size, captured.Bounds())
		}
	}
}

func TestResultUsesCompleteUnitsAndCalculationBasis(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	view := buildReferenceView(app.NewWindow("test"), modeOil)
	view.quantityEntry.SetText("10")
	test.Tap(view.calculateButton)
	if !strings.Contains(view.co2Value.Text, "kg CO₂/kWh") {
		t.Fatalf("expected complete intensity unit, got %q", view.co2Value.Text)
	}
	if !strings.Contains(view.costValue.Text, "€ brutto") {
		t.Fatalf("expected gross cost unit, got %q", view.costValue.Text)
	}
	joinedBasis := strings.Join(view.resultBasis.values(), "|")
	for _, expected := range []string{"kg CO₂/MJ", "€/t", "2022", strconv.Itoa(view.result.CalculationYear)} {
		if !strings.Contains(joinedBasis, expected) {
			t.Fatalf("missing calculation basis %q in %q", expected, joinedBasis)
		}
	}
}
