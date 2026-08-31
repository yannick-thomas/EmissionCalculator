package calculation

import (
	"emissioncalculator/internal/models"
	"emissioncalculator/internal/version"
	"time"
)

// Calculate computes CO₂ emissions, energy content, and cost for any registered fuel.
// quantity is in the fuel's native unit (litres for volume-based fuels, tonnes for mass-based
// fuels, distinguished by whether the resolved FuelFactor has a non-zero Density).
func Calculate(fuel FuelType, quantity float64, cfg Config) (models.CalculationRecord, error) {
	f, err := FactorFor(fuel, cfg.asOf())
	if err != nil {
		return models.CalculationRecord{}, err
	}
	massKg := quantity * kgPerTonne
	if f.Density > 0 {
		massKg = quantity * f.Density
	}
	energyMJ := f.CalorificValue * massKg
	emissions := energyMJ * f.EmissionFactor
	energyContent := energyMJ / mjPerKWh

	record := models.CalculationRecord{
		Valid:         true,
		FuelType:      string(fuel),
		Quantity:      quantity,
		Unit:          f.Unit,
		Emissions:     models.KgCO2(emissions),
		EmissionCost:  (emissions / kgPerTonne) * cfg.CO2PricePerTonne * vatMultiplier,
		EnergyContent: models.KWh(energyContent),
		CO2PerKWh:     emissions / energyContent,
		CO2Price:      cfg.CO2PricePerTonne,
		FactorYear:    cfg.Year,
		Factor: models.FactorSnapshot{
			CalorificValue: f.CalorificValue,
			Density:        f.Density,
			EmissionFactor: f.EmissionFactor,
			Source:         f.Source,
			ValidFrom:      f.ValidFrom,
		},
		Source:     f.Source,
		AppVersion: version.Version,
		ComputedAt: time.Now(),
	}
	record.AuditHash = models.ComputeAuditHash(record)
	return record, nil
}
