package calculation

import "emissioncalculator/internal/models"

func CalculateOil(quantity float64) models.CalculationResult {
	result := models.CalculationResult{Valid: true}

	emissions := 42.8 * 0.074 * 0.845 * quantity
	emissionComponentResult := (emissions / 1000.0) * 45.0 * 1.19
	energyContent := 42.8 * 0.845 / 1000.0 * quantity * 277.778

	result.Emissions = FormatWithLocale(emissions, false) + " kg CO2"
	result.EmissionCost = FormatWithLocale(emissionComponentResult, false) + " €"
	result.EnergyContent = FormatWithLocale(energyContent, true) + " kWh"
	result.CO2PerKWh = "0,2664"
	return result
}
