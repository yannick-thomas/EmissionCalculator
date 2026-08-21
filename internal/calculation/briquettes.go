package calculation

import "emissioncalculator/internal/models"

func CalculateBriquettes(quantity float64) models.CalculationResult {
	result := models.CalculationResult{Valid: true}

	emissions := 19.0 * 0.0992 * 1000.0 * quantity
	emissionComponentResult := emissions * 45.0 * 1.19 / 1000.0
	energyContent := 19.0 * 277.778 * quantity

	result.Emissions = FormatWithLocale(emissions, false) + " kg CO2"
	result.EmissionCost = FormatWithLocale(emissionComponentResult, false) + " €"
	result.EnergyContent = FormatWithLocale(energyContent, true) + " kWh"
	result.CO2PerKWh = "0,3571"
	return result
}
