package calculation

import "emissioncalculator/internal/models"

// CalculateBriquettes computes CO₂ emissions for lignite briquettes (Braunkohlebriketts).
// quantity is in metric tonnes.
func CalculateBriquettes(quantity float64, cfg Config) models.CalculationRecord {
	record, err := Calculate(FuelBriquettes, quantity, cfg)
	if err != nil {
		return models.CalculationRecord{Valid: false}
	}
	return record
}
