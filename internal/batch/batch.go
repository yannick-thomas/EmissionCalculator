package batch

import (
	"emissioncalculator/internal/calculation"
	"emissioncalculator/internal/models"
	"emissioncalculator/internal/validation"
	"fmt"
)

// Input is one row of a batch calculation request, as read from CSV or JSON.
// Fuel accepts the same identifiers as calculation.ParseFuelType, and Quantity accepts both
// plain and German decimal-comma numbers, as parsed by validation.ParseQuantity.
type Input struct {
	Row      int    `json:"-"`
	Fuel     string `json:"fuel"`
	Quantity string `json:"quantity"`
}

// RowError reports a validation or calculation failure for a single input row.
// Row is 1-based and refers to the data row, excluding any header.
type RowError struct {
	Row     int
	Message string
}

func (e RowError) Error() string {
	return fmt.Sprintf("Zeile %d: %s", e.Row, e.Message)
}

// Process validates and calculates every input row against cfg. Each row yields either a
// CalculationRecord or a RowError, so one malformed row never aborts the whole batch.
func Process(inputs []Input, cfg calculation.Config) ([]models.CalculationRecord, []RowError) {
	records := make([]models.CalculationRecord, 0, len(inputs))
	var errs []RowError
	for i, input := range inputs {
		row := input.Row
		if row == 0 {
			row = i + 1
		}
		fuel, err := calculation.ParseFuelType(input.Fuel)
		if err != nil {
			errs = append(errs, RowError{Row: row, Message: err.Error()})
			continue
		}
		quantity, err := validation.ParseQuantity(input.Quantity)
		if err != nil {
			errs = append(errs, RowError{Row: row, Message: err.Error()})
			continue
		}
		record, err := calculation.Calculate(fuel, quantity, cfg)
		if err != nil {
			errs = append(errs, RowError{Row: row, Message: err.Error()})
			continue
		}
		records = append(records, record)
	}
	return records, errs
}
