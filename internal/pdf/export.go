package pdf

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func ExportLabel(emissions, emissionCost, energyContent, co2PerKWh string) (string, error) {
	tempDir := os.TempDir()
	fileName := fmt.Sprintf("emission_label_%d.pdf", time.Now().UnixNano())
	filePath := filepath.Join(tempDir, fileName)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.MultiCell(0, 8, "CO2 Abgabe: "+emissionCost+" (Brutto)", "", "L", false)
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 7, "CO2 kg der Lieferung: "+emissions, "", "L", false)
	pdf.MultiCell(0, 7, "CO2 kg pro kWh: "+co2PerKWh, "", "L", false)
	pdf.MultiCell(0, 7, "kWh der Lieferung: "+energyContent, "", "L", false)

	if err := pdf.OutputFileAndClose(filePath); err != nil {
		return "", err
	}

	return filePath, nil
}
