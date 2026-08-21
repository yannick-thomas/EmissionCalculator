package ui

import (
	"emissioncalculator/internal/calculation"
	"emissioncalculator/internal/models"
	"emissioncalculator/internal/pdf"
	"emissioncalculator/internal/validation"
	"fmt"
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
	spacingTiny   float32
	spacingSmall  float32
	spacingMedium float32
	spacingLarge  float32
	cardRadius    float32
	railWidth     float32
}

type palette struct {
	background     color.Color
	railBackground color.Color
	surface        color.Color
	resultSurface  color.Color
	resultShadow   color.Color
	border         color.Color
	textPrimary    color.Color
	textSecondary  color.Color
	textMuted      color.Color
	textDisabled   color.Color
	accent         color.Color
	accentHover    color.Color
	accentSoft     color.Color
	coral          color.Color
	error          color.Color
	success        color.Color
	white          color.Color
	whiteMuted     color.Color
}

var ui = designTokens{
	spacingTiny:   4,
	spacingSmall:  8,
	spacingMedium: 16,
	spacingLarge:  24,
	cardRadius:    30,
	railWidth:     94,
}

var appPalette = palette{
	background:     color.NRGBA{R: 0xf4, G: 0xf1, B: 0xe8, A: 255},
	railBackground: color.NRGBA{R: 0x31, G: 0x55, B: 0xd7, A: 255},
	surface:        color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 255},
	resultSurface:  color.NRGBA{R: 0xdf, G: 0xfd, B: 0x72, A: 255},
	resultShadow:   color.NRGBA{R: 0xdc, G: 0xd7, B: 0xca, A: 255},
	border:         color.NRGBA{R: 0xd9, G: 0xd5, B: 0xcb, A: 255},
	textPrimary:    color.NRGBA{R: 0x18, G: 0x21, B: 0x3d, A: 255},
	textSecondary:  color.NRGBA{R: 0x68, G: 0x6d, B: 0x7d, A: 255},
	textMuted:      color.NRGBA{R: 0x8b, G: 0x8f, B: 0x99, A: 255},
	textDisabled:   color.NRGBA{R: 0xab, G: 0xae, B: 0xb5, A: 255},
	accent:         color.NRGBA{R: 0x31, G: 0x55, B: 0xd7, A: 255},
	accentHover:    color.NRGBA{R: 0x24, G: 0x45, B: 0xbd, A: 255},
	accentSoft:     color.NRGBA{R: 0xe7, G: 0xec, B: 0xff, A: 255},
	coral:          color.NRGBA{R: 0xff, G: 0x77, B: 0x5f, A: 255},
	error:          color.NRGBA{R: 0xc7, G: 0x3d, B: 0x46, A: 255},
	success:        color.NRGBA{R: 0x3d, G: 0x7e, B: 0x52, A: 255},
	white:          color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 255},
	whiteMuted:     color.NRGBA{R: 0xd6, G: 0xdd, B: 0xff, A: 255},
}

type emissionTheme struct{ fyne.Theme }

func (themeOverride emissionTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary, theme.ColorNameFocus, theme.ColorNameSelection:
		return appPalette.accent
	case theme.ColorNameHover:
		return appPalette.accentHover
	case theme.ColorNameBackground:
		return appPalette.background
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
	app.Settings().SetTheme(emissionTheme{Theme: theme.LightTheme()})
	win := app.NewWindow("Emissionsrechner")
	win.Resize(fyne.NewSize(1080, 700))
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
	resultHint      *canvas.Text
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

	view.resultValue = canvas.NewText("—", appPalette.textPrimary)
	view.resultValue.TextSize = 43
	view.resultValue.TextStyle = fyne.TextStyle{Bold: true}
	view.resultHint = canvas.NewText("Bereit für deine Eingabe", appPalette.textSecondary)
	view.resultHint.TextSize = 12
	view.costValue = detailValue("—")
	view.energyValue = detailValue("—")
	view.co2Value = detailValue("—")

	view.calculateButton = widget.NewButton("Jetzt berechnen  →", view.calculate)
	view.calculateButton.Importance = widget.HighImportance
	view.printButton = widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), view.print)
	view.printButton.Importance = widget.LowImportance
	view.printButton.Disable()
	view.quantityEntry.OnSubmitted = func(string) { view.calculate() }
	view.quantityEntry.OnChanged = func(string) { view.status.Hide() }

	view.content = buildAppShell(view, win)
	return view
}

func buildAppShell(view *referenceView, win fyne.Window) fyne.CanvasObject {
	rail := buildNavigationRail(win, view.mode)
	header := buildHeader(view)
	workspace := buildWorkspace(view)

	mainBackground := canvas.NewRectangle(appPalette.background)
	mainContent := container.NewBorder(header, nil, nil, nil, workspace)
	main := container.NewStack(mainBackground, mainContent)

	return container.NewBorder(nil, nil, rail, nil, main)
}

func buildNavigationRail(win fyne.Window, activeMode string) fyne.CanvasObject {
	logoBackground := canvas.NewRectangle(appPalette.resultSurface)
	logoBackground.CornerRadius = 15
	logoText := canvas.NewText("e°", appPalette.textPrimary)
	logoText.TextSize = 19
	logoText.TextStyle = fyne.TextStyle{Bold: true}
	logo := container.NewGridWrap(
		fyne.NewSize(50, 50),
		container.NewStack(logoBackground, container.NewCenter(logoText)),
	)

	separator := canvas.NewRectangle(color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 70})
	separator.SetMinSize(fyne.NewSize(1, 38))

	briquettes := newFuelNavigationButton("B", "Briketts", activeMode == "briketts", func() {
		win.SetContent(buildView(win, "briketts"))
	})
	oil := newFuelNavigationButton("O", "Heizöl", activeMode == "oil", func() {
		win.SetContent(buildView(win, "oil"))
	})

	footer := canvas.NewText("LYC.REU", appPalette.whiteMuted)
	footer.TextSize = 9
	footer.TextStyle = fyne.TextStyle{Bold: true}

	navigation := container.NewVBox(
		container.NewCenter(logo),
		verticalGap(ui.spacingMedium),
		container.NewCenter(container.NewGridWrap(fyne.NewSize(1, 38), separator)),
		verticalGap(ui.spacingMedium),
		container.NewCenter(container.NewGridWrap(fyne.NewSize(66, 66), briquettes)),
		verticalGap(ui.spacingSmall),
		container.NewCenter(container.NewGridWrap(fyne.NewSize(66, 66), oil)),
	)
	inner := container.NewBorder(navigation, container.NewCenter(footer), nil, nil, nil)
	background := canvas.NewRectangle(appPalette.railBackground)
	background.SetMinSize(fyne.NewSize(ui.railWidth, 0))
	return container.NewStack(background, container.NewPadded(inner))
}

func buildHeader(view *referenceView) fyne.CanvasObject {
	brandBars := container.NewHBox(
		coloredBar(appPalette.accent, 8, 19),
		coloredBar(appPalette.coral, 8, 12),
	)
	brand := canvas.NewText("Emissionsrechner", appPalette.textPrimary)
	brand.TextSize = 14
	brand.TextStyle = fyne.TextStyle{Bold: true}
	brandGroup := container.NewHBox(brandBars, brand)

	readyDot := canvas.NewCircle(appPalette.success)
	readyIndicator := container.NewGridWrap(fyne.NewSize(9, 9), readyDot)
	ready := canvas.NewText("Bereit", appPalette.textSecondary)
	ready.TextSize = 12
	actions := container.NewHBox(readyIndicator, ready, view.printButton)
	headerContent := container.NewBorder(nil, nil, brandGroup, actions, nil)

	separator := canvas.NewRectangle(appPalette.border)
	separator.SetMinSize(fyne.NewSize(1, 1))
	return container.NewVBox(
		container.NewPadded(headerContent),
		separator,
	)
}

func buildWorkspace(view *referenceView) fyne.CanvasObject {
	form := buildForm(view)
	result := buildResultCard(view)
	columns := container.NewGridWithColumns(2,
		container.NewPadded(form),
		container.NewPadded(container.NewCenter(result)),
	)
	frame := container.NewGridWrap(fyne.NewSize(900, 525), columns)
	return container.NewCenter(frame)
}

func buildForm(view *referenceView) fyne.CanvasObject {
	stepCircle := canvas.NewCircle(color.Transparent)
	stepCircle.StrokeColor = appPalette.accent
	stepCircle.StrokeWidth = 1
	stepNumber := canvas.NewText("01", appPalette.accent)
	stepNumber.TextSize = 9
	stepNumber.TextStyle = fyne.TextStyle{Bold: true}
	step := container.NewHBox(
		container.NewGridWrap(fyne.NewSize(30, 30), container.NewStack(stepCircle, container.NewCenter(stepNumber))),
		canvasText("LIEFERMENGE", 10, appPalette.accent, true),
	)

	title := container.NewVBox(
		canvasText(titleForMode(view.mode)+".", 38, appPalette.textPrimary, true),
		canvasText("Einfach berechnet.", 38, appPalette.textPrimary, true),
	)
	description := container.NewVBox(
		canvasText("Menge eingeben und CO2-Ausstoß, Energiegehalt", 12, appPalette.textSecondary, false),
		canvasText("sowie Kostenanteil auf einen Blick erhalten.", 12, appPalette.textSecondary, false),
	)

	quantityLabel := canvasText("MENGE EINGEBEN", 10, appPalette.textSecondary, true)
	unit := canvasText(unitForMode(view.mode), 14, appPalette.accent, true)
	inputLine := canvas.NewRectangle(appPalette.textPrimary)
	inputLine.SetMinSize(fyne.NewSize(1, 2))
	input := container.NewBorder(nil, inputLine, nil, unit, view.quantityEntry)

	return container.NewVBox(
		step,
		verticalGap(ui.spacingLarge),
		title,
		verticalGap(ui.spacingMedium),
		description,
		verticalGap(ui.spacingLarge),
		quantityLabel,
		verticalGap(ui.spacingSmall),
		input,
		verticalGap(ui.spacingSmall),
		view.status,
		verticalGap(ui.spacingMedium),
		container.NewHBox(view.calculateButton),
	)
}

func buildResultCard(view *referenceView) fyne.CanvasObject {
	label := canvasText("GESAMTEMISSIONEN", 10, appPalette.textSecondary, true)
	live := canvasText("Live-Ergebnis", 11, appPalette.textSecondary, false)
	metrics := container.NewGridWithColumns(3,
		resultMetric("CO2-Kostenanteil", view.costValue),
		resultMetric("Energiegehalt", view.energyValue),
		resultMetric("CO2 pro kWh", view.co2Value),
	)

	separator := canvas.NewRectangle(color.NRGBA{R: 0x18, G: 0x21, B: 0x3d, A: 42})
	separator.SetMinSize(fyne.NewSize(1, 1))
	cardContent := container.NewVBox(
		container.NewBorder(nil, nil, label, live, nil),
		verticalGap(52),
		view.resultValue,
		verticalGap(ui.spacingSmall),
		view.resultHint,
		verticalGap(42),
		separator,
		verticalGap(ui.spacingMedium),
		metrics,
	)

	cardBackground := canvas.NewRectangle(appPalette.resultSurface)
	cardBackground.CornerRadius = ui.cardRadius
	card := container.NewStack(cardBackground, container.NewPadded(cardContent))

	shadowBackground := canvas.NewRectangle(appPalette.resultShadow)
	shadowBackground.CornerRadius = ui.cardRadius
	shadow := container.NewBorder(verticalGap(12), nil, horizontalGap(12), nil, shadowBackground)
	foreground := container.NewBorder(nil, verticalGap(12), nil, horizontalGap(12), card)
	stack := container.NewStack(shadow, foreground)
	return container.NewGridWrap(fyne.NewSize(420, 410), stack)
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
	view.resultValue.Refresh()
	view.resultHint.Text = fmt.Sprintf("Berechnet für %s %s %s", view.quantityEntry.Text, unitForMode(view.mode), titleForMode(view.mode))
	view.resultHint.Refresh()
	view.costValue.Text = view.result.EmissionCost
	view.costValue.Refresh()
	view.energyValue.Text = view.result.EnergyContent
	view.energyValue.Refresh()
	view.co2Value.Text = view.result.CO2PerKWh
	view.co2Value.Refresh()
	view.printButton.Enable()
	setStatus(view.status, "Ergebnis aktualisiert", appPalette.success)
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

func titleForMode(mode string) string {
	if mode == "oil" {
		return "Heizöl"
	}
	return "Briketts"
}

func unitForMode(mode string) string {
	if mode == "oil" {
		return "Liter"
	}
	return "Tonnen"
}

func detailValue(value string) *canvas.Text {
	text := canvas.NewText(value, appPalette.textPrimary)
	text.TextSize = 12
	text.TextStyle = fyne.TextStyle{Bold: true}
	return text
}

func resultMetric(label string, value *canvas.Text) fyne.CanvasObject {
	return container.NewVBox(
		canvasText(label, 9, appPalette.textSecondary, false),
		verticalGap(ui.spacingSmall),
		value,
	)
}

func canvasText(value string, size float32, textColor color.Color, bold bool) *canvas.Text {
	text := canvas.NewText(value, textColor)
	text.TextSize = size
	text.TextStyle = fyne.TextStyle{Bold: bold}
	return text
}

func coloredBar(fill color.Color, width, height float32) fyne.CanvasObject {
	bar := canvas.NewRectangle(fill)
	bar.CornerRadius = width / 2
	return container.NewGridWrap(fyne.NewSize(width, height), bar)
}

func verticalGap(height float32) fyne.CanvasObject {
	gap := canvas.NewRectangle(color.Transparent)
	gap.SetMinSize(fyne.NewSize(1, height))
	return gap
}

func horizontalGap(width float32) fyne.CanvasObject {
	gap := canvas.NewRectangle(color.Transparent)
	gap.SetMinSize(fyne.NewSize(width, 1))
	return gap
}

func newStatusText() *canvas.Text {
	status := canvas.NewText("", appPalette.textMuted)
	status.TextSize = 11
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

type fuelNavigationButton struct {
	widget.BaseWidget
	iconLabel string
	label     string
	active    bool
	onTapped  func()
}

func newFuelNavigationButton(iconLabel, label string, active bool, onTapped func()) *fuelNavigationButton {
	button := &fuelNavigationButton{
		iconLabel: iconLabel,
		label:     label,
		active:    active,
		onTapped:  onTapped,
	}
	button.ExtendBaseWidget(button)
	return button
}

func (button *fuelNavigationButton) Tapped(*fyne.PointEvent) {
	if button.onTapped != nil {
		button.onTapped()
	}
}

func (button *fuelNavigationButton) CreateRenderer() fyne.WidgetRenderer {
	var backgroundColor color.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0}
	textColor := appPalette.whiteMuted
	if button.active {
		backgroundColor = appPalette.white
		textColor = appPalette.accent
	}
	background := canvas.NewRectangle(backgroundColor)
	background.CornerRadius = 17
	icon := canvasText(button.iconLabel, 17, textColor, true)
	label := canvasText(button.label, 9, textColor, true)
	content := container.NewCenter(container.NewVBox(
		container.NewCenter(icon),
		verticalGap(ui.spacingTiny),
		container.NewCenter(label),
	))
	return widget.NewSimpleRenderer(container.NewStack(background, content))
}
