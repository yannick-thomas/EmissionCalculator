package calculation

import (
	"emissioncalculator/internal/models"
	"fmt"
	"strings"
	"time"
)

// FuelType identifies a supported fuel.
type FuelType string

const (
	FuelOil        FuelType = "oil"
	FuelBriquettes FuelType = "briquettes"
	FuelNaturalGas FuelType = "natural_gas"
	FuelLPG        FuelType = "lpg"
)

// factorHistory is built from the validated, versioned factor pack. Publish a new pack
// when official factors change; never rewrite a released pack.
var factorHistory = factorHistoryFromPack(bundledFactorPack)

func factorHistoryFromPack(pack FactorPack) map[FuelType][]FuelFactor {
	history := make(map[FuelType][]FuelFactor)
	for _, factor := range pack.Factors {
		history[factor.Fuel] = append(history[factor.Fuel], factor)
	}
	return history
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
	Fuel      FuelType
	Label     string
	Unit      models.Unit
	Dimension QuantityDimension
}

// Catalog lists every fuel currently supported end-to-end and is derived from the factor pack.
var Catalog = catalogFromPack(bundledFactorPack)

func catalogFromPack(pack FactorPack) []FuelDescriptor {
	seen := make(map[FuelType]bool)
	var catalog []FuelDescriptor
	for _, factor := range pack.Factors {
		if seen[factor.Fuel] {
			continue
		}
		seen[factor.Fuel] = true
		catalog = append(catalog, FuelDescriptor{
			Fuel: factor.Fuel, Label: factor.Label, Unit: factor.DefaultUnit, Dimension: factor.Dimension,
		})
	}
	return catalog
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
	for _, descriptor := range Catalog {
		if normalized == string(descriptor.Fuel) || normalized == strings.ToLower(descriptor.Label) {
			return descriptor.Fuel, nil
		}
	}
	switch normalized {
	case "heizoel":
		return FuelOil, nil
	case "briketts", "braunkohlenbrikett":
		return FuelBriquettes, nil
	case "erdgas", "gas":
		return FuelNaturalGas, nil
	case "flüssiggas", "fluessiggas", "propangas", "propan":
		return FuelLPG, nil
	default:
		return "", fmt.Errorf("unbekannter Brennstoff: %q", value)
	}
}
