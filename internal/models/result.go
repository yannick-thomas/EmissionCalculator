package models

import "time"

const CalculationRecordSchemaVersion = 1

// CalculationStep is one human-readable, reproducible step of a calculation.
type CalculationStep struct {
	Title      string  `json:"title"`
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
	Unit       string  `json:"unit"`
}

// CalculationTrace explains how inputs, factors and prices produced the result.
type CalculationTrace struct {
	FactorPackID string            `json:"factor_pack_id"`
	Steps        []CalculationStep `json:"steps"`
}

// FactorSnapshot captures the exact physical constants used for one calculation,
// so results remain reproducible and auditable even if factors are revised later.
type FactorSnapshot struct {
	CalorificValue  float64   `json:"calorific_value_mj_per_kg"`
	Density         float64   `json:"density_kg_per_l,omitempty"` // 0 for mass-based fuels
	EmissionFactor  float64   `json:"emission_factor_kg_per_mj"`
	EnergyMJPerUnit float64   `json:"energy_mj_per_input_unit"`
	InputUnit       Unit      `json:"input_unit"`
	Source          string    `json:"source"`
	SourceURL       string    `json:"source_url,omitempty"`
	SourceYear      int       `json:"source_year"`
	ValidFrom       time.Time `json:"valid_from"`
}

// PriceSnapshot captures the monetary assumption independently from the
// calculation year and physical emission-factor provenance.
type PriceSnapshot struct {
	EURPerTonne   float64   `json:"eur_per_tonne"`
	ReferenceYear int       `json:"reference_year"`
	Source        string    `json:"source"`
	SourceURL     string    `json:"source_url,omitempty"`
	ValidFrom     time.Time `json:"valid_from,omitempty"`
	ValidUntil    time.Time `json:"valid_until,omitempty"`
	RangeMin      float64   `json:"range_min_eur_per_tonne,omitempty"`
	RangeMax      float64   `json:"range_max_eur_per_tonne,omitempty"`
	IsAssumption  bool      `json:"is_assumption"`
}

// CalculationRecord is the canonical output of a fuel calculation.
// It carries all inputs, physical results, and provenance metadata for UI, PDF, and export.
type CalculationRecord struct {
	SchemaVersion   int              `json:"schema_version"`
	Valid           bool             `json:"valid"`
	FuelType        string           `json:"fuel_type"`
	Quantity        float64          `json:"quantity"` // input quantity in the fuel's native unit
	Unit            string           `json:"unit"`     // unit label, e.g. "L" or "t"
	Emissions       KgCO2            `json:"emissions_kg"`
	EmissionCost    float64          `json:"emission_cost_eur"` // incl. VAT
	EnergyContent   KWh              `json:"energy_content_kwh"`
	CO2PerKWh       float64          `json:"co2_per_kwh"`             // intensity: kg CO₂/kWh
	CO2Price        float64          `json:"co2_price_eur_per_tonne"` // compatible alias for Price.EURPerTonne
	CalculationYear int              `json:"calculation_year"`
	FactorYear      int              `json:"factor_year,omitempty"` // deprecated compatibility alias for CalculationYear
	Price           PriceSnapshot    `json:"price"`
	Factor          FactorSnapshot   `json:"factor"`
	Trace           CalculationTrace `json:"trace"`
	Source          string           `json:"source"`                   // physical constants reference (kept for backward-compatible display)
	ScenarioLabel   string           `json:"scenario_label,omitempty"` // set when part of a scenario comparison
	AppVersion      string           `json:"app_version"`
	AuditHash       string           `json:"audit_hash"` // fingerprint for internal traceability, not a legal certification
	ComputedAt      time.Time        `json:"computed_at"`
}
