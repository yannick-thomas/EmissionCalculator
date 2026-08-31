package calculation

import (
	"strings"
	"testing"
)

func TestEmbeddedFactorPackIsValidAndVersioned(t *testing.T) {
	pack := CurrentFactorPack()
	if pack.SchemaVersion != FactorPackSchemaVersion || pack.ID == "" {
		t.Fatalf("unexpected factor pack metadata: %+v", pack)
	}
	if len(pack.Factors) != 4 {
		t.Fatalf("expected four bundled factors, got %d", len(pack.Factors))
	}
	if err := ValidateFactorPack(pack); err != nil {
		t.Fatalf("embedded factor pack is invalid: %v", err)
	}
}

func TestFactorPackRejectsUnknownFieldsAndMissingConversion(t *testing.T) {
	unknown := `{"schema_version":1,"id":"x","published_at":"2026-01-01T00:00:00Z","factors":[],"unexpected":true}`
	if _, err := LoadFactorPack(strings.NewReader(unknown)); err == nil {
		t.Fatal("expected unknown field to be rejected")
	}

	pack := CurrentFactorPack()
	pack.Factors = append([]FuelFactor(nil), pack.Factors...)
	pack.Factors[0].Conversions = nil
	if err := ValidateFactorPack(pack); err == nil {
		t.Fatal("expected missing default-unit conversion to be rejected")
	}
}

func TestCurrentFactorPackReturnsAnIndependentCopy(t *testing.T) {
	first := CurrentFactorPack()
	first.Factors[0].Label = "verändert"
	first.Factors[0].Conversions[0].EnergyMJPerUnit = 1
	second := CurrentFactorPack()
	if second.Factors[0].Label == "verändert" || second.Factors[0].Conversions[0].EnergyMJPerUnit == 1 {
		t.Fatal("expected callers not to mutate the bundled factor pack")
	}
}
