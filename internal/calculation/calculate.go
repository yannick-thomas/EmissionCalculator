package calculation

import (
	"emissioncalculator/internal/models"
	"emissioncalculator/internal/version"
	"fmt"
	"math"
	"strconv"
	"time"
)

// Calculate computes CO₂ emissions, energy content, and cost for any registered fuel.
// quantity is interpreted in the fuel's default unit declared by the active factor pack.
func Calculate(fuel FuelType, quantity float64, cfg Config) (models.CalculationRecord, error) {
	factor, err := FactorFor(fuel, cfg.asOf())
	if err != nil {
		return models.CalculationRecord{}, err
	}
	return CalculateQuantity(fuel, models.Quantity{Value: quantity, Unit: factor.DefaultUnit}, cfg)
}

// CalculateQuantity computes a record from an explicitly typed input quantity.
func CalculateQuantity(fuel FuelType, quantity models.Quantity, cfg Config) (models.CalculationRecord, error) {
	f, err := FactorFor(fuel, cfg.asOf())
	if err != nil {
		return models.CalculationRecord{}, err
	}
	if quantity.Value <= 0 || math.IsNaN(quantity.Value) || math.IsInf(quantity.Value, 0) {
		return models.CalculationRecord{}, fmt.Errorf("Menge muss eine positive endliche Zahl sein")
	}
	conversion, ok := f.conversionFor(quantity.Unit)
	if !ok {
		return models.CalculationRecord{}, fmt.Errorf("Einheit %s wird für %s nicht unterstützt", quantity.Unit, f.Label)
	}
	energyMJ := quantity.Value * conversion.EnergyMJPerUnit
	emissions := energyMJ * f.EmissionFactor
	energyContent := energyMJ / mjPerKWh
	calculationYear := cfg.effectiveYear()
	priceReference := cfg.resolvedPriceReference()
	trace := buildTrace(quantity, f, conversion, energyMJ, emissions, cfg.CO2PricePerTonne)

	record := models.CalculationRecord{
		SchemaVersion:   models.CalculationRecordSchemaVersion,
		Valid:           true,
		FuelType:        string(fuel),
		Quantity:        quantity.Value,
		Unit:            string(quantity.Unit),
		Emissions:       models.KgCO2(emissions),
		EmissionCost:    (emissions / kgPerTonne) * cfg.CO2PricePerTonne * vatMultiplier,
		EnergyContent:   models.KWh(energyContent),
		CO2PerKWh:       emissions / energyContent,
		CO2Price:        cfg.CO2PricePerTonne,
		CalculationYear: calculationYear,
		FactorYear:      calculationYear,
		Price: models.PriceSnapshot{
			EURPerTonne:   priceReference.Amount,
			ReferenceYear: priceReference.ReferenceYear,
			Source:        priceReference.Source,
			SourceURL:     priceReference.SourceURL,
			ValidFrom:     priceReference.ValidFrom,
			ValidUntil:    priceReference.ValidUntil,
			RangeMin:      priceReference.Min,
			RangeMax:      priceReference.Max,
			IsAssumption:  priceReference.IsAssumption,
		},
		Factor: models.FactorSnapshot{
			CalorificValue:  f.CalorificValue,
			Density:         f.Density,
			EmissionFactor:  f.EmissionFactor,
			EnergyMJPerUnit: conversion.EnergyMJPerUnit,
			InputUnit:       quantity.Unit,
			Source:          f.Source,
			SourceURL:       f.SourceURL,
			SourceYear:      f.SourceYear,
			ValidFrom:       f.ValidFrom,
		},
		Trace:      trace,
		Source:     f.Source,
		AppVersion: version.Version,
		ComputedAt: time.Now(),
	}
	record.AuditHash = models.ComputeAuditHash(record)
	return record, nil
}

func buildTrace(quantity models.Quantity, factor FuelFactor, conversion UnitConversion, energyMJ, emissions, price float64) models.CalculationTrace {
	steps := make([]models.CalculationStep, 0, 5)
	input := traceNumber(quantity.Value) + " " + string(quantity.Unit)
	if factor.Density > 0 && quantity.Unit == models.UnitLitre {
		mass := quantity.Value * factor.Density
		steps = append(steps, models.CalculationStep{
			Title: "Brennstoffmasse", Expression: input + " × " + traceNumber(factor.Density) + " kg/L",
			Result: mass, Unit: "kg",
		})
		steps = append(steps, models.CalculationStep{
			Title: "Energiegehalt", Expression: traceNumber(mass) + " kg × " + traceNumber(factor.CalorificValue) + " MJ/kg",
			Result: energyMJ, Unit: "MJ",
		})
	} else if factor.Dimension == DimensionMass && factor.CalorificValue > 0 {
		massKg := quantity.Value
		if quantity.Unit == models.UnitTonne {
			massKg *= kgPerTonne
		}
		steps = append(steps, models.CalculationStep{
			Title: "Brennstoffmasse", Expression: input,
			Result: massKg, Unit: "kg",
		})
		steps = append(steps, models.CalculationStep{
			Title: "Energiegehalt", Expression: traceNumber(massKg) + " kg × " + traceNumber(factor.CalorificValue) + " MJ/kg",
			Result: energyMJ, Unit: "MJ",
		})
	} else {
		steps = append(steps, models.CalculationStep{
			Title: "Energiegehalt", Expression: input + " × " + traceNumber(conversion.EnergyMJPerUnit) + " MJ/" + string(quantity.Unit),
			Result: energyMJ, Unit: "MJ",
		})
	}
	steps = append(steps,
		models.CalculationStep{
			Title: "Gesamtemissionen", Expression: traceNumber(energyMJ) + " MJ × " + traceNumber(factor.EmissionFactor) + " kg CO₂/MJ",
			Result: emissions, Unit: "kg CO₂",
		},
		models.CalculationStep{
			Title: "CO₂-Kosten netto", Expression: traceNumber(emissions) + " kg ÷ 1.000 × " + traceNumber(price) + " €/t",
			Result: emissions / kgPerTonne * price, Unit: "€",
		},
		models.CalculationStep{
			Title: "CO₂-Kosten brutto", Expression: traceNumber(emissions/kgPerTonne*price) + " € × 1,19",
			Result: emissions / kgPerTonne * price * vatMultiplier, Unit: "€",
		},
	)
	return models.CalculationTrace{FactorPackID: bundledFactorPack.ID, Steps: steps}
}

func traceNumber(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
