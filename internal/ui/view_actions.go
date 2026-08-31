package ui

import (
	"emissioncalculator/internal/calculation"
	"emissioncalculator/internal/models"
	"emissioncalculator/internal/validation"
	"image/color"

	"fyne.io/fyne/v2/canvas"
)

func (view *referenceView) calculate() {
	input, err := validation.ParseQuantity(view.quantityEntry.Text)
	if err != nil {
		view.clearResult()
		view.quantityControl.SetError(true)
		setStatus(view.status, err.Error(), appPalette.error)
		view.setHeaderStatus("Eingabe prüfen", appPalette.error)
		return
	}

	view.quantityControl.SetError(false)
	cfg := calculation.DefaultConfig()
	if view.configProvider != nil {
		cfg = view.configProvider()
	}
	if view.mode == modeOil {
		view.result = calculation.CalculateOil(input, cfg)
	} else {
		view.result = calculation.CalculateBriquettes(input, cfg)
	}
	emissionsFormatted := formatFloat(float64(view.result.Emissions), 2)
	view.resultValue.Text = emissionsFormatted
	view.resultValue.TextSize = resultTextSize(emissionsFormatted)
	view.resultUnit.Text = "kg CO₂"
	view.resultHint.Text = "Berechnet für " + formatQuantityDisplay(input) + " " + unitForMode(view.mode) + " " + titleForMode(view.mode)
	view.costValue.Text = formatFloat(view.result.EmissionCost, 2) + " € brutto"
	view.energyValue.Text = formatFloat(float64(view.result.EnergyContent), 2) + " kWh"
	view.co2Value.Text = formatFloat(view.result.CO2PerKWh, 4) + " kg CO₂/kWh"
	fitDetailText(view.costValue)
	fitDetailText(view.energyValue)
	fitDetailText(view.co2Value)
	view.resultBackground.Resource = resultPanelResource(true)
	view.printButton.Enable()
	view.saveButton.Enable()
	view.scenarioButton.Enable()
	view.state = resultStateCalculated
	view.resultBadge.Text = "Berechnet"
	view.resultBadge.Color = appPalette.textSecondary
	view.resultBadgeBackground.StrokeColor = color.NRGBA{R: 0x62, G: 0x77, B: 0x2e, A: 90}
	view.resultBasis.SetRecord(view.result)
	refreshTexts(view.resultValue, view.resultUnit, view.resultHint, view.costValue, view.energyValue, view.co2Value, view.resultBadge)
	view.resultValueRow.Refresh()
	view.resultBackground.Refresh()
	canvas.Refresh(view.resultBadgeBackground)
	setStatus(view.status, " ", appPalette.textSecondary)
	view.setHeaderStatus("Berechnet", appPalette.success)
}

func (view *referenceView) refreshForSettingsChange() {
	if view.result.Valid {
		view.markResultStale(resultStateStale, "Eingabe geändert – neu berechnen", appPalette.accent)
	}
}

func (view *referenceView) markResultStale(state resultState, message string, statusColor color.Color) {
	view.state = state
	view.printButton.Disable()
	view.saveButton.Disable()
	view.scenarioButton.Disable()
	if view.result.Valid {
		view.resultHint.Text = "Eingabe geändert – neu berechnen"
		view.resultBadge.Text = "Ergebnis veraltet"
		view.resultBadge.Color = appPalette.textPrimary
		view.resultBadgeBackground.StrokeColor = appPalette.accent
		refreshTexts(view.resultHint, view.resultBadge)
		canvas.Refresh(view.resultBadgeBackground)
	}
	view.setHeaderStatus(message, statusColor)
}

func (view *referenceView) clearResult() {
	view.result = models.CalculationRecord{}
	view.state = resultStateInvalid
	view.resultValue.Text = "—"
	view.resultValue.TextSize = 72
	view.resultUnit.Text = ""
	view.resultHint.Text = "Noch keine Berechnung"
	view.costValue.Text = "—"
	view.energyValue.Text = "—"
	view.co2Value.Text = "—"
	fitDetailText(view.costValue)
	fitDetailText(view.energyValue)
	fitDetailText(view.co2Value)
	view.resultBackground.Resource = resultPanelResource(false)
	view.printButton.Disable()
	view.saveButton.Disable()
	view.scenarioButton.Disable()
	view.resultBadge.Color = color.Transparent
	view.resultBadgeBackground.StrokeColor = color.Transparent
	view.resultBasis.Clear()
	refreshTexts(view.resultValue, view.resultUnit, view.resultHint, view.costValue, view.energyValue, view.co2Value, view.resultBadge)
	view.resultValueRow.Refresh()
	view.resultBackground.Refresh()
	canvas.Refresh(view.resultBadgeBackground)
}

func (view *referenceView) hasCurrentResult() bool {
	return view.result.Valid && (view.state == resultStateCalculated || view.state == resultStatePDFCreated)
}

func (view *referenceView) print() {
	if !view.hasCurrentResult() {
		setStatus(view.status, "Bitte zuerst eine gültige Berechnung durchführen.", appPalette.error)
		view.setHeaderStatus("Eingabe erforderlich", appPalette.error)
		return
	}
	view.setHeaderStatus("PDF wird erstellt", appPalette.accent)
	path, err := exportLabel(view.result)
	if err != nil {
		setStatus(view.status, "Fehler beim PDF-Export: "+err.Error(), appPalette.error)
		view.setHeaderStatus("PDF-Fehler", appPalette.error)
		return
	}
	view.state = resultStatePDFCreated
	if err := openExportedFile(path); err != nil {
		setStatus(view.status, "PDF erstellt, konnte aber nicht geöffnet werden: "+err.Error(), appPalette.error)
		view.setHeaderStatus("PDF erstellt", appPalette.success)
		return
	}
	setStatus(view.status, "PDF wurde erstellt.", appPalette.success)
	view.setHeaderStatus("PDF erstellt", appPalette.success)
}

func (view *referenceView) setHeaderStatus(message string, statusColor color.Color) {
	view.headerStatus.Text = message
	view.headerStatusDot.FillColor = statusColor
	view.headerStatus.Refresh()
	view.headerStatusDot.Refresh()
}

func (view *referenceView) saveAs() {
	if !view.hasCurrentResult() {
		setStatus(view.status, "Bitte zuerst eine gültige Berechnung durchführen.", appPalette.error)
		view.setHeaderStatus("Eingabe erforderlich", appPalette.error)
		return
	}
	view.setHeaderStatus("Speichern...", appPalette.accent)
	saveLabelAs(view.window, view.result, func(path string, err error) {
		if err != nil {
			setStatus(view.status, "Fehler beim Speichern: "+err.Error(), appPalette.error)
			view.setHeaderStatus("Speicherfehler", appPalette.error)
			return
		}
		if path == "" {
			view.setHeaderStatus("Berechnet", appPalette.success)
			return
		}
		setStatus(view.status, "PDF gespeichert: "+path, appPalette.success)
		view.state = resultStatePDFCreated
		view.setHeaderStatus("PDF erstellt", appPalette.success)
	})
}

func (view *referenceView) showScenarios() {
	if !view.hasCurrentResult() {
		setStatus(view.status, "Bitte zuerst eine gültige Berechnung durchführen.", appPalette.error)
		view.setHeaderStatus("Eingabe erforderlich", appPalette.error)
		return
	}
	showScenarioDialog(view)
}
