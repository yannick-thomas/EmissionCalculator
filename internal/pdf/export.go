package pdf

import (
	"bytes"
	"emissioncalculator/internal/models"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// RenderLabel renders a calculation record as PDF bytes.
func RenderLabel(r models.CalculationRecord) ([]byte, error) {
	doc := gofpdf.New("P", "mm", "A4", "")
	doc.AddPage()
	doc.SetFont("Arial", "B", 16)
	doc.MultiCell(0, 8, fmt.Sprintf("CO2-Kosten: %.2f \u20ac brutto", r.EmissionCost), "", "L", false)
	doc.SetFont("Arial", "", 12)
	doc.MultiCell(0, 7, fmt.Sprintf("Gesamtemissionen: %.2f kg CO2", float64(r.Emissions)), "", "L", false)
	doc.MultiCell(0, 7, fmt.Sprintf("Emissionsintensitaet: %.4f kg CO2/kWh", r.CO2PerKWh), "", "L", false)
	doc.MultiCell(0, 7, fmt.Sprintf("Energiegehalt: %.2f kWh", float64(r.EnergyContent)), "", "L", false)

	doc.Ln(4)
	doc.SetFont("Arial", "B", 11)
	doc.MultiCell(0, 6, "Berechnungsnachweis", "", "L", false)
	doc.SetFont("Arial", "", 10)
	doc.MultiCell(0, 6, fmt.Sprintf("Menge: %.2f %s", r.Quantity, r.Unit), "", "L", false)
	doc.MultiCell(0, 6, fmt.Sprintf("Heizwert Hu: %.1f MJ/kg", r.Factor.CalorificValue), "", "L", false)
	if r.Factor.Density > 0 {
		doc.MultiCell(0, 6, fmt.Sprintf("Dichte: %.3f kg/L", r.Factor.Density), "", "L", false)
	}
	doc.MultiCell(0, 6, fmt.Sprintf("Emissionsfaktor: %.4f kg CO2/MJ", r.Factor.EmissionFactor), "", "L", false)
	doc.MultiCell(0, 6, fmt.Sprintf("CO2-Preis: %.2f EUR/t zzgl. 19%% MwSt.", r.CO2Price), "", "L", false)
	doc.MultiCell(0, 6, fmt.Sprintf("Faktorjahr: %d (gültig ab %s)", r.FactorYear, r.Factor.ValidFrom.Format("02.01.2006")), "", "L", false)
	doc.MultiCell(0, 6, "Quelle: "+r.Source, "", "L", false)
	doc.MultiCell(0, 6, "Berechnet am: "+r.ComputedAt.Format("02.01.2006 15:04"), "", "L", false)

	doc.Ln(2)
	doc.SetFont("Arial", "", 8)
	auditID := r.AuditHash
	if len(auditID) > 16 {
		auditID = auditID[:16]
	}
	doc.MultiCell(0, 5, fmt.Sprintf("Audit-ID: %s · App-Version: %s", auditID, r.AppVersion), "", "L", false)
	doc.MultiCell(0, 5, "Hinweis: Audit-ID und App-Version dienen der internen Nachvollziehbarkeit und stellen keine amtliche Zertifizierung dar.", "", "L", false)

	var buf bytes.Buffer
	if err := doc.Output(&buf); err != nil {
		return nil, fmt.Errorf("pdf render: %w", err)
	}
	return buf.Bytes(), nil
}

// ExportLabel renders a calculation record as PDF, saves it to a temp file, and returns the path.
func ExportLabel(r models.CalculationRecord) (string, error) {
	data, err := RenderLabel(r)
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("emission_label_%d.pdf", time.Now().UnixNano())
	filePath := filepath.Join(os.TempDir(), fileName)
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return "", fmt.Errorf("pdf save %s: %w", filePath, err)
	}
	return filePath, nil
}
