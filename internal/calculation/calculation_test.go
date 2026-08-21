package calculation

import (
	"emissioncalculator/internal/validation"
	"math"
	"testing"
)

func TestCalculateBriquettes(t *testing.T) {
	result := CalculateBriquettes(10, DefaultConfig())
	if !result.Valid {
		t.Fatal("expected valid result")
	}
	if result.Emissions != 18848.0 {
		t.Fatalf("unexpected emissions: %v", result.Emissions)
	}
	if math.Abs(result.CO2PerKWh-0.3571) > 0.0001 {
		t.Fatalf("unexpected CO2 per kWh: %v", result.CO2PerKWh)
	}
}

func TestCalculateOil(t *testing.T) {
	result := CalculateOil(10, DefaultConfig())
	if !result.Valid {
		t.Fatal("expected valid result")
	}
	if math.Abs(result.Emissions-26.76) > 0.01 {
		t.Fatalf("unexpected emissions: %v", result.Emissions)
	}
	if math.Abs(result.CO2PerKWh-0.2664) > 0.0001 {
		t.Fatalf("unexpected CO2 per kWh: %v", result.CO2PerKWh)
	}
}

func TestParseQuantityGermanFormat(t *testing.T) {
	cases := []struct {
		input    string
		expected float64
	}{
		{"1,5", 1.5},
		{"12.500", 12500},
		{"12.500,50", 12500.5},
		{"12.5", 12.5},
		{"1000", 1000},
		{"1.234.567", 1234567},
	}
	for _, tc := range cases {
		value, err := validation.ParseQuantity(tc.input)
		if err != nil {
			t.Fatalf("ParseQuantity(%q) unexpected error: %v", tc.input, err)
		}
		if value != tc.expected {
			t.Fatalf("ParseQuantity(%q) = %v, expected %v", tc.input, value, tc.expected)
		}
	}
}
