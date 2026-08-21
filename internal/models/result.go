package models

type CalculationResult struct {
	Valid         bool
	Emissions     string
	EmissionCost  string
	EnergyContent string
	CO2PerKWh     string
	ErrorMessage  string
}
