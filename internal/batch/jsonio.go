package batch

import (
	"emissioncalculator/internal/models"
	"encoding/json"
	"fmt"
	"io"
)

// ImportJSON decodes a JSON array of {"fuel": "...", "quantity": "..."} objects.
func ImportJSON(r io.Reader) ([]Input, error) {
	var inputs []Input
	if err := json.NewDecoder(r).Decode(&inputs); err != nil {
		return nil, fmt.Errorf("json lesen: %w", err)
	}
	for i := range inputs {
		inputs[i].Row = i + 1
	}
	return inputs, nil
}

// ExportJSON writes the given calculation records as an indented JSON array, including the
// applied factor snapshot and audit hash for each record.
func ExportJSON(w io.Writer, records []models.CalculationRecord) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(records); err != nil {
		return fmt.Errorf("json schreiben: %w", err)
	}
	return nil
}
