package ui

import (
	"emissioncalculator/internal/calculation"
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestFuelMenuStartsWithOilAndBriquettes(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	app.Preferences().RemoveValue(enabledFuelsPreferenceKey)

	store := newFuelMenuStore(app.Preferences())
	enabled := store.enabledDescriptors()
	if len(enabled) != 2 || enabled[0].Fuel != calculation.FuelOil || enabled[1].Fuel != calculation.FuelBriquettes {
		t.Fatalf("unexpected initial menu fuels: %+v", enabled)
	}
	if len(store.availableDescriptors()) == 0 {
		t.Fatal("expected further catalog fuels to be available for the menu")
	}
}

func TestFuelMenuPersistsAddedFuel(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	app.Preferences().RemoveValue(enabledFuelsPreferenceKey)

	store := newFuelMenuStore(app.Preferences())
	store.enable(calculation.FuelNaturalGas)
	reloaded := newFuelMenuStore(app.Preferences())
	if !reloaded.enabled[calculation.FuelNaturalGas] {
		t.Fatal("expected an added fuel to remain visible after restart")
	}
}
