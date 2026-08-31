package calculation

import "emissioncalculator/internal/models"

// CalculateOil computes CO₂ emissions for heating oil (Heizöl EL).
// quantity is in litres.
func CalculateOil(quantity float64, cfg Config) models.CalculationRecord {
	record, err := Calculate(FuelOil, quantity, cfg)
	if err != nil {
		return models.CalculationRecord{Valid: false}
	}
	return record
}
