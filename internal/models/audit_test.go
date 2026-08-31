package models

import "testing"

func sampleRecord() CalculationRecord {
	return CalculationRecord{
		Valid:         true,
		FuelType:      "oil",
		Quantity:      10,
		Unit:          "L",
		Emissions:     KgCO2(26.76),
		EmissionCost:  1.43,
		EnergyContent: KWh(100.45),
		CO2PerKWh:     0.2664,
		CO2Price:      45,
		FactorYear:    2026,
		Factor: FactorSnapshot{
			CalorificValue: 42.8,
			Density:        0.845,
			EmissionFactor: 0.074,
			Source:         "UBA 2022, DIN 51603-1",
		},
		Source:     "UBA 2022, DIN 51603-1",
		AppVersion: "test",
	}
}

func TestComputeAuditHashIsDeterministic(t *testing.T) {
	record := sampleRecord()
	first := ComputeAuditHash(record)
	second := ComputeAuditHash(record)
	if first != second {
		t.Fatalf("expected the same record to hash identically, got %q and %q", first, second)
	}
	if len(first) != 64 {
		t.Fatalf("expected a 64-character hex sha256 digest, got %d chars", len(first))
	}
}

func TestComputeAuditHashChangesWithInputs(t *testing.T) {
	base := sampleRecord()
	changed := base
	changed.Quantity = 20

	if ComputeAuditHash(base) == ComputeAuditHash(changed) {
		t.Fatal("expected a different quantity to produce a different audit hash")
	}
}
