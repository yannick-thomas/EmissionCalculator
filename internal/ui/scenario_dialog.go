package ui

import (
	"emissioncalculator/internal/calculation"
	"emissioncalculator/internal/models"
	"emissioncalculator/internal/validation"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// showScenarioDialog compares configurable price assumptions and calculation years.
func showScenarioDialog(view *referenceView) {
	baseConfig := calculation.DefaultConfig()
	if view.configProvider != nil {
		baseConfig = view.configProvider()
	}
	fuel, err := calculation.ParseFuelType(view.mode)
	if err != nil {
		dialog.ShowError(err, view.window)
		return
	}
	yearSelect := widget.NewSelect(yearSelectOptions(), nil)
	yearSelect.SetSelected(strconv.Itoa(view.result.CalculationYear))
	pricesEntry := widget.NewEntry()
	pricesEntry.SetText(defaultScenarioPrices(view.result))
	pricesEntry.SetPlaceHolder("z. B. 55; 60; 65")
	message := widget.NewLabel("")
	rows := container.NewVBox()

	render := func() {
		year, parseYearErr := strconv.Atoi(yearSelect.Selected)
		if parseYearErr != nil {
			message.SetText("Bitte ein gültiges Jahr wählen.")
			return
		}
		prices, parsePricesErr := parseScenarioPrices(pricesEntry.Text)
		if parsePricesErr != nil {
			message.SetText(parsePricesErr.Error())
			return
		}
		cfg := baseConfig
		cfg.CalculationYear, cfg.Year = year, year
		scenarios := make([]calculation.PriceScenario, 0, len(prices))
		for _, price := range prices {
			scenarios = append(scenarios, calculation.PriceScenario{Label: formatFloat(price, 2) + " €/t", CO2PricePerTonne: price})
		}
		records, compareErr := calculation.CompareCO2Prices(fuel, view.result.Quantity, cfg, scenarios)
		if compareErr != nil {
			message.SetText(compareErr.Error())
			return
		}
		message.SetText("")
		rows.RemoveAll()
		rows.Add(container.NewGridWithColumns(5,
			canvasText("Szenario", 12, appPalette.textSecondary, true),
			canvasText("Jahr", 12, appPalette.textSecondary, true),
			canvasText("CO₂-Preis", 12, appPalette.textSecondary, true),
			canvasText("Kosten", 12, appPalette.textSecondary, true),
			canvasText("Differenz", 12, appPalette.textSecondary, true),
		))
		for _, record := range records {
			delta := record.EmissionCost - view.result.EmissionCost
			rows.Add(container.NewGridWithColumns(5,
				canvasText(record.ScenarioLabel, 13, appPalette.textPrimary, false),
				canvasText(strconv.Itoa(record.CalculationYear), 13, appPalette.textPrimary, false),
				canvasText(formatFloat(record.CO2Price, 2)+" €/t", 13, appPalette.textPrimary, false),
				canvasText(formatFloat(record.EmissionCost, 2)+" €", 13, appPalette.textPrimary, false),
				canvasText(signedCurrency(delta), 13, appPalette.textPrimary, false),
			))
		}
		rows.Refresh()
	}
	yearSelect.OnChanged = func(string) { render() }
	updateButton := widget.NewButton("Vergleich aktualisieren", render)
	form := container.NewGridWithColumns(2,
		container.NewVBox(widget.NewLabel("Berechnungsjahr"), yearSelect),
		container.NewVBox(widget.NewLabel("CO₂-Preise in €/t (mit Semikolon trennen)"), pricesEntry),
	)
	render()
	content := container.NewVBox(
		canvasText("Vergleich für "+formatQuantityDisplay(view.result.Quantity)+" "+unitForMode(view.mode)+" "+titleForMode(view.mode), 13, appPalette.textSecondary, false),
		verticalGap(10), form, updateButton, message, verticalGap(8), container.NewVScroll(rows),
	)
	content.Resize(fyne.NewSize(820, 500))
	dialog.ShowCustom("Szenariovergleich", "Schließen", content, view.window)
}

func defaultScenarioPrices(record models.CalculationRecord) string {
	minimum, maximum := record.Price.RangeMin, record.Price.RangeMax
	if minimum <= 0 || maximum <= 0 {
		minimum, maximum = record.CO2Price*0.8, record.CO2Price*1.2
	}
	return fmt.Sprintf("%s; %s; %s", formatFloat(minimum, 2), formatFloat(record.CO2Price, 2), formatFloat(maximum, 2))
}

func parseScenarioPrices(input string) ([]float64, error) {
	parts := strings.Split(input, ";")
	if len(parts) < 2 || len(parts) > 8 {
		return nil, fmt.Errorf("Bitte 2 bis 8 Preise mit Semikolon trennen.")
	}
	prices := make([]float64, 0, len(parts))
	for _, part := range parts {
		price, err := validation.ParseQuantity(part)
		if err != nil {
			return nil, fmt.Errorf("Ungültiger CO₂-Preis %q: %w", strings.TrimSpace(part), err)
		}
		prices = append(prices, price)
	}
	return prices, nil
}

func signedCurrency(value float64) string {
	if value > 0 {
		return "+" + formatFloat(value, 2) + " €"
	}
	return formatFloat(value, 2) + " €"
}
