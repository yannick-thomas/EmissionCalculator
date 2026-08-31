package pdf

import (
	"bytes"
	"emissioncalculator/internal/models"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/gomedium"
)

const (
	pageWidthMM  = 210.0
	pageHeightMM = 148.0
	fontRegular  = "GoRegular"
	fontBold     = "GoBold"
)

// RenderLabel renders a calculation record as PDF bytes.
func RenderLabel(r models.CalculationRecord) ([]byte, error) {
	doc := gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: "P",
		UnitStr:        "mm",
		Size:           gofpdf.SizeType{Wd: pageWidthMM, Ht: pageHeightMM},
	})
	doc.SetMargins(14, 12, 14)
	doc.SetAutoPageBreak(false, 10)
	doc.AddUTF8FontFromBytes(fontRegular, "", gomedium.TTF)
	doc.AddUTF8FontFromBytes(fontBold, "", gobold.TTF)
	doc.AddPage()

	doc.SetTextColor(24, 33, 61)
	doc.SetFont(fontBold, "", 17)
	doc.CellFormat(0, 8, "Emissionsrechner · Berechnungsnachweis", "", 1, "L", false, 0, "")
	doc.SetFont(fontRegular, "", 9)
	doc.SetTextColor(97, 102, 115)
	doc.CellFormat(0, 5, fmt.Sprintf("Berechnet am %s · Berechnungsjahr %d", r.ComputedAt.Format("02.01.2006 15:04"), r.CalculationYear), "", 1, "L", false, 0, "")

	doc.SetFillColor(217, 238, 138)
	doc.RoundedRect(14, 28, 182, 30, 3, "1234", "F")
	doc.SetXY(20, 34)
	doc.SetTextColor(24, 33, 61)
	doc.SetFont(fontBold, "", 21)
	doc.CellFormat(94, 9, formatNumber(float64(r.Emissions), 2)+" kg CO₂", "", 0, "L", false, 0, "")
	doc.SetFont(fontBold, "", 16)
	doc.CellFormat(62, 9, formatNumber(r.EmissionCost, 2)+" € brutto", "", 1, "R", false, 0, "")
	doc.SetX(20)
	doc.SetFont(fontRegular, "", 9)
	doc.CellFormat(94, 6, "Gesamtemissionen", "", 0, "L", false, 0, "")
	doc.CellFormat(62, 6, "CO₂-Kostenanteil", "", 1, "R", false, 0, "")

	doc.SetXY(14, 66)
	doc.SetFont(fontBold, "", 11)
	doc.CellFormat(88, 6, "Physische Berechnungsgrundlage", "", 0, "L", false, 0, "")
	doc.SetX(108)
	doc.CellFormat(88, 6, "Preisgrundlage", "", 1, "L", false, 0, "")

	left := []string{
		"Menge: " + formatNumber(r.Quantity, 2) + " " + r.Unit,
		"Energiegehalt: " + formatNumber(float64(r.EnergyContent), 2) + " kWh",
		"Emissionsintensität: " + formatNumber(r.CO2PerKWh, 4) + " kg CO₂/kWh",
		"Heizwert Hu: " + formatNumber(r.Factor.CalorificValue, 1) + " MJ/kg",
		"Emissionsfaktor: " + formatNumber(r.Factor.EmissionFactor, 4) + " kg CO₂/MJ",
	}
	if r.Factor.Density > 0 {
		left = append(left, "Dichte: "+formatNumber(r.Factor.Density, 3)+" kg/L")
	}
	left = append(left,
		fmt.Sprintf("Faktor gültig seit: %s", formatDate(r.Factor.ValidFrom)),
		fmt.Sprintf("Quellenstand Faktor: %d", r.Factor.SourceYear),
	)

	priceKind := "Festpreis"
	if r.Price.IsAssumption {
		priceKind = "Annahme"
	}
	right := []string{
		"CO₂-Preis: " + formatNumber(r.CO2Price, 2) + " €/t netto",
		fmt.Sprintf("Preisstand: %d · %s", r.Price.ReferenceYear, priceKind),
	}
	if r.Price.RangeMin > 0 && r.Price.RangeMax > 0 {
		right = append(right, "Offizieller Korridor: "+formatNumber(r.Price.RangeMin, 0)+"–"+formatNumber(r.Price.RangeMax, 0)+" €/t")
	}
	right = append(right,
		"Preisquelle: "+r.Price.Source,
		"Faktorquelle: "+r.Factor.Source,
		"Umsatzsteuer auf Kostenanteil: 19 %",
		fmt.Sprintf("App-Version: %s", r.AppVersion),
	)

	doc.SetFont(fontRegular, "", 8.5)
	writeLines(doc, 14, 73, 88, left)
	writeLines(doc, 108, 73, 88, right)

	doc.SetDrawColor(212, 209, 200)
	doc.Line(14, 130, 196, 130)
	doc.SetXY(14, 133)
	doc.SetFont(fontRegular, "", 7)
	doc.SetTextColor(97, 102, 115)
	auditID := r.AuditHash
	if len(auditID) > 16 {
		auditID = auditID[:16]
	}
	doc.CellFormat(0, 4, "Audit-ID: "+auditID+" · Interne Nachvollziehbarkeit, keine amtliche Zertifizierung.", "", 1, "L", false, 0, "")

	var buf bytes.Buffer
	if err := doc.Output(&buf); err != nil {
		return nil, fmt.Errorf("pdf render: %w", err)
	}
	return buf.Bytes(), nil
}

func writeLines(doc *gofpdf.Fpdf, x, y, width float64, lines []string) {
	doc.SetXY(x, y)
	for _, line := range lines {
		doc.MultiCell(width, 4.7, line, "", "L", false)
		doc.SetX(x)
	}
}

func formatDate(value time.Time) string {
	if value.IsZero() {
		return "nicht angegeben"
	}
	return value.Format("02.01.2006")
}

func formatNumber(value float64, decimals int) string {
	formatted := strconv.FormatFloat(value, 'f', decimals, 64)
	parts := strings.SplitN(formatted, ".", 2)
	integer := parts[0]
	start := 0
	if strings.HasPrefix(integer, "-") {
		start = 1
	}
	for index := len(integer) - 3; index > start; index -= 3 {
		integer = integer[:index] + "." + integer[index:]
	}
	if len(parts) == 1 {
		return integer
	}
	return integer + "," + parts[1]
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
