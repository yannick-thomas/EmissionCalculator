package pdf

import (
	"bytes"
	"emissioncalculator/internal/models"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"
)

func TestExportLabelCreatesPDF(t *testing.T) {
	record := models.CalculationRecord{
		Valid:           true,
		FuelType:        "oil",
		Quantity:        10,
		Unit:            "L",
		Emissions:       models.KgCO2(26.76),
		EmissionCost:    1.43,
		EnergyContent:   models.KWh(100.45),
		CO2PerKWh:       0.2664,
		CO2Price:        45.0,
		CalculationYear: 2024,
		Price: models.PriceSnapshot{
			EURPerTonne: 45, ReferenceYear: 2024, Source: "Bundesregierung / BEHG",
		},
		Factor: models.FactorSnapshot{
			CalorificValue: 42.8, Density: 0.845, EmissionFactor: 0.074,
			Source: "UBA 2022, DIN 51603-1", SourceYear: 2022,
			ValidFrom: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Source:     "UBA 2022, DIN 51603-1",
		ComputedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		AuditHash:  "0123456789abcdef0123456789abcdef",
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

func TestRenderLabelUsesCompactA5LandscapeAndUnicodeFont(t *testing.T) {
	record := models.CalculationRecord{
		Valid: true, FuelType: "oil", Quantity: 100, Unit: "L",
		Emissions: 267.63, EmissionCost: 19.11, EnergyContent: 1004.61, CO2PerKWh: 0.2664,
		CO2Price: 60, CalculationYear: 2026,
		Price:      models.PriceSnapshot{EURPerTonne: 60, ReferenceYear: 2026, Source: "Bundesregierung / BEHG", RangeMin: 55, RangeMax: 65, IsAssumption: true},
		Factor:     models.FactorSnapshot{CalorificValue: 42.8, Density: 0.845, EmissionFactor: 0.074, Source: "UBA 2022, DIN 51603-1", SourceYear: 2022, ValidFrom: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
		ComputedAt: time.Date(2026, 8, 31, 12, 0, 0, 0, time.UTC), AuditHash: "0123456789abcdef",
	}
	data, err := RenderLabel(record)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	if !bytes.HasPrefix(data, []byte("%PDF-")) {
		t.Fatal("expected PDF output")
	}
	if !bytes.Contains(data, []byte("/ToUnicode")) {
		t.Fatal("expected an embedded Unicode font with a ToUnicode map")
	}
	match := regexp.MustCompile(`/MediaBox\s*\[\s*0\s+0\s+([0-9.]+)\s+([0-9.]+)\s*\]`).FindSubmatch(data)
	if len(match) != 3 {
		t.Fatal("expected a readable PDF MediaBox")
	}
	width, _ := strconv.ParseFloat(string(match[1]), 64)
	height, _ := strconv.ParseFloat(string(match[2]), 64)
	if width < 594 || width > 596 || height < 419 || height > 421 {
		t.Fatalf("expected A5 landscape dimensions, got %.2f x %.2f pt", width, height)
	}
}
