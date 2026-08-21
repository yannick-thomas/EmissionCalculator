package models

type CalculationResult struct {
	Valid         bool
	Emissions     float64 // kg CO2
	EmissionCost  float64 // €
	EnergyContent float64 // kWh
	CO2PerKWh     float64 // kg CO2/kWh
}
