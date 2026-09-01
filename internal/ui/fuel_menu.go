package ui

import (
	"emissioncalculator/internal/calculation"
	"strings"

	"fyne.io/fyne/v2"
)

const enabledFuelsPreferenceKey = "ui.enabled_fuels"

// fuelMenuStore persists the catalog entries a user has chosen to show in navigation.
type fuelMenuStore struct {
	preferences fyne.Preferences
	enabled     map[calculation.FuelType]bool
}

func newFuelMenuStore(preferences fyne.Preferences) *fuelMenuStore {
	store := &fuelMenuStore{
		preferences: preferences,
		enabled: map[calculation.FuelType]bool{
			calculation.FuelOil:        true,
			calculation.FuelBriquettes: true,
		},
	}
	if preferences == nil {
		return store
	}

	persisted := preferences.StringWithFallback(enabledFuelsPreferenceKey, "")
	if persisted == "" {
		return store
	}
	store.enabled = make(map[calculation.FuelType]bool)
	for _, value := range strings.Split(persisted, ",") {
		fuel, err := calculation.ParseFuelType(value)
		if err == nil {
			store.enabled[fuel] = true
		}
	}
	store.enabled[calculation.FuelOil] = true
	store.enabled[calculation.FuelBriquettes] = true
	return store
}

func (store *fuelMenuStore) enabledDescriptors() []calculation.FuelDescriptor {
	descriptors := make([]calculation.FuelDescriptor, 0, len(store.enabled))
	for _, descriptor := range calculation.Catalog {
		if store.enabled[descriptor.Fuel] {
			descriptors = append(descriptors, descriptor)
		}
	}
	return descriptors
}

func (store *fuelMenuStore) availableDescriptors() []calculation.FuelDescriptor {
	descriptors := make([]calculation.FuelDescriptor, 0, len(calculation.Catalog))
	for _, descriptor := range calculation.Catalog {
		if !store.enabled[descriptor.Fuel] {
			descriptors = append(descriptors, descriptor)
		}
	}
	return descriptors
}

func (store *fuelMenuStore) enable(fuel calculation.FuelType) {
	store.enabled[fuel] = true
	if store.preferences == nil {
		return
	}

	values := make([]string, 0, len(store.enabled))
	for _, descriptor := range calculation.Catalog {
		if store.enabled[descriptor.Fuel] {
			values = append(values, string(descriptor.Fuel))
		}
	}
	store.preferences.SetString(enabledFuelsPreferenceKey, strings.Join(values, ","))
}
