package calculation

import (
	"fmt"
	"strings"
	"time"
)

// FuelType identifies a supported fuel.
type FuelType string

const (
	FuelOil        FuelType = "oil"
	FuelBriquettes FuelType = "briquettes"
)

// FuelFactor holds physical constants and regulatory metadata for a single fuel,
// valid from ValidFrom until superseded by a later entry for the same fuel.
type FuelFactor struct {
	Fuel           FuelType
	Unit           string  // native input unit label ("L" or "t")
	CalorificValue float64 // lower heating value in MJ/kg
	Density        float64 // kg/L (0 for mass-based fuels)
	EmissionFactor float64 // kg CO₂/MJ
	Source         string
	ValidFrom      time.Time
}

// factorHistory holds every known factor version per fuel, in no particular order.
// Append new entries when official emission factors are revised; never mutate existing ones,
// so that past calculations stay reproducible from their recorded FactorYear.
var factorHistory = map[FuelType][]FuelFactor{
	FuelOil: {
		{
			Fuel:           FuelOil,
			Unit:           "L",
			CalorificValue: 42.8,
			Density:        0.845,
			EmissionFactor: 0.074,
			Source:         "UBA 2022, DIN 51603-1",
			ValidFrom:      time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	},
	FuelBriquettes: {
		{
			Fuel:           FuelBriquettes,
			Unit:           "t",
			CalorificValue: 19.0,
			EmissionFactor: 0.0992,
			Source:         "UBA 2022",
			ValidFrom:      time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	},
}

// FactorFor returns the factor set for fuel that was valid at the given point in time:
// the entry with the latest ValidFrom that is not after at.
func FactorFor(fuel FuelType, at time.Time) (FuelFactor, error) {
	history := factorHistory[fuel]
	var selected FuelFactor
	found := false
	for _, candidate := range history {
		if candidate.ValidFrom.After(at) {
			continue
		}
		if !found || candidate.ValidFrom.After(selected.ValidFrom) {
			selected = candidate
			found = true
		}
	}
	if !found {
		return FuelFactor{}, fmt.Errorf("kein Emissionsfaktor für %s zum Stichtag %s hinterlegt", fuel, at.Format("2006-01-02"))
	}
	return selected, nil
}

// AvailableYears returns the sorted calendar years for which at least one fuel has a factor
// set, ranging from the earliest registered ValidFrom through the current year.
func AvailableYears() []int {
	now := time.Now()
	earliest := now.Year()
	for _, history := range factorHistory {
		for _, f := range history {
			if f.ValidFrom.Year() < earliest {
				earliest = f.ValidFrom.Year()
			}
		}
	}
	years := make([]int, 0, now.Year()-earliest+1)
	for year := earliest; year <= now.Year(); year++ {
		years = append(years, year)
	}
	return years
}

// FuelDescriptor describes a fuel selectable by users, the CLI, or batch import, independent of
// its factor history.
type FuelDescriptor struct {
	Fuel  FuelType
	Label string // German display label
	Unit  string // native input unit label ("L" or "t")
}

// Catalog lists every fuel currently supported end-to-end. To add a further Brennstoff (e.g. an
// EBeV briquette subtype), register its verified FuelFactor in factorHistory and add an entry
// here; never invent emission factor numbers without an official source.
var Catalog = []FuelDescriptor{
	{Fuel: FuelOil, Label: "Heizöl", Unit: "L"},
	{Fuel: FuelBriquettes, Label: "Briketts", Unit: "t"},
}

// FuelByType looks up a catalog descriptor by FuelType.
func FuelByType(fuel FuelType) (FuelDescriptor, bool) {
	for _, descriptor := range Catalog {
		if descriptor.Fuel == fuel {
			return descriptor, true
		}
	}
	return FuelDescriptor{}, false
}

// ParseFuelType parses a case-insensitive fuel identifier, accepting both the FuelType value and
// common German labels (as used in CSV/JSON batch import).
func ParseFuelType(value string) (FuelType, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case string(FuelOil), "heizöl", "heizoel":
		return FuelOil, nil
	case string(FuelBriquettes), "briketts":
		return FuelBriquettes, nil
	default:
		return "", fmt.Errorf("unbekannter Brennstoff: %q", value)
	}
}
