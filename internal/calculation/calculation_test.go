package calculation

import (
	"emissioncalculator/internal/models"
	"emissioncalculator/internal/validation"
	"math"
	"testing"
	"time"
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
	if result.Factor.EmissionFactor != 0.0992 {
		t.Fatalf("unexpected factor snapshot: %v", result.Factor)
	}
}

func TestCalculateOil(t *testing.T) {
	result := CalculateOil(10, DefaultConfig())
	if !result.Valid {
		t.Fatal("expected valid result")
	}
	if math.Abs(float64(result.Emissions)-26.76) > 0.01 {
		t.Fatalf("unexpected emissions: %v", result.Emissions)
	}
	if math.Abs(result.CO2PerKWh-0.2664) > 0.0001 {
		t.Fatalf("unexpected CO2 per kWh: %v", result.CO2PerKWh)
	}
}

func TestCalculateAdditionalFuelsAndTypedUnits(t *testing.T) {
	gas, err := CalculateQuantity(FuelNaturalGas, models.Quantity{Value: 1000, Unit: models.UnitKWh}, DefaultConfig())
	if err != nil {
		t.Fatalf("natural gas calculation failed: %v", err)
	}
	if math.Abs(float64(gas.Emissions)-200.88) > 0.001 || gas.Unit != "kWh" {
		t.Fatalf("unexpected natural gas result: %+v", gas)
	}

	lpg, err := CalculateQuantity(FuelLPG, models.Quantity{Value: 100, Unit: models.UnitKilogram}, DefaultConfig())
	if err != nil {
		t.Fatalf("LPG calculation failed: %v", err)
	}
	if math.Abs(float64(lpg.Emissions)-302.328) > 0.001 || lpg.Unit != "kg" {
		t.Fatalf("unexpected LPG result: %+v", lpg)
	}
	if _, err := CalculateQuantity(FuelNaturalGas, models.Quantity{Value: 10, Unit: models.UnitLitre}, DefaultConfig()); err == nil {
		t.Fatal("expected unsupported unit to fail")
	}
}

func TestCalculationIncludesReproducibleTrace(t *testing.T) {
	record, err := Calculate(FuelOil, 100, DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	if record.SchemaVersion != models.CalculationRecordSchemaVersion || record.Trace.FactorPackID == "" {
		t.Fatalf("missing schema or factor pack metadata: %+v", record.Trace)
	}
	if len(record.Trace.Steps) < 4 {
		t.Fatalf("expected calculation steps, got %+v", record.Trace.Steps)
	}
}

func TestCalculationRecordSeparatesYearAndProvenance(t *testing.T) {
	cfg := DefaultConfig()
	record, err := Calculate(FuelOil, 10, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if record.CalculationYear != cfg.CalculationYear {
		t.Fatalf("unexpected calculation year: %d", record.CalculationYear)
	}
	if record.Factor.ValidFrom.Year() != 2022 || record.Factor.SourceYear != 2022 {
		t.Fatalf("factor validity and source year were not preserved: %+v", record.Factor)
	}
	if record.Price.ReferenceYear == 0 || record.Price.Source == "" || record.Price.EURPerTonne != record.CO2Price {
		t.Fatalf("price provenance was not preserved: %+v", record.Price)
	}
	if cfg.CalculationYear == 2026 && (!record.Price.IsAssumption || record.Price.RangeMin != 55 || record.Price.RangeMax != 65) {
		t.Fatalf("2026 default must be labelled as a corridor assumption: %+v", record.Price)
	}
}

func TestCalculateOilRejectsYearBeforeAnyFactor(t *testing.T) {
	cfg := Config{CO2PricePerTonne: 45, Year: 1999}
	result := CalculateOil(10, cfg)
	if result.Valid {
		t.Fatal("expected an invalid result for a year with no registered factor")
	}
}

func TestFactorForSelectsLatestApplicableVersion(t *testing.T) {
	factor, err := FactorFor(FuelOil, time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if factor.EmissionFactor != 0.074 {
		t.Fatalf("unexpected factor: %v", factor)
	}
	if _, err := FactorFor(FuelOil, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)); err == nil {
		t.Fatal("expected an error for a date before the earliest factor")
	}
}

func TestAvailableYearsIncludesCurrentYear(t *testing.T) {
	years := AvailableYears()
	currentYear := time.Now().Year()
	found := false
	for _, year := range years {
		if year == currentYear {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected %d to be in %v", currentYear, years)
	}
}

func TestParseFuelTypeAcceptsKnownIdentifiersAndLabels(t *testing.T) {
	cases := []struct {
		input    string
		expected FuelType
	}{
		{"oil", FuelOil},
		{"Heizöl", FuelOil},
		{"heizoel", FuelOil},
		{"briquettes", FuelBriquettes},
		{"Briketts", FuelBriquettes},
		{"Erdgas", FuelNaturalGas},
		{"Flüssiggas", FuelLPG},
	}
	for _, tc := range cases {
		fuel, err := ParseFuelType(tc.input)
		if err != nil {
			t.Fatalf("ParseFuelType(%q) unexpected error: %v", tc.input, err)
		}
		if fuel != tc.expected {
			t.Fatalf("ParseFuelType(%q) = %v, expected %v", tc.input, fuel, tc.expected)
		}
	}
	if _, err := ParseFuelType("wasserstoff"); err == nil {
		t.Fatal("expected an error for an unregistered fuel")
	}
}

func TestCalculateGenericMatchesOilAndBriquettesWrappers(t *testing.T) {
	cfg := DefaultConfig()
	genericOil, err := Calculate(FuelOil, 10, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	oilWrapper := CalculateOil(10, cfg)
	if genericOil.Emissions != oilWrapper.Emissions || genericOil.EnergyContent != oilWrapper.EnergyContent ||
		genericOil.EmissionCost != oilWrapper.EmissionCost || genericOil.Factor != oilWrapper.Factor {
		t.Fatalf("Calculate(FuelOil, ...) diverged from CalculateOil: %+v vs %+v", genericOil, oilWrapper)
	}
	genericBriquettes, err := Calculate(FuelBriquettes, 10, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	briquettesWrapper := CalculateBriquettes(10, cfg)
	if genericBriquettes.Emissions != briquettesWrapper.Emissions || genericBriquettes.Factor != briquettesWrapper.Factor {
		t.Fatal("Calculate(FuelBriquettes, ...) diverged from CalculateBriquettes")
	}
}

func TestCompareCO2PricesVariesOnlyThePrice(t *testing.T) {
	cfg := DefaultConfig()
	scenarios := []PriceScenario{
		{Label: "niedrig", CO2PricePerTonne: 30},
		{Label: "hoch", CO2PricePerTonne: 60},
	}
	records, err := CompareCO2Prices(FuelOil, 10, cfg, scenarios)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].Emissions != records[1].Emissions {
		t.Fatal("expected emissions to be unaffected by the CO2 price scenario")
	}
	if records[0].EmissionCost >= records[1].EmissionCost {
		t.Fatal("expected the higher CO2 price scenario to produce a higher cost")
	}
	if records[0].ScenarioLabel != "niedrig" || records[1].ScenarioLabel != "hoch" {
		t.Fatalf("unexpected scenario labels: %q, %q", records[0].ScenarioLabel, records[1].ScenarioLabel)
	}
}

func TestParseQuantityGermanFormat(t *testing.T) {
	cases := []struct {
		input    string
		expected float64
	}{
		{"1,5", 1.5},
		{"12.500,50", 12500.5},
		{"12.5", 12.5},
		{"1000", 1000},
		{"1.000.000", 1000000},
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
