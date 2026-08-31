package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// ComputeAuditHash derives a deterministic fingerprint over a record's inputs, results, and the
// exact factor set applied. It exists for internal traceability (e.g. "which exact factors and
// app version produced this number") and is not a legal or regulatory certification.
func ComputeAuditHash(r CalculationRecord) string {
	canonical, _ := json.Marshal(struct {
		SchemaVersion   int
		FuelType        string
		Quantity        float64
		Unit            string
		Emissions       KgCO2
		EmissionCost    float64
		EnergyContent   KWh
		CO2PerKWh       float64
		CalculationYear int
		Price           PriceSnapshot
		Factor          FactorSnapshot
		Trace           CalculationTrace
		AppVersion      string
		ComputedAt      string
	}{
		SchemaVersion: r.SchemaVersion,
		FuelType:      r.FuelType, Quantity: r.Quantity, Unit: r.Unit,
		Emissions: r.Emissions, EmissionCost: r.EmissionCost,
		EnergyContent: r.EnergyContent, CO2PerKWh: r.CO2PerKWh,
		CalculationYear: r.CalculationYear, Price: r.Price, Factor: r.Factor, Trace: r.Trace,
		AppVersion: r.AppVersion,
		ComputedAt: r.ComputedAt.UTC().Format("2006-01-02T15:04:05.000000000Z07:00"),
	})
	sum := sha256.Sum256(canonical)
	return hex.EncodeToString(sum[:])
}
