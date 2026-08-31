package calculation

import "time"

const federalGovernmentPriceSource = "Bundesregierung / Brennstoffemissionshandelsgesetz (BEHG)"

// PriceReference documents where a configured CO₂ price comes from. Amount is the
// value used by the calculation. Min/Max describe an official corridor when no
// single statutory price exists; IsAssumption makes that distinction explicit.
type PriceReference struct {
	Amount        float64
	ReferenceYear int
	Source        string
	SourceURL     string
	ValidFrom     time.Time
	ValidUntil    time.Time
	Min           float64
	Max           float64
	IsAssumption  bool
}

// Config holds user-adjustable parameters for emission cost calculation.
type Config struct {
	CO2PricePerTonne float64
	CalculationYear  int
	Year             int // Deprecated: use CalculationYear; retained for callers using older Config literals.
	PriceReference   PriceReference
}

func DefaultConfig() Config {
	year := time.Now().Year()
	reference := defaultPriceReference(year)
	return Config{CO2PricePerTonne: reference.Amount, CalculationYear: year, Year: year, PriceReference: reference}
}

// WithCO2Price returns cfg with a documented manual price assumption. It is used
// for user-entered prices and scenarios so an official reference is never retained
// accidentally after the numeric value has changed.
func (cfg Config) WithCO2Price(price float64, description string) Config {
	cfg.CO2PricePerTonne = price
	if description == "" {
		description = "Manuelle Preisannahme"
	}
	cfg.PriceReference = PriceReference{
		Amount:        price,
		ReferenceYear: cfg.effectiveYear(),
		Source:        description,
		IsAssumption:  true,
	}
	return cfg
}

// resolvedPriceReference normalizes legacy Config literals that only contain
// CO2PricePerTonne and Year.
func (cfg Config) resolvedPriceReference() PriceReference {
	reference := cfg.PriceReference
	if reference.Amount == 0 {
		reference = cfg.WithCO2Price(cfg.CO2PricePerTonne, "Manuelle Preisannahme").PriceReference
	}
	reference.Amount = cfg.CO2PricePerTonne
	return reference
}

func (cfg Config) effectiveYear() int {
	if cfg.CalculationYear != 0 {
		return cfg.CalculationYear
	}
	if cfg.Year != 0 {
		return cfg.Year
	}
	return time.Now().Year()
}

// defaultPriceReference returns the newest documented default that is applicable
// to the requested year. For 2026 no fixed statutory price exists; 60 EUR/t is an
// explicitly labelled midpoint assumption within the official 55–65 EUR/t corridor.
func defaultPriceReference(year int) PriceReference {
	references := []PriceReference{
		{
			Amount: 30, ReferenceYear: 2022, Source: federalGovernmentPriceSource,
			ValidFrom: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), ValidUntil: time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			Amount: 45, ReferenceYear: 2024, Source: federalGovernmentPriceSource,
			SourceURL: "https://www.bundesregierung.de/breg-de/aktuelles/neues-gebaeudeenergiegesetz-2184942",
			ValidFrom: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), ValidUntil: time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			Amount: 55, ReferenceYear: 2025, Source: federalGovernmentPriceSource,
			SourceURL: "https://www.bundesregierung.de/breg-de/aktuelles/neues-gebaeudeenergiegesetz-2184942",
			ValidFrom: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), ValidUntil: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			Amount: 60, ReferenceYear: 2026, Source: federalGovernmentPriceSource,
			SourceURL: "https://www.bundesregierung.de/breg-de/suche/gesetzliche-neuregelungen-januar-2026-2399838",
			ValidFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), ValidUntil: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
			Min: 55, Max: 65, IsAssumption: true,
		},
	}
	selected := references[0]
	for _, candidate := range references {
		if candidate.ReferenceYear <= year && candidate.ReferenceYear >= selected.ReferenceYear {
			selected = candidate
		}
	}
	return selected
}

// asOf returns the point in time used to resolve the factor set valid for the calculation year.
func (cfg Config) asOf() time.Time {
	year := cfg.effectiveYear()
	return time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
}

const (
	vatMultiplier = 1.19   // inkl. 19 % Mehrwertsteuer
	kgPerTonne    = 1000.0 // 1 t = 1000 kg
	mjPerKWh      = 3.6    // 1 kWh = 3,6 MJ
)
