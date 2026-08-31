package ui

import (
	"emissioncalculator/internal/calculation"
	"strconv"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

// showScenarioDialog compares the current result's CO2 price against a lower and higher
// assumption for the same fuel, quantity, and factor year — useful for Beratung/Angebotsvergleich.
func showScenarioDialog(view *referenceView) {
	cfg := calculation.DefaultConfig()
	if view.configProvider != nil {
		cfg = view.configProvider()
	}
	fuel := calculation.FuelOil
	if view.mode == modeBriquettes {
		fuel = calculation.FuelBriquettes
	}
	basePrice := view.result.CO2Price
	scenarios := []calculation.PriceScenario{
		{Label: "-20 %", CO2PricePerTonne: basePrice * 0.8},
		{Label: "aktuell", CO2PricePerTonne: basePrice},
		{Label: "+20 %", CO2PricePerTonne: basePrice * 1.2},
	}
	records, err := calculation.CompareCO2Prices(fuel, view.result.Quantity, cfg, scenarios)
	if err != nil {
		setStatus(view.status, "Szenarien konnten nicht berechnet werden: "+err.Error(), appPalette.error)
		view.setHeaderStatus("Fehler", appPalette.error)
		return
	}

	rows := container.NewVBox(
		container.NewGridWithColumns(4,
			canvasText("Szenario", 12, appPalette.textSecondary, true),
			canvasText("Jahr", 12, appPalette.textSecondary, true),
			canvasText("CO₂-Preis", 12, appPalette.textSecondary, true),
			canvasText("Kosten", 12, appPalette.textSecondary, true),
		),
	)
	for _, record := range records {
		rows.Add(container.NewGridWithColumns(4,
			canvasText(record.ScenarioLabel, 13, appPalette.textPrimary, false),
			canvasText(strconv.Itoa(record.FactorYear), 13, appPalette.textPrimary, false),
			canvasText(formatFloat(record.CO2Price, 2)+" €/t", 13, appPalette.textPrimary, false),
			canvasText(formatFloat(record.EmissionCost, 2)+" €", 13, appPalette.textPrimary, false),
		))
	}
	content := container.NewVBox(
		canvasText("Vergleich für "+formatQuantityDisplay(view.result.Quantity)+" "+unitForMode(view.mode)+" "+titleForMode(view.mode), 13, appPalette.textSecondary, false),
		verticalGap(12),
		rows,
	)
	dialog.ShowCustom("Preisszenarien", "Schließen", content, view.window)
}
