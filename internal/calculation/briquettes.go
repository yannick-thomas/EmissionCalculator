package calculation

import "emissioncalculator/internal/models"

func CalculateBriquettes(quantity float64, cfg Config) models.CalculationResult {
	result := models.CalculationResult{Valid: true}
	emissions := 19.0 * 0.0992 * 1000.0 * quantity
	energyContent := 19.0 * 277.778 * quantity
	result.Emissions = emissions
	result.EmissionCost = emissions * cfg.CO2PricePerTonne * 1.19 / 1000.0
	result.EnergyContent = energyContent
	result.CO2PerKWh = emissions / energyContent
	return result
}
