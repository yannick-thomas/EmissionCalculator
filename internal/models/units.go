package models

// Unit identifies an input unit without relying on display text comparisons.
type Unit string

const (
	UnitLitre    Unit = "L"
	UnitTonne    Unit = "t"
	UnitKilogram Unit = "kg"
	UnitKWh      Unit = "kWh"
)

// Quantity is a typed calculation input.
type Quantity struct {
	Value float64 `json:"value"`
	Unit  Unit    `json:"unit"`
}

// Liters represents a volume in litres.
type Liters float64

// Tonnes represents a mass in metric tonnes.
type Tonnes float64

// KgCO2 represents a mass of CO₂ in kilograms.
type KgCO2 float64

// KWh represents an energy quantity in kilowatt-hours.
type KWh float64
