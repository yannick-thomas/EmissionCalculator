package ui

import (
	"emissioncalculator/internal/models"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// calculationBasis renders each provenance value in its own layout cell. This
// intentionally avoids canvas.Text with embedded newlines, which Fyne renders as
// a single overflowing line.
type calculationBasis struct {
	content          *fyne.Container
	emissionFactor   *canvas.Text
	factorValidFrom  *canvas.Text
	factorSourceYear *canvas.Text
	price            *canvas.Text
	priceReference   *canvas.Text
	calculationYear  *canvas.Text
}

func newCalculationBasis() *calculationBasis {
	basis := &calculationBasis{
		emissionFactor:   basisValue(),
		factorValidFrom:  basisValue(),
		factorSourceYear: basisValue(),
		price:            basisValue(),
		priceReference:   basisValue(),
		calculationYear:  basisValue(),
	}
	basis.content = container.NewGridWithColumns(2,
		basisLabel("Emissionsfaktor"), basis.emissionFactor,
		basisLabel("Faktor gültig seit"), basis.factorValidFrom,
		basisLabel("Quellenstand Faktor"), basis.factorSourceYear,
		basisLabel("CO₂-Preis (netto)"), basis.price,
		basisLabel("Preisstand"), basis.priceReference,
		basisLabel("Berechnungsjahr"), basis.calculationYear,
	)
	return basis
}

func basisLabel(text string) *canvas.Text {
	return canvasText(text, 12, appPalette.textSecondary, false)
}

func basisValue() *canvas.Text {
	return canvasText("—", 12, appPalette.textPrimary, true)
}

func (basis *calculationBasis) SetRecord(record models.CalculationRecord) {
	basis.emissionFactor.Text = formatFloat(record.Factor.EmissionFactor, 4) + " kg CO₂/MJ"
	basis.factorValidFrom.Text = record.Factor.ValidFrom.Format("02.01.2006")
	basis.factorSourceYear.Text = strconv.Itoa(record.Factor.SourceYear)
	basis.price.Text = formatFloat(record.CO2Price, 2) + " €/t"
	priceKind := "Festpreis"
	if record.Price.IsAssumption {
		priceKind = "Annahme"
	}
	basis.priceReference.Text = strconv.Itoa(record.Price.ReferenceYear) + " · " + priceKind
	basis.calculationYear.Text = strconv.Itoa(record.CalculationYear)
	refreshTexts(
		basis.emissionFactor,
		basis.factorValidFrom,
		basis.factorSourceYear,
		basis.price,
		basis.priceReference,
		basis.calculationYear,
	)
	basis.content.Refresh()
}

func (basis *calculationBasis) Clear() {
	for _, value := range []*canvas.Text{
		basis.emissionFactor,
		basis.factorValidFrom,
		basis.factorSourceYear,
		basis.price,
		basis.priceReference,
		basis.calculationYear,
	} {
		value.Text = "—"
		value.Refresh()
	}
	basis.content.Refresh()
}

func (basis *calculationBasis) values() []string {
	return []string{
		basis.emissionFactor.Text,
		basis.factorValidFrom.Text,
		basis.factorSourceYear.Text,
		basis.price.Text,
		basis.priceReference.Text,
		basis.calculationYear.Text,
	}
}
