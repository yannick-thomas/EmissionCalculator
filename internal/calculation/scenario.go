package calculation

import "emissioncalculator/internal/models"

// PriceScenario is one point in a CO2-price comparison for the same fuel and quantity.
type PriceScenario struct {
	Label            string
	CO2PricePerTonne float64
}

// CompareCO2Prices recalculates the same fuel, quantity, and calculation year for each given CO2
// price, e.g. to compare an offer against a lower or higher CO2 price assumption. Each returned
// record carries the scenario's Label so it can be attributed after the call.
func CompareCO2Prices(fuel FuelType, quantity float64, cfg Config, scenarios []PriceScenario) ([]models.CalculationRecord, error) {
	records := make([]models.CalculationRecord, 0, len(scenarios))
	for _, scenario := range scenarios {
		scenarioCfg := cfg
		scenarioCfg = scenarioCfg.WithCO2Price(scenario.CO2PricePerTonne, "Preisszenario: "+scenario.Label)
		record, err := Calculate(fuel, quantity, scenarioCfg)
		if err != nil {
			return nil, err
		}
		record.ScenarioLabel = scenario.Label
		records = append(records, record)
	}
	return records, nil
}
