package pdf

import (
	"os"
	"testing"
)

func TestExportLabelCreatesPDF(t *testing.T) {
	path, err := ExportLabel("26,76 kg CO2", "1,43 €", "9.995,30 kWh", "0,2664")
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
