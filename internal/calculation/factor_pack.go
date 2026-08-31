package calculation

import (
	"bytes"
	_ "embed"
	"emissioncalculator/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/url"
	"time"
)

const FactorPackSchemaVersion = 1

// QuantityDimension describes the physical dimension of a fuel's default input.
type QuantityDimension string

const (
	DimensionVolume QuantityDimension = "volume"
	DimensionMass   QuantityDimension = "mass"
	DimensionEnergy QuantityDimension = "energy"
)

// UnitConversion converts one supported input unit directly to thermal energy.
type UnitConversion struct {
	Unit            models.Unit `json:"unit"`
	EnergyMJPerUnit float64     `json:"energy_mj_per_unit"`
}

// FuelFactor contains one versioned set of calculation constants for a fuel.
type FuelFactor struct {
	Fuel           FuelType          `json:"fuel"`
	Label          string            `json:"label"`
	Dimension      QuantityDimension `json:"dimension"`
	DefaultUnit    models.Unit       `json:"default_unit"`
	Conversions    []UnitConversion  `json:"conversions"`
	CalorificValue float64           `json:"calorific_value_mj_per_kg,omitempty"`
	Density        float64           `json:"density_kg_per_l,omitempty"`
	EmissionFactor float64           `json:"emission_factor_kg_per_mj"`
	Source         string            `json:"source"`
	SourceURL      string            `json:"source_url"`
	SourceYear     int               `json:"source_year"`
	ValidFrom      time.Time         `json:"valid_from"`
}

// FactorPack is the versioned, auditable source of every physical factor shipped
// with the application.
type FactorPack struct {
	SchemaVersion int          `json:"schema_version"`
	ID            string       `json:"id"`
	PublishedAt   time.Time    `json:"published_at"`
	Factors       []FuelFactor `json:"factors"`
}

//go:embed data/factors-v1.json
var embeddedFactorPackJSON []byte

var bundledFactorPack = mustLoadFactorPack()

func mustLoadFactorPack() FactorPack {
	pack, err := LoadFactorPackJSON(embeddedFactorPackJSON)
	if err != nil {
		panic("invalid embedded factor pack: " + err.Error())
	}
	return pack
}

// LoadFactorPack decodes and validates an externally supplied factor pack.
func LoadFactorPack(reader io.Reader) (FactorPack, error) {
	var pack FactorPack
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&pack); err != nil {
		return FactorPack{}, fmt.Errorf("Faktorpaket lesen: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return FactorPack{}, fmt.Errorf("Faktorpaket lesen: weitere JSON-Daten nach dem Paket")
		}
		return FactorPack{}, fmt.Errorf("Faktorpaket lesen: %w", err)
	}
	if err := ValidateFactorPack(pack); err != nil {
		return FactorPack{}, err
	}
	return pack, nil
}

// LoadFactorPackJSON decodes a factor pack from bytes.
func LoadFactorPackJSON(data []byte) (FactorPack, error) {
	return LoadFactorPack(bytes.NewReader(data))
}

// ValidateFactorPack rejects incomplete, duplicate or physically invalid factors.
func ValidateFactorPack(pack FactorPack) error {
	if pack.SchemaVersion != FactorPackSchemaVersion {
		return fmt.Errorf("nicht unterstützte Faktorpaket-Version %d", pack.SchemaVersion)
	}
	if pack.ID == "" || pack.PublishedAt.IsZero() || len(pack.Factors) == 0 {
		return fmt.Errorf("Faktorpaket benötigt ID, Veröffentlichungsdatum und Faktoren")
	}
	seen := make(map[string]struct{}, len(pack.Factors))
	for _, factor := range pack.Factors {
		key := string(factor.Fuel) + "|" + factor.ValidFrom.UTC().Format(time.RFC3339)
		if _, exists := seen[key]; exists {
			return fmt.Errorf("doppelter Faktor %s", key)
		}
		seen[key] = struct{}{}
		if factor.Fuel == "" || factor.Label == "" || factor.DefaultUnit == "" || factor.ValidFrom.IsZero() {
			return fmt.Errorf("unvollständiger Faktor %s", factor.Fuel)
		}
		switch factor.Dimension {
		case DimensionVolume:
			if factor.Density <= 0 || factor.CalorificValue <= 0 {
				return fmt.Errorf("Volumenfaktor %s benötigt Dichte und Heizwert", factor.Fuel)
			}
		case DimensionMass:
			if factor.CalorificValue <= 0 {
				return fmt.Errorf("Massenfaktor %s benötigt einen Heizwert", factor.Fuel)
			}
		case DimensionEnergy:
		default:
			return fmt.Errorf("unbekannte Mengendimension %q für %s", factor.Dimension, factor.Fuel)
		}
		if factor.EmissionFactor <= 0 || math.IsNaN(factor.EmissionFactor) || math.IsInf(factor.EmissionFactor, 0) {
			return fmt.Errorf("ungültiger Emissionsfaktor für %s", factor.Fuel)
		}
		if factor.Source == "" || factor.SourceYear <= 0 {
			return fmt.Errorf("fehlende Quellenangabe für %s", factor.Fuel)
		}
		if factor.SourceURL != "" {
			parsed, err := url.ParseRequestURI(factor.SourceURL)
			if err != nil || (parsed.Scheme != "https" && parsed.Scheme != "http") {
				return fmt.Errorf("ungültige Quellen-URL für %s", factor.Fuel)
			}
		}
		defaultFound := false
		units := make(map[models.Unit]struct{}, len(factor.Conversions))
		for _, conversion := range factor.Conversions {
			if conversion.Unit == "" || conversion.EnergyMJPerUnit <= 0 || math.IsNaN(conversion.EnergyMJPerUnit) || math.IsInf(conversion.EnergyMJPerUnit, 0) {
				return fmt.Errorf("ungültige Einheitenumrechnung für %s", factor.Fuel)
			}
			if _, exists := units[conversion.Unit]; exists {
				return fmt.Errorf("doppelte Einheitenumrechnung %s für %s", conversion.Unit, factor.Fuel)
			}
			units[conversion.Unit] = struct{}{}
			if conversion.Unit == factor.DefaultUnit {
				defaultFound = true
			}
		}
		if !defaultFound {
			return fmt.Errorf("Standardeinheit %s fehlt bei %s", factor.DefaultUnit, factor.Fuel)
		}
	}
	return nil
}

// CurrentFactorPack returns the immutable factor pack bundled with this build.
func CurrentFactorPack() FactorPack {
	pack := bundledFactorPack
	pack.Factors = append([]FuelFactor(nil), bundledFactorPack.Factors...)
	for index := range pack.Factors {
		pack.Factors[index].Conversions = append([]UnitConversion(nil), bundledFactorPack.Factors[index].Conversions...)
	}
	return pack
}

func (factor FuelFactor) conversionFor(unit models.Unit) (UnitConversion, bool) {
	for _, conversion := range factor.Conversions {
		if conversion.Unit == unit {
			return conversion, true
		}
	}
	return UnitConversion{}, false
}
