package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// ComputeAuditHash derives a deterministic fingerprint over a record's inputs, results, and the
// exact factor set applied. It exists for internal traceability (e.g. "which exact factors and
// app version produced this number") and is not a legal or regulatory certification.
func ComputeAuditHash(r CalculationRecord) string {
	canonical := fmt.Sprintf(
		"%s|%.6f|%s|%.6f|%.6f|%.6f|%.6f|%.6f|%d|%.6f|%.6f|%.6f|%s|%s|%s|%s",
		r.FuelType, r.Quantity, r.Unit,
		float64(r.Emissions), r.EmissionCost, float64(r.EnergyContent), r.CO2PerKWh, r.CO2Price,
		r.FactorYear, r.Factor.CalorificValue, r.Factor.Density, r.Factor.EmissionFactor,
		r.Factor.ValidFrom.UTC().Format("2006-01-02"),
		r.AppVersion,
		r.ComputedAt.UTC().Format("2006-01-02T15:04:05.000000000Z07:00"),
		r.Source,
	)
	sum := sha256.Sum256([]byte(canonical))
	return hex.EncodeToString(sum[:])
}
