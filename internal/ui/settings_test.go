package ui

import (
	"math"
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestSettingsStorePersistsAndResetsCO2Price(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	preferences := app.Preferences()
	preferences.RemoveValue(co2PricePreferenceKey)

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
	initialCost := view.result.EmissionCost

	store.SetCO2Price(90)
	view.refreshForSettingsChange()
	if math.Abs(view.result.EmissionCost-initialCost*2) > 0.000001 {
		t.Fatalf("expected the updated price to recalculate the result: %v", view.result.EmissionCost)
	}
	if view.headerStatus.Text != "Aktualisiert" {
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
