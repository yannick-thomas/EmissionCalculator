package calculation

import (
	"emissioncalculator/internal/validation"
	"testing"
)

func TestCalculateBriquettes(t *testing.T) {
	result := CalculateBriquettes(10)
	if !result.Valid {
		t.Fatal("expected valid result")
	}
	if result.Emissions != "18848,00 kg CO2" {
		t.Fatalf("unexpected emissions: %s", result.Emissions)
	}
	if result.CO2PerKWh != "0,3571" {
		t.Fatalf("unexpected CO2 per kWh: %s", result.CO2PerKWh)
	}
}

func TestCalculateOil(t *testing.T) {
	result := CalculateOil(10)
	if !result.Valid {
		t.Fatal("expected valid result")
	}
	if result.Emissions != "26,76 kg CO2" {
		t.Fatalf("unexpected emissions: %s", result.Emissions)
	}
	if result.CO2PerKWh != "0,2664" {
		t.Fatalf("unexpected CO2 per kWh: %s", result.CO2PerKWh)
	}
}

func TestParseQuantity(t *testing.T) {
	value, err := validation.ParseQuantity("1,5")
	if err != nil {
		t.Fatalf("expected valid quantity, got error: %v", err)
	}
	if value != 1.5 {
		t.Fatalf("expected 1.5, got %v", value)
	}
}
