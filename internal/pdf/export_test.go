package pdf

import (
	"emissioncalculator/internal/models"
	"os"
	"testing"
)

func TestExportLabelCreatesPDF(t *testing.T) {
	record := models.CalculationRecord{
		Valid:         true,
		FuelType:      "oil",
		Quantity:      10,
		Unit:          "L",
		Emissions:     models.KgCO2(26.76),
		EmissionCost:  1.43,
		EnergyContent: models.KWh(100.45),
		CO2PerKWh:     0.2664,
		CO2Price:      45.0,
		Source:        "UBA 2022, DIN 51603-1",
	}
	path, err := ExportLabel(record)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}
	defer os.Remove(path)

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("expected exported PDF: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("expected exported PDF to contain data")
	}
}
