package calculation

import "time"

// Config holds user-adjustable parameters for emission cost calculation.
type Config struct {
	CO2PricePerTonne float64
	Year             int // calendar year to select the applicable factor set for
}

func DefaultConfig() Config {
	return Config{CO2PricePerTonne: 45.0, Year: time.Now().Year()}
}

// asOf returns the point in time used to resolve the factor set valid for cfg.Year.
func (cfg Config) asOf() time.Time {
	year := cfg.Year
	if year == 0 {
		year = time.Now().Year()
	}
	return time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
}

const (
	vatMultiplier = 1.19   // inkl. 19 % Mehrwertsteuer
	kgPerTonne    = 1000.0 // 1 t = 1000 kg
	mjPerKWh      = 3.6    // 1 kWh = 3,6 MJ
)
