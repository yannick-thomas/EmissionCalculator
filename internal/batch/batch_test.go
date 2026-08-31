package batch

import (
	"bytes"
	"emissioncalculator/internal/calculation"
	"strings"
	"testing"
)

func TestImportCSVAndProcessReportsPerRowErrors(t *testing.T) {
	csvData := "fuel,quantity,unit\noil,10,L\nbriquettes,2.5,t\nerdgas,5,L\noil,-3,L\n"
	inputs, err := ImportCSV(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("unexpected import error: %v", err)
	}
	if len(inputs) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(inputs))
	}

	records, rowErrors := Process(inputs, calculation.DefaultConfig())
	if len(records) != 2 {
		t.Fatalf("expected 2 successful records, got %d", len(records))
	}
	if len(rowErrors) != 2 {
		t.Fatalf("expected 2 row errors, got %d: %v", len(rowErrors), rowErrors)
	}
	if rowErrors[0].Row != 3 || rowErrors[1].Row != 4 {
		t.Fatalf("unexpected row numbers: %+v", rowErrors)
	}
}

func TestImportCSVRejectsMissingHeader(t *testing.T) {
	if _, err := ImportCSV(strings.NewReader("a,b\n1,2\n")); err == nil {
		t.Fatal("expected an error for a header without fuel/quantity columns")
	}
}

func TestImportJSONAndProcess(t *testing.T) {
	jsonData := `[{"fuel":"oil","quantity":"10"},{"fuel":"briquettes","quantity":"2,5"}]`
	inputs, err := ImportJSON(strings.NewReader(jsonData))
	if err != nil {
		t.Fatalf("unexpected import error: %v", err)
	}
	records, rowErrors := Process(inputs, calculation.DefaultConfig())
	if len(rowErrors) != 0 {
		t.Fatalf("expected no row errors, got %v", rowErrors)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
}

func TestExportCSVAndExportJSONRoundtrip(t *testing.T) {
	records, rowErrors := Process([]Input{{Fuel: "oil", Quantity: "10"}}, calculation.DefaultConfig())
	if len(rowErrors) != 0 {
		t.Fatalf("unexpected row errors: %v", rowErrors)
	}

	var csvBuf bytes.Buffer
	if err := ExportCSV(&csvBuf, records); err != nil {
		t.Fatalf("unexpected CSV export error: %v", err)
	}
	if !strings.Contains(csvBuf.String(), "oil,10.00") {
		t.Fatalf("unexpected CSV output: %s", csvBuf.String())
	}

	var jsonBuf bytes.Buffer
	if err := ExportJSON(&jsonBuf, records); err != nil {
		t.Fatalf("unexpected JSON export error: %v", err)
	}
	if !strings.Contains(jsonBuf.String(), `"fuel_type": "oil"`) {
		t.Fatalf("unexpected JSON output: %s", jsonBuf.String())
	}
}
