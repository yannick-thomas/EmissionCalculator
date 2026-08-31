package batch

import (
	"emissioncalculator/internal/models"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ImportCSV reads rows with a header containing "fuel"/"brennstoff" and "quantity"/"menge"
// columns (order-insensitive). Missing or unparsable values are kept as empty/raw strings so
// Process can report a precise, row-attributed error instead of silently dropping the row.
func ImportCSV(r io.Reader) ([]Input, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("csv lesen: %w", err)
	}
	if len(rows) == 0 {
		return nil, nil
	}
	fuelIdx, quantityIdx, unitIdx, err := csvColumnIndexes(rows[0])
	if err != nil {
		return nil, err
	}
	inputs := make([]Input, 0, len(rows)-1)
	for i, row := range rows[1:] {
		input := Input{Row: i + 1}
		if fuelIdx < len(row) {
			input.Fuel = strings.TrimSpace(row[fuelIdx])
		}
		if quantityIdx < len(row) {
			input.Quantity = strings.TrimSpace(row[quantityIdx])
		}
		if unitIdx >= 0 && unitIdx < len(row) {
			input.Unit = strings.TrimSpace(row[unitIdx])
		}
		inputs = append(inputs, input)
	}
	return inputs, nil
}

func csvColumnIndexes(header []string) (fuelIdx, quantityIdx, unitIdx int, err error) {
	fuelIdx, quantityIdx, unitIdx = -1, -1, -1
	for i, col := range header {
		switch strings.ToLower(strings.TrimSpace(col)) {
		case "fuel", "brennstoff":
			fuelIdx = i
		case "quantity", "menge":
			quantityIdx = i
		case "unit", "einheit":
			unitIdx = i
		}
	}
	if fuelIdx == -1 || quantityIdx == -1 {
		return -1, -1, -1, fmt.Errorf("csv-header muss die Spalten 'fuel' (oder 'brennstoff') und 'quantity' (oder 'menge') enthalten")
	}
	return fuelIdx, quantityIdx, unitIdx, nil
}

// ExportCSV writes one row per calculation record, including the audit hash for traceability.
func ExportCSV(w io.Writer, records []models.CalculationRecord) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()
	header := []string{"schema_version", "fuel", "quantity", "unit", "emissions_kg", "cost_eur", "energy_kwh", "co2_per_kwh", "calculation_year", "price_reference_year", "factor_pack_id", "factor_valid_from", "factor_source_year", "factor_source_url", "price_source_url", "audit_hash"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("csv schreiben: %w", err)
	}
	for _, r := range records {
		row := []string{
			strconv.Itoa(r.SchemaVersion),
			r.FuelType,
			strconv.FormatFloat(r.Quantity, 'f', 2, 64),
			r.Unit,
			strconv.FormatFloat(float64(r.Emissions), 'f', 2, 64),
			strconv.FormatFloat(r.EmissionCost, 'f', 2, 64),
			strconv.FormatFloat(float64(r.EnergyContent), 'f', 2, 64),
			strconv.FormatFloat(r.CO2PerKWh, 'f', 4, 64),
			strconv.Itoa(r.CalculationYear),
			strconv.Itoa(r.Price.ReferenceYear),
			r.Trace.FactorPackID,
			r.Factor.ValidFrom.Format("2006-01-02"),
			strconv.Itoa(r.Factor.SourceYear),
			r.Factor.SourceURL,
			r.Price.SourceURL,
			r.AuditHash,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("csv schreiben: %w", err)
		}
	}
	if err := writer.Error(); err != nil {
		return fmt.Errorf("csv schreiben: %w", err)
	}
	return nil
}
