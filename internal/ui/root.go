package ui

import (
	"emissioncalculator/internal/calculation"
	"emissioncalculator/internal/models"
	"emissioncalculator/internal/pdf"
	"emissioncalculator/internal/validation"
	"image/color"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type designTokens struct {
	spacingSmall  float32
	spacingMedium float32
	spacingLarge  float32
	cardRadius    float32
	sidebarWidth  float32
	contentWidth  float32
	contentHeight float32
}

type palette struct {
	bg            color.Color
	sidebarBg     color.Color
	surface       color.Color
	resultSurface color.Color
	border        color.Color
	textPrimary   color.Color
	textSecondary color.Color
	textMuted     color.Color
	textDisabled  color.Color
	accent        color.Color
	accentHover   color.Color
	accentSoft    color.Color
	error         color.Color
	success       color.Color
}

var ui = designTokens{spacingSmall: 8, spacingMedium: 16, spacingLarge: 24, cardRadius: 12, sidebarWidth: 164, contentWidth: 820, contentHeight: 370}

var appPalette = palette{
	bg:            color.NRGBA{R: 0xf7, G: 0xf6, B: 0xef, A: 255},
	sidebarBg:     color.NRGBA{R: 0x36, G: 0x58, B: 0xd4, A: 255},
	surface:       color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 255},
	resultSurface: color.NRGBA{R: 0xdc, G: 0xff, B: 0x57, A: 255},
	border:        color.NRGBA{R: 0xdf, G: 0xe1, B: 0xdc, A: 255},
	textPrimary:   color.NRGBA{R: 0x14, G: 0x20, B: 0x45, A: 255},
	textSecondary: color.NRGBA{R: 0x6f, G: 0x76, B: 0x87, A: 255},
	textMuted:     color.NRGBA{R: 0x91, G: 0x96, B: 0xa3, A: 255},
	textDisabled:  color.NRGBA{R: 0xaf, G: 0xb8, B: 0xbd, A: 255},
	accent:        color.NRGBA{R: 0x36, G: 0x58, B: 0xd4, A: 255},
	accentHover:   color.NRGBA{R: 0x25, G: 0x45, B: 0xbc, A: 255},
	accentSoft:    color.NRGBA{R: 0xe8, G: 0xee, B: 0xff, A: 255},
	error:         color.NRGBA{R: 0xc7, G: 0x3d, B: 0x46, A: 255},
	success:       color.NRGBA{R: 0x1d, G: 0x76, B: 0x70, A: 255},
}

type calmTheme struct{ fyne.Theme }

func (themeOverride calmTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary, theme.ColorNameFocus, theme.ColorNameSelection:
		return appPalette.accent
	case theme.ColorNameHover:
		return appPalette.accentHover
	case theme.ColorNameBackground:
		return appPalette.bg
	case theme.ColorNameInputBackground, theme.ColorNameButton:
		return appPalette.surface
	case theme.ColorNameDisabledButton:
		return appPalette.accentSoft
	case theme.ColorNameDisabled:
		return appPalette.textDisabled
	case theme.ColorNameForeground:
		return appPalette.textPrimary
	case theme.ColorNamePlaceHolder:
		return appPalette.textMuted
	}
	return themeOverride.Theme.Color(name, variant)
}

func NewRootWindow(app fyne.App) fyne.Window {
	app.Settings().SetTheme(calmTheme{Theme: theme.LightTheme()})
	win := app.NewWindow("Emissionsrechner")
	win.Resize(fyne.NewSize(780, 860))
	win.SetContent(buildView(win, "oil"))
	return win
}

func buildView(win fyne.Window, mode string) fyne.CanvasObject {
	return buildReferenceView(win, mode).content
}

var exportLabel = pdf.ExportLabel
var openExportedFile = openFile

type referenceView struct {
	content         fyne.CanvasObject
	quantityEntry   *widget.Entry
	calculateButton *widget.Button
	printButton     *widget.Button
	status          *canvas.Text
	result          models.CalculationResult
	mode            string
	resultValue     *canvas.Text
	costValue       *canvas.Text
	energyValue     *canvas.Text
	co2Value        *canvas.Text
}

func buildReferenceView(win fyne.Window, mode string) *referenceView {
	view := &referenceView{mode: mode}
	view.quantityEntry = widget.NewEntry()
	view.quantityEntry.SetPlaceHolder("z. B. 12.500 oder 12,5")
	view.status = newStatusText()
	view.status.Hide()

	view.resultValue = canvas.NewText("—", appPalette.textMuted)
	view.resultValue.TextSize = 46
	view.resultValue.TextStyle = fyne.TextStyle{Bold: true}
	view.costValue = referenceDetailValue("—")
	view.energyValue = referenceDetailValue("—")
	view.co2Value = referenceDetailValue("—")

	view.calculateButton = widget.NewButton("Jetzt berechnen  →", view.calculate)
	view.calculateButton.Importance = widget.HighImportance
	view.printButton = widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), view.print)
	view.printButton.Importance = widget.LowImportance
	view.printButton.Disable()
	view.quantityEntry.OnSubmitted = func(string) { view.calculate() }
	view.quantityEntry.OnChanged = func(string) { view.status.Hide() }

	view.content = buildReferenceContent(view, win)
	return view
}

func buildReferenceContent(view *referenceView, win fyne.Window) fyne.CanvasObject {
	logoBackground := canvas.NewRectangle(appPalette.resultSurface)
	logoBackground.CornerRadius = 10
	logoBackground.SetMinSize(fyne.NewSize(42, 42))
	logoText := canvas.NewText("e°", appPalette.textPrimary)
	logoText.TextSize = 17
	logoText.TextStyle = fyne.TextStyle{Bold: true}
	logo := container.NewStack(logoBackground, container.NewCenter(logoText))
	bricketts := widget.NewButton("Briketts", func() { win.SetContent(buildView(win, "briketts")) })
	oil := widget.NewButton("Heizöl", func() { win.SetContent(buildView(win, "oil")) })
	bricketts.Importance = widget.LowImportance
	oil.Importance = widget.LowImportance
	if view.mode == "oil" {
		oil.Importance = widget.HighImportance
	} else {
		bricketts.Importance = widget.HighImportance
	}
	topBackground := canvas.NewRectangle(appPalette.sidebarBg)
	topBar := container.NewStack(topBackground, container.NewPadded(container.NewBorder(nil, nil, logo, container.NewHBox(bricketts, oil), nil)))

	brandDot := canvas.NewText("●", color.NRGBA{R: 0xff, G: 0x78, B: 0x66, A: 255})
	brandDot.TextSize = 15
	brand := canvas.NewText("Emissionsrechner", appPalette.textPrimary)
	brand.TextSize = 16
	brand.TextStyle = fyne.TextStyle{Bold: true}
	readyDot := canvas.NewCircle(appPalette.success)
	readyIndicator := container.NewGridWrap(fyne.NewSize(10, 10), readyDot)
	ready := canvas.NewText("Bereit", appPalette.textSecondary)
	ready.TextSize = 13
	header := container.NewBorder(nil, nil, container.NewHBox(brandDot, brand), container.NewHBox(readyIndicator, ready, view.printButton), nil)

	step := canvas.NewText("01   LIEFERMENGE", appPalette.accent)
	step.TextSize = 12
	step.TextStyle = fyne.TextStyle{Bold: true}
	title := canvas.NewText(titleForMode(view.mode)+".\nEinfach berechnet.", appPalette.textPrimary)
	title.TextSize = 36
	title.TextStyle = fyne.TextStyle{Bold: true}
	description := canvas.NewText("Menge eingeben und CO2-Ausstoß, Energiegehalt sowie\nKostenanteil auf einen Blick erhalten.", appPalette.textSecondary)
	description.TextSize = 14
	quantityLabel := canvas.NewText("MENGE EINGEBEN", appPalette.textSecondary)
	quantityLabel.TextSize = 11
	quantityLabel.TextStyle = fyne.TextStyle{Bold: true}
	unit := canvas.NewText(referenceUnitForMode(view.mode), appPalette.accent)
	unit.TextSize = 15
	unit.TextStyle = fyne.TextStyle{Bold: true}
	inputLine := canvas.NewRectangle(appPalette.textPrimary)
	inputLine.SetMinSize(fyne.NewSize(1, 2))
	input := container.NewBorder(nil, inputLine, nil, unit, view.quantityEntry)

	content := container.NewVBox(
		header,
		verticalGap(ui.spacingMedium),
		widget.NewSeparator(),
		verticalGap(ui.spacingLarge),
		step,
		verticalGap(ui.spacingMedium),
		title,
		verticalGap(ui.spacingSmall),
		description,
		verticalGap(ui.spacingLarge),
		quantityLabel,
		verticalGap(ui.spacingSmall),
		input,
		verticalGap(ui.spacingSmall),
		view.status,
		verticalGap(ui.spacingMedium),
		container.NewHBox(view.calculateButton),
		verticalGap(ui.spacingLarge),
		buildReferenceResultCard(view),
	)
	background := canvas.NewRectangle(appPalette.bg)
	frame := container.NewGridWrap(fyne.NewSize(720, 0), container.NewPadded(content))
	body := container.NewStack(background, container.NewVScroll(container.NewCenter(frame)))
	return container.NewBorder(topBar, nil, nil, nil, body)
}

func buildReferenceResultCard(view *referenceView) fyne.CanvasObject {
	label := canvas.NewText("GESAMTEMISSIONEN", appPalette.textSecondary)
	label.TextSize = 11
	label.TextStyle = fyne.TextStyle{Bold: true}
	live := canvas.NewText("Live-Ergebnis", appPalette.textSecondary)
	live.TextSize = 12
	hint := canvas.NewText("Berechnet für die eingegebene Liefermenge", appPalette.textSecondary)
	hint.TextSize = 13
	metrics := container.NewGridWithColumns(3,
		referenceMetric("CO2-Kostenanteil", view.costValue),
		referenceMetric("Energiegehalt", view.energyValue),
		referenceMetric("CO2 pro kWh", view.co2Value),
	)
	background := canvas.NewRectangle(appPalette.resultSurface)
	background.CornerRadius = ui.cardRadius
	cardContent := container.NewVBox(
		container.NewBorder(nil, nil, label, live, nil),
		verticalGap(ui.spacingLarge),
		view.resultValue,
		verticalGap(ui.spacingSmall),
		hint,
		verticalGap(ui.spacingMedium),
		widget.NewSeparator(),
		verticalGap(ui.spacingMedium),
		metrics,
	)
	return container.NewStack(background, container.NewPadded(cardContent))
}

func (view *referenceView) calculate() {
	input, err := validation.ParseQuantity(view.quantityEntry.Text)
	if err != nil {
		setStatus(view.status, err.Error(), appPalette.error)
		view.printButton.Disable()
		return
	}
	if view.mode == "oil" {
		view.result = calculation.CalculateOil(input)
	} else {
		view.result = calculation.CalculateBriquettes(input)
	}
	view.resultValue.Text = view.result.Emissions
	view.resultValue.Color = appPalette.textPrimary
	view.resultValue.Refresh()
	view.costValue.Text = view.result.EmissionCost
	view.costValue.Refresh()
	view.energyValue.Text = view.result.EnergyContent
	view.energyValue.Refresh()
	view.co2Value.Text = view.result.CO2PerKWh
	view.co2Value.Refresh()
	view.printButton.Enable()
	setStatus(view.status, "Ergebnis berechnet.", appPalette.success)
}

func (view *referenceView) print() {
	if !view.result.Valid {
		setStatus(view.status, "Bitte zuerst eine gültige Berechnung durchführen.", appPalette.error)
		return
	}
	path, err := exportLabel(view.result.Emissions, view.result.EmissionCost, view.result.EnergyContent, view.result.CO2PerKWh)
	if err != nil {
		setStatus(view.status, "Fehler beim PDF-Export: "+err.Error(), appPalette.error)
		return
	}
	if err := openExportedFile(path); err != nil {
		setStatus(view.status, "PDF erstellt: "+path, appPalette.success)
		return
	}
	setStatus(view.status, "PDF wurde erstellt.", appPalette.success)
}

func referenceUnitForMode(mode string) string {
	if mode == "oil" {
		return "Liter"
	}
	return "Tonnen"
}

func referenceDetailValue(value string) *canvas.Text {
	text := canvas.NewText(value, appPalette.textPrimary)
	text.TextSize = 13
	text.TextStyle = fyne.TextStyle{Bold: true}
	return text
}

func referenceMetric(label string, value *canvas.Text) fyne.CanvasObject {
	metricLabel := canvas.NewText(label, appPalette.textSecondary)
	metricLabel.TextSize = 10
	return container.NewVBox(metricLabel, verticalGap(ui.spacingSmall), value)
}

func buildSidebar(win fyne.Window, currentMode string) fyne.CanvasObject {
	brand := canvas.NewText("Emissionsrechner", color.White)
	brand.TextSize = 18
	brand.TextStyle = fyne.TextStyle{Bold: true}
	brandCaption := canvas.NewText("CO2-Kennzahlen", color.NRGBA{R: 0xb8, G: 0xc4, B: 0xc9, A: 255})
	brandCaption.TextSize = 12

	buttons := []*widget.Button{
		widget.NewButton("Briketts", func() { win.SetContent(buildView(win, "briketts")) }),
		widget.NewButton("Heizöl", func() { win.SetContent(buildView(win, "oil")) }),
	}

	for i, btn := range buttons {
		btn.Importance = widget.MediumImportance
		btn.Resize(fyne.NewSize(ui.sidebarWidth-ui.spacingLarge, 36))
		if (currentMode == "briketts" && i == 0) || (currentMode == "oil" && i == 1) {
			btn.Importance = widget.HighImportance
		}
	}

	footer := canvas.NewText("© Lycr.eu", color.NRGBA{R: 0x95, G: 0xa2, B: 0xa8, A: 255})
	footer.TextSize = 11
	header := container.NewVBox(brand, brandCaption)
	navigation := container.NewVBox(buttons[0], verticalGap(ui.spacingSmall), buttons[1])
	inner := container.NewBorder(header, footer, nil, nil, navigation)

	sidebarBg := canvas.NewRectangle(appPalette.sidebarBg)
	sidebarBg.SetMinSize(fyne.NewSize(ui.sidebarWidth, 0))
	return container.NewStack(sidebarBg, container.NewPadded(inner))
}

func buildCalculationView(win fyne.Window, mode string) fyne.CanvasObject {
	result := models.CalculationResult{}

	labelTitle := widget.NewLabelWithStyle("Emissionsberechnung für "+titleForMode(mode), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	labelTitle.Importance = widget.HighImportance

	subtitle := widget.NewLabel("Liefermenge erfassen und die zugehörigen CO2-Kennzahlen berechnen.")
	subtitle.Importance = widget.LowImportance

	quantityEntry := widget.NewEntry()
	quantityEntry.PlaceHolder = "z. B. 12,5"
	quantityEntry.SetPlaceHolder("z. B. 12,5")

	statusLabel := newStatusText()
	statusLabel.Hide()

	unitLabel := widget.NewLabel("t")
	if mode == "oil" {
		unitLabel.SetText("l")
	}
	unitLabel.Alignment = fyne.TextAlignCenter
	unitLabel.TextStyle = fyne.TextStyle{Bold: true}

	inputField := quantityInput(quantityEntry, unitLabel)
	inputPanel := newPanel(
		"Liefermenge",
		"Die Einheit ist bereits mit der gewählten Brennstoffart verknüpft.",
		container.NewVBox(inputField.content, statusLabel),
		appPalette.surface,
	)

	resultValue := canvas.NewText("—", appPalette.textMuted)
	resultValue.TextSize = 30
	resultValue.TextStyle = fyne.TextStyle{Bold: true}
	resultMeta := widget.NewLabel("Nach der Berechnung wird die Emissionsmenge hier angezeigt.")
	resultMeta.Importance = widget.LowImportance

	costValue := widget.NewLabel("--")
	energyValue := widget.NewLabel("--")
	co2Value := widget.NewLabel("--")
	resultPanel := newPanel(
		"Gesamtemissionen",
		"Ergebnis der eingegebenen Liefermenge.",
		container.NewVBox(
			resultValue,
			resultMeta,
			widget.NewSeparator(),
			resultDetail("CO2-Kostenanteil", costValue),
			resultDetail("Energiegehalt", energyValue),
			resultDetail("CO2 pro kWh", co2Value),
		),
		appPalette.resultSurface,
	)

	primaryButton := widget.NewButton("Berechnen", func() {})
	primaryButton.Importance = widget.HighImportance
	primaryButton.Resize(fyne.NewSize(144, 36))

	printButton := widget.NewButtonWithIcon("Drucken", theme.DocumentSaveIcon(), func() {
		if !result.Valid {
			setStatus(statusLabel, "Bitte zuerst eine gültige Berechnung durchführen.", appPalette.error)
			return
		}

		path, err := pdf.ExportLabel(result.Emissions, result.EmissionCost, result.EnergyContent, result.CO2PerKWh)
		if err != nil {
			setStatus(statusLabel, "Fehler beim PDF-Export: "+err.Error(), appPalette.error)
			return
		}
		if err := openFile(path); err != nil {
			setStatus(statusLabel, "PDF erstellt: "+path, appPalette.success)
			return
		}
		setStatus(statusLabel, "PDF wurde erstellt.", appPalette.success)
	})
	printButton.Importance = widget.LowImportance
	printButton.Disable()

	calculateForMode := func() {
		input, err := validation.ParseQuantity(quantityEntry.Text)
		if err != nil {
			inputField.SetError(true)
			setStatus(statusLabel, err.Error(), appPalette.error)
			printButton.Disable()
			return
		}
		inputField.SetError(false)
		if mode == "oil" {
			result = calculation.CalculateOil(input)
		} else {
			result = calculation.CalculateBriquettes(input)
		}
		resultValue.Text = result.Emissions
		resultValue.Color = appPalette.textPrimary
		resultValue.Refresh()
		resultMeta.SetText("Berechnete Emissionen der eingegebenen Liefermenge.")
		costValue.SetText(result.EmissionCost)
		energyValue.SetText(result.EnergyContent)
		co2Value.SetText(result.CO2PerKWh)
		printButton.Enable()
		setStatus(statusLabel, "Ergebnis berechnet.", appPalette.success)
	}

	primaryButton.OnTapped = func() { calculateForMode() }
	quantityEntry.OnSubmitted = func(string) { calculateForMode() }
	quantityEntry.OnChanged = func(string) {
		inputField.SetError(false)
		statusLabel.Hide()
	}

	inputColumn := container.NewVBox(inputPanel, verticalGap(ui.spacingMedium), container.NewHBox(primaryButton))
	resultColumn := container.NewVBox(resultPanel, verticalGap(ui.spacingMedium), container.NewHBox(printButton))
	contentGrid := container.NewGridWithColumns(2, inputColumn, resultColumn)

	mainContent := container.NewVBox(
		labelTitle,
		verticalGap(ui.spacingSmall),
		subtitle,
		verticalGap(ui.spacingMedium),
		widget.NewSeparator(),
		verticalGap(ui.spacingMedium),
		contentGrid,
	)

	bg := canvas.NewRectangle(appPalette.bg)
	contentFrame := container.NewGridWrap(fyne.NewSize(ui.contentWidth, ui.contentHeight), container.NewPadded(mainContent))
	return container.NewStack(bg, container.NewCenter(contentFrame))
}

func titleForMode(mode string) string {
	if mode == "oil" {
		return "Heizöl"
	}
	return "Briketts"
}

type quantityControl struct {
	content fyne.CanvasObject
	border  *canvas.Rectangle
}

func (control *quantityControl) SetError(hasError bool) {
	if hasError {
		control.border.StrokeColor = appPalette.error
		control.border.StrokeWidth = 1.5
	} else {
		control.border.StrokeColor = appPalette.border
		control.border.StrokeWidth = 1
	}
	control.border.Refresh()
}

func newPanel(title, supporting string, content fyne.CanvasObject, background color.Color) fyne.CanvasObject {
	panelTitle := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	panelSupporting := widget.NewLabel(supporting)
	panelSupporting.Importance = widget.LowImportance

	bg := canvas.NewRectangle(background)
	bg.CornerRadius = ui.cardRadius
	bg.StrokeColor = appPalette.border
	bg.StrokeWidth = 1
	panelContent := container.NewVBox(
		panelTitle,
		verticalGap(ui.spacingSmall/2),
		panelSupporting,
		verticalGap(ui.spacingMedium),
		content,
	)
	return container.NewStack(bg, container.NewPadded(panelContent))
}

func quantityInput(entry *widget.Entry, unit *widget.Label) *quantityControl {
	unitBg := canvas.NewRectangle(appPalette.accentSoft)
	unitBg.CornerRadius = ui.cardRadius - 2
	unitBg.SetMinSize(fyne.NewSize(52, 36))
	unitBox := container.NewStack(unitBg, container.NewCenter(unit))

	bg := canvas.NewRectangle(appPalette.surface)
	bg.CornerRadius = ui.cardRadius
	bg.StrokeColor = appPalette.border
	bg.StrokeWidth = 1
	field := container.NewBorder(nil, nil, nil, unitBox, entry)
	return &quantityControl{content: container.NewStack(bg, container.NewPadded(field)), border: bg}
}

func resultDetail(label string, value *widget.Label) fyne.CanvasObject {
	labelText := canvas.NewText(label, appPalette.textSecondary)
	labelText.TextSize = 12
	value.TextStyle = fyne.TextStyle{Bold: true}
	return container.NewBorder(nil, nil, labelText, value)
}

func verticalGap(height float32) fyne.CanvasObject {
	gap := canvas.NewRectangle(color.Transparent)
	gap.SetMinSize(fyne.NewSize(1, height))
	return gap
}

func newStatusText() *canvas.Text {
	status := canvas.NewText("", appPalette.textMuted)
	status.TextSize = 12
	return status
}

func setStatus(status *canvas.Text, message string, colorValue color.Color) {
	status.Text = message
	status.Color = colorValue
	status.TextStyle = fyne.TextStyle{Bold: true}
	status.Show()
	status.Refresh()
}

func openFile(path string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", path).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", path).Start()
	default:
		return exec.Command("xdg-open", path).Start()
	}
}
