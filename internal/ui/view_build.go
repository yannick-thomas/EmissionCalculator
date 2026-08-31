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
	view.quantityEntry.SetPlaceHolder("z. B. 12.500 oder 12,5")
	view.status = newStatusText()

	view.resultValue = canvasText("—", 72, appPalette.textPrimary, true)
	view.resultUnit = canvasText("", 23, appPalette.textPrimary, true)
	view.resultHint = canvasText("Noch keine Berechnung", 16, appPalette.textSecondary, false)
	view.costValue = detailValue("—")
	view.energyValue = detailValue("—")
	view.co2Value = detailValue("—")
	view.resultBasis = canvasText("", 14, color.Transparent, false)

	view.calculateButton = newActionButton("Jetzt berechnen", view.calculate)
	view.printButton = newCircleIconButton(printIconResource(), view.print)
	view.printButton.Disable()
	view.saveButton = newCircleIconButton(saveIconResource(), view.saveAs)
	view.saveButton.Disable()
	view.scenarioButton = newCircleIconButton(scenarioIconResource(), view.showScenarios)
	view.scenarioButton.Disable()

	view.content = buildAppShell(view)
	view.quantityEntry.OnSubmitted = func(string) { view.calculate() }
	view.quantityEntry.OnChanged = func(string) {
		if _, err := validation.ParseQuantity(view.quantityEntry.Text); err != nil {
			view.quantityControl.SetError(true)
			setStatus(view.status, err.Error(), appPalette.error)
			view.markResultStale(resultStateInvalid, "Eingabe prüfen", appPalette.error)
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
	form := container.NewGridWrap(fyne.NewSize(ui.formColWidth, ui.colHeight), buildForm(view))
	result := container.NewGridWrap(fyne.NewSize(ui.resultColWidth, ui.colHeight), buildResultCard(view))
	columns := container.NewHBox(form, horizontalGap(20), result)
	page := container.NewVBox(
		header,
		verticalGap(36),
		columns,
		verticalGap(36),
	)
	scroll := container.NewVScroll(container.NewCenter(page))
	view.scroll = scroll
	return container.NewStack(canvas.NewRectangle(appPalette.background), scroll)
}

func buildHeader(view *referenceView) fyne.CanvasObject {
	brandBars := container.NewHBox(
		coloredBar(appPalette.accent, 8, 19),
		coloredBar(appPalette.coral, 8, 12),
	)
	brand := canvasText("Emissionsrechner", 18, appPalette.textPrimary, true)
	brandGroup := container.NewHBox(brandBars, brand)

	view.headerStatusDot = canvas.NewCircle(appPalette.success)
	statusHalo := canvas.NewCircle(color.NRGBA{R: 0x83, G: 0xa8, B: 0x31, A: 30})
	readyIndicator := container.NewGridWrap(fyne.NewSize(24, 24), container.NewStack(statusHalo, container.NewPadded(view.headerStatusDot)))
	view.headerStatus = canvasText("Bereit", 16, appPalette.textSecondary, true)
	statusGroup := container.NewGridWrap(fyne.NewSize(190, 54), container.NewCenter(container.NewHBox(readyIndicator, view.headerStatus)))
	actions := container.NewHBox(
		statusGroup,
		horizontalGap(8),
		headerAction(view.scenarioButton, "Details"),
		horizontalGap(8),
		headerAction(view.saveButton, "PDF speichern"),
		horizontalGap(8),
		headerAction(view.printButton, "Drucken"),
	)
	headerContent := container.NewBorder(nil, nil, container.NewCenter(brandGroup), container.NewCenter(actions), nil)
	headerFrame := container.NewGridWrap(fyne.NewSize(ui.contentWidth, 112), headerContent)

	separator := canvas.NewRectangle(appPalette.border)
	separator.SetMinSize(fyne.NewSize(1, 1))
	return container.NewVBox(headerFrame, separator)
}

func buildForm(view *referenceView) fyne.CanvasObject {
	stepCircle := canvas.NewCircle(color.Transparent)
	stepCircle.StrokeColor = appPalette.accent
	stepCircle.StrokeWidth = 1
	stepNumber := canvasText("01", 12, appPalette.accent, true)
	step := container.NewHBox(
		container.NewGridWrap(fyne.NewSize(42, 42), container.NewStack(stepCircle, container.NewCenter(stepNumber))),
		horizontalGap(10),
		canvasText("L I E F E R M E N G E", 13, appPalette.accent, true),
	)

	title := container.NewVBox(
		canvasText(titleForMode(view.mode)+".", 36, appPalette.textPrimary, true),
		canvasText("Einfach berechnet.", 36, appPalette.textPrimary, true),
	)
	description := container.NewVBox(
		canvasText("Menge eingeben und CO₂-Ausstoß, Energiegehalt sowie", 15, appPalette.textSecondary, false),
		canvasText("Kostenanteil auf einen Blick erhalten.", 15, appPalette.textSecondary, false),
	)

	view.quantityControl = newQuantityControl(view.quantityEntry, unitForMode(view.mode))
	statusFrame := container.NewGridWrap(fyne.NewSize(ui.formColWidth, 20), view.status)
	return container.NewVBox(
		step,
		verticalGap(30),
		title,
		verticalGap(20),
		description,
		verticalGap(38),
		canvasText("Menge eingeben", 14, appPalette.textPrimary, true),
		verticalGap(12),
		view.quantityControl.content,
		verticalGap(8),
		canvasText("Beispiel: 1.250 oder 1250,5", 14, appPalette.textSecondary, false),
		verticalGap(8),
		statusFrame,
		verticalGap(16),
		container.NewHBox(view.calculateButton),
	)
}

func buildResultCard(view *referenceView) fyne.CanvasObject {
	label := canvasText("Gesamtemissionen", 15, appPalette.textSecondary, true)
	view.resultBadge = canvasText("", 14, color.Transparent, true)
	view.resultBadgeBackground = canvas.NewRectangle(color.Transparent)
	view.resultBadgeBackground.CornerRadius = 22
	view.resultBadgeBackground.StrokeColor = color.Transparent
	view.resultBadgeBackground.StrokeWidth = 1
	badge := container.NewGridWrap(
		fyne.NewSize(136, 44),
		container.NewStack(view.resultBadgeBackground, container.NewCenter(view.resultBadge)),
	)
	metrics := container.NewGridWithColumns(3,
		resultMetric("CO₂-Kosten (brutto)", view.costValue),
		resultMetric("Energiegehalt", view.energyValue),
		resultMetric("Emissionsintensität", view.co2Value),
	)
	separator := canvas.NewRectangle(color.NRGBA{R: 0x18, G: 0x21, B: 0x3d, A: 36})
	separator.SetMinSize(fyne.NewSize(1, 1))
	view.resultValueRow = container.NewVBox(view.resultValue, view.resultUnit)
	cardContent := container.NewVBox(
		container.NewBorder(nil, nil, label, badge, nil),
		verticalGap(24),
		view.resultValueRow,
		verticalGap(10),
		view.resultHint,
		verticalGap(20),
		separator,
		verticalGap(16),
		metrics,
		verticalGap(18),
		canvasText("Berechnungsgrundlage", 15, appPalette.textPrimary, true),
		verticalGap(8),
		view.resultBasis,
	)

	view.resultBackground = canvas.NewImageFromResource(resultPanelResource(false))
	view.resultBackground.FillMode = canvas.ImageFillStretch
	contentInset := container.NewBorder(verticalGap(32), verticalGap(30), horizontalGap(36), horizontalGap(36), cardContent)
	card := container.NewStack(view.resultBackground, contentInset)
	shadowImage := canvas.NewImageFromResource(resultShadowResource())
	shadowImage.FillMode = canvas.ImageFillStretch
	return container.NewStack(shadowImage, card)
}

func headerAction(button *circleIconButton, label string) fyne.CanvasObject {
	return container.NewGridWrap(fyne.NewSize(126, 54), container.NewBorder(nil, nil, button, nil, container.NewCenter(canvasText(label, 13, appPalette.textPrimary, true))))
}
