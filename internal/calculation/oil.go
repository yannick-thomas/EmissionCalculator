package calculation

import "emissioncalculator/internal/models"

func CalculateOil(quantity float64, cfg Config) models.CalculationResult {
	result := models.CalculationResult{Valid: true}
	emissions := 42.8 * 0.074 * 0.845 * quantity
	energyContent := 42.8 * 0.845 / 1000.0 * quantity * 277.778
	result.Emissions = emissions
	result.EmissionCost = (emissions / 1000.0) * cfg.CO2PricePerTonne * 1.19
	result.EnergyContent = energyContent
	result.CO2PerKWh = emissions / energyContent
	return result
}
