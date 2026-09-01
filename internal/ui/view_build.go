package ui

import (
	"emissioncalculator/internal/calculation"
	"emissioncalculator/internal/validation"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func buildReferenceView(window fyne.Window, mode string) *referenceView {
	return buildReferenceViewWithConfig(window, mode, calculation.DefaultConfig)
}

func buildReferenceViewWithConfig(window fyne.Window, mode string, configProvider func() calculation.Config) *referenceView {
	view := &referenceView{mode: mode, window: window, configProvider: configProvider}
	view.quantityEntry = newFocusEntry()
	view.quantityEntry.SetPlaceHolder("z. B. 12500 oder 12,5")
	view.status = newStatusText()

	view.resultValue = canvasText("—", 60, appPalette.textPrimary, true)
	view.resultUnit = canvasText("", 20, appPalette.textPrimary, true)
	view.resultHint = canvasText("Ergebnis erscheint nach der Berechnung", 15, appPalette.textSecondary, false)
	view.costValue = detailValue("—")
	view.energyValue = detailValue("—")
	view.co2Value = detailValue("—")
	view.resultBasis = newCalculationBasis()

	view.calculateButton = newActionButton("Berechnen", view.calculate)
	view.printButton = newCircleIconButton(printIconResource(), view.print)
	view.printButton.Disable()
	view.saveButton = newCircleIconButton(saveIconResource(), view.saveAs)
	view.saveButton.Disable()
	view.scenarioButton = newCircleIconButton(scenarioIconResource(), view.showScenarios)
	view.scenarioButton.Disable()
	view.traceButton = newCircleIconButton(traceIconResource(), view.showTrace)
	view.traceButton.Disable()

	view.content = buildAppShell(view)
	view.quantityEntry.OnSubmitted = func(string) { view.calculate() }
	view.quantityEntry.OnChanged = func(string) {
		if _, err := validation.ParseQuantity(view.quantityEntry.Text); err != nil {
			// Do not punish incomplete input while the user is still typing. Once a
			// submission failed, keep the validation visible until it is corrected.
			if view.state == resultStateInvalid {
				view.quantityControl.SetError(true)
				setStatus(view.status, err.Error(), appPalette.error)
				view.setHeaderStatus("Eingabe prüfen", appPalette.error)
				return
			}
			view.quantityControl.SetError(false)
			setStatus(view.status, " ", appPalette.textSecondary)
			if view.result.Valid {
				view.markResultStale(resultStateStale, "Eingabe geändert – neu berechnen", appPalette.accent)
				return
			}
			view.setHeaderStatus("Bereit", appPalette.success)
			return
		}
		view.quantityControl.SetError(false)
		setStatus(view.status, " ", appPalette.textSecondary)
		if view.result.Valid {
			view.markResultStale(resultStateStale, "Eingabe geändert – neu berechnen", appPalette.accent)
			return
		}
		view.state = resultStateReady
		view.setHeaderStatus("Bereit", appPalette.success)
	}
	view.setHeaderStatus("Bereit", appPalette.success)
	return view
}

func buildAppShell(view *referenceView) fyne.CanvasObject {
	header := buildHeader(view)
	page := newResponsivePage(header, buildForm(view), buildResultCard(view))
	scroll := container.NewVScroll(page)
	view.scroll = scroll
	return container.NewStack(canvas.NewRectangle(appPalette.background), scroll)
}

func buildHeader(view *referenceView) fyne.CanvasObject {
	brandGroup := canvasText("Emissionsrechner", 16, appPalette.textPrimary, true)

	view.headerStatusDot = canvas.NewCircle(appPalette.success)
	readyIndicator := container.NewGridWrap(fyne.NewSize(14, 14), container.NewPadded(view.headerStatusDot))
	view.headerStatus = canvasText("Bereit", 14, appPalette.textSecondary, false)
	statusGroup := container.NewCenter(container.NewHBox(readyIndicator, view.headerStatus))
	separator := canvas.NewRectangle(appPalette.border)
	separator.SetMinSize(fyne.NewSize(1, 1))
	return container.New(&responsiveHeaderLayout{}, container.NewCenter(brandGroup), statusGroup, container.NewHBox(), separator)
}

func buildForm(view *referenceView) fyne.CanvasObject {
	view.quantityControl = newQuantityControl(view.quantityEntry, unitForMode(view.mode))
	statusFrame := container.NewGridWrap(fyne.NewSize(ui.formColWidth, 20), view.status)
	return container.NewVBox(
		canvasText(titleForMode(view.mode), 13, appPalette.textSecondary, false),
		verticalGap(10),
		canvasText("Liefermenge", 30, appPalette.textPrimary, true),
		verticalGap(8),
		canvasText("Menge eingeben und Ergebnis nachvollziehen.", 14, appPalette.textSecondary, false),
		verticalGap(30),
		canvasText("Menge", 14, appPalette.textPrimary, true),
		verticalGap(12),
		view.quantityControl.content,
		verticalGap(8),
		canvasText("Beispiel: 1.250 oder 1250,5", 14, appPalette.textSecondary, false),
		verticalGap(8),
		statusFrame,
		verticalGap(14),
		container.NewHBox(view.calculateButton),
	)
}

func buildResultCard(view *referenceView) fyne.CanvasObject {
	label := canvasText("Gesamtemissionen", 13, appPalette.textSecondary, true)
	view.resultBadge = canvasText("", 12, color.Transparent, true)
	view.resultBadgeBackground = canvas.NewRectangle(color.Transparent)
	view.resultBadgeBackground.CornerRadius = 16
	view.resultBadgeBackground.StrokeColor = color.Transparent
	view.resultBadgeBackground.StrokeWidth = 1
	badge := container.NewGridWrap(
		fyne.NewSize(124, 32),
		container.NewStack(view.resultBadgeBackground, container.NewCenter(view.resultBadge)),
	)
	metrics := container.NewGridWithColumns(3,
		resultMetric("CO₂-Kosten (brutto)", view.costValue),
		resultMetric("Energiegehalt", view.energyValue),
		resultMetric("Emissionsintensität", view.co2Value),
	)
	separator := canvas.NewRectangle(appPalette.border)
	separator.SetMinSize(fyne.NewSize(1, 1))
	actions := container.NewGridWithColumns(2,
		headerAction(view.traceButton, "Rechenweg"),
		headerAction(view.scenarioButton, "Szenarien"),
		headerAction(view.saveButton, "PDF speichern"),
		headerAction(view.printButton, "Drucken"),
	)
	view.resultValueRow = container.NewVBox(view.resultValue, view.resultUnit)
	cardContent := container.NewVBox(
		container.NewBorder(nil, nil, label, badge, nil),
		verticalGap(20),
		view.resultValueRow,
		verticalGap(10),
		view.resultHint,
		verticalGap(18),
		separator,
		verticalGap(16),
		metrics,
		verticalGap(16),
		canvasText("Berechnungsgrundlage", 15, appPalette.textPrimary, true),
		verticalGap(8),
		view.resultBasis.content,
		verticalGap(14),
		actions,
	)

	view.resultBackground = canvas.NewImageFromResource(resultPanelResource(false))
	view.resultBackground.FillMode = canvas.ImageFillStretch
	contentInset := container.NewBorder(verticalGap(26), verticalGap(26), horizontalGap(30), horizontalGap(30), cardContent)
	return container.NewStack(view.resultBackground, contentInset)
}

func headerAction(button *circleIconButton, label string) fyne.CanvasObject {
	return container.NewGridWrap(fyne.NewSize(145, 38), container.NewBorder(nil, nil, button, nil, container.NewCenter(canvasText(label, 13, appPalette.textPrimary, false))))
}
