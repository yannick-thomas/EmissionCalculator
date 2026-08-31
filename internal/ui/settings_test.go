package ui

import (
	"math"
	"strconv"
	"testing"

	"emissioncalculator/internal/calculation"
	"emissioncalculator/internal/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestSettingsStorePersistsAndResetsCO2Price(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	preferences := app.Preferences()
	preferences.RemoveValue(co2PricePreferenceKey)
	preferences.RemoveValue(factorYearPreferenceKey)

	store := newSettingsStore(preferences)
	if store.Config().CO2PricePerTonne != 45 {
		t.Fatalf("unexpected default price: %v", store.Config().CO2PricePerTonne)
	}

	store.SetCO2Price(72.5)
	reloaded := newSettingsStore(preferences)
	if reloaded.Config().CO2PricePerTonne != 72.5 {
		t.Fatalf("expected persisted price, got %v", reloaded.Config().CO2PricePerTonne)
	}

	reloaded.Reset()
	if reloaded.Config().CO2PricePerTonne != 45 {
		t.Fatalf("expected reset default price, got %v", reloaded.Config().CO2PricePerTonne)
	}
	if newSettingsStore(preferences).Config().CO2PricePerTonne != 45 {
		t.Fatal("expected reset to remove the persisted override")
	}
}

func TestSettingsPanelValidatesAndAppliesSharedPrice(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	app.Preferences().RemoveValue(co2PricePreferenceKey)
	store := newSettingsStore(app.Preferences())
	saved := 0
	panel := newSettingsPanel(store, func() { saved++ })

	panel.priceEntry.SetText("ungültig")
	panel.apply()
	if !panel.priceControl.invalid || saved != 0 || store.Config().CO2PricePerTonne != 45 {
		t.Fatal("expected invalid settings to be rejected")
	}

	panel.priceEntry.SetText("62,50")
	panel.apply()
	if panel.priceControl.invalid || saved != 1 || store.Config().CO2PricePerTonne != 62.5 {
		t.Fatal("expected valid settings to be saved")
	}

	panel = newSettingsPanel(store, func() { saved++ })
	panel.restoreDefault()
	panel.apply()
	if saved != 2 || store.Config().CO2PricePerTonne != 45 {
		t.Fatal("expected the default action to restore the central default")
	}
	if persisted := app.Preferences().FloatWithFallback(co2PricePreferenceKey, -1); persisted != -1 {
		t.Fatalf("expected the persisted override to be removed, got %v", persisted)
	}
}

func TestSettingsChangeRecalculatesExistingResult(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	app.Preferences().RemoveValue(co2PricePreferenceKey)
	window := app.NewWindow("test")
	store := newSettingsStore(app.Preferences())
	view := buildReferenceViewWithConfig(window, modeOil, store.Config)
	view.quantityEntry.SetText("10")
	view.calculate()
	initialCost := view.result.EmissionCost

	store.SetCO2Price(90)
	view.refreshForSettingsChange()
	if math.Abs(view.result.EmissionCost-initialCost) > 0.000001 {
		t.Fatalf("expected settings change not to alter a visible result: %v", view.result.EmissionCost)
	}
	if view.headerStatus.Text != "Eingabe geändert – neu berechnen" || !view.saveButton.Disabled() || !view.printButton.Disabled() {
		t.Fatalf("unexpected header status: %s", view.headerStatus.Text)
	}
}

func TestFormatSettingsValue(t *testing.T) {
	if actual := formatSettingsValue(45); actual != "45" {
		t.Fatalf("unexpected whole settings value: %s", actual)
	}
	if actual := formatSettingsValue(72.5); actual != "72,5" {
		t.Fatalf("unexpected decimal settings value: %s", actual)
	}
}

func TestSettingsStorePersistsAndResetsFactorYear(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	preferences := app.Preferences()
	preferences.RemoveValue(factorYearPreferenceKey)
	defaultYear := calculation.DefaultConfig().Year

	store := newSettingsStore(preferences)
	if store.Config().Year != defaultYear {
		t.Fatalf("unexpected default year: %v", store.Config().Year)
	}

	availableYears := calculation.AvailableYears()
	if len(availableYears) == 0 {
		t.Fatal("expected at least one available year")
	}
	store.SetFactorYear(availableYears[0])
	reloaded := newSettingsStore(preferences)
	if reloaded.Config().Year != availableYears[0] {
		t.Fatalf("expected persisted year, got %v", reloaded.Config().Year)
	}

	reloaded.Reset()
	if reloaded.Config().Year != defaultYear {
		t.Fatalf("expected reset default year, got %v", reloaded.Config().Year)
	}
}

func TestSettingsPanelAppliesSelectedFactorYear(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	app.Preferences().RemoveValue(co2PricePreferenceKey)
	app.Preferences().RemoveValue(factorYearPreferenceKey)
	store := newSettingsStore(app.Preferences())
	panel := newSettingsPanel(store, func() {})

	availableYears := calculation.AvailableYears()
	targetYear := availableYears[len(availableYears)-1]
	panel.yearSelect.SetSelected(strconv.Itoa(targetYear))
	panel.apply()
	if store.Config().Year != targetYear {
		t.Fatalf("expected the selected year to be applied, got %v", store.Config().Year)
	}
}

func TestReferenceViewSaveAsWritesRenderedPDF(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	window := app.NewWindow("test")

	original := saveLabelAs
	defer func() { saveLabelAs = original }()

	var savedRecord models.CalculationRecord
	saveLabelAs = func(_ fyne.Window, record models.CalculationRecord, onResult func(string, error)) {
		savedRecord = record
		onResult("/tmp/emission_label.pdf", nil)
	}

	view := buildReferenceView(window, modeOil)
	view.quantityEntry.SetText("10")
	view.calculate()
	test.Tap(view.saveButton)
	if !savedRecord.Valid || savedRecord.FuelType != modeOil {
		t.Fatal("expected save-as to receive the calculated oil result")
	}
	if view.headerStatus.Text != "PDF erstellt" {
		t.Fatalf("unexpected header status after save: %s", view.headerStatus.Text)
	}
}
