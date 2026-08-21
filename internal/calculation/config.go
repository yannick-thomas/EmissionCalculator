package calculation

// Config holds user-adjustable parameters for emission cost calculation.
type Config struct {
	CO2PricePerTonne float64
}

func DefaultConfig() Config {
	return Config{CO2PricePerTonne: 45.0}
}
