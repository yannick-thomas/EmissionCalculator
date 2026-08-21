package ui

import (
	"emissioncalculator/internal/calculation"
	"emissioncalculator/internal/models"
	"emissioncalculator/internal/pdf"
	"emissioncalculator/internal/validation"
	"fmt"
	"image/color"
	"math"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

var exportLabel = pdf.ExportLabel
var openExportedFile = openFile

type referenceView struct {
	content               fyne.CanvasObject
	scroll                *container.Scroll
	quantityEntry         *focusEntry
	quantityControl       *quantityControl
	calculateButton       *actionButton
	printButton           *circleIconButton
	status                *canvas.Text
	headerStatus          *canvas.Text
	headerStatusDot       *canvas.Circle
	result                models.CalculationResult
	mode                  string
	resultValue           *canvas.Text
	resultUnit            *canvas.Text
	resultValueRow        *fyne.Container
	resultHint            *canvas.Text
	resultBackground      *canvas.Image
	resultBadge           *canvas.Text
	resultBadgeBackground *canvas.Rectangle
	costValue             *canvas.Text
	energyValue           *canvas.Text
	co2Value              *canvas.Text
}

func NewRootWindow(app fyne.App) fyne.Window {
	app.Settings().SetTheme(emissionTheme{Theme: theme.LightTheme()})
	window := app.NewWindow("Emissionsrechner")
	window.Resize(fyne.NewSize(1080, 720))
	window.SetContent(buildView(window, "oil"))
	return window
}

func buildView(window fyne.Window, mode string) fyne.CanvasObject {
	return buildReferenceView(window, mode).content
}

func buildReferenceView(window fyne.Window, mode string) *referenceView {
	view := &referenceView{mode: mode}
	view.quantityEntry = newFocusEntry()
	view.quantityEntry.SetPlaceHolder("z. B. 12.500 oder 12,5")
	view.status = newStatusText()

	view.resultValue = canvasText("—", 72, appPalette.textPrimary, true)
	view.resultUnit = canvasText("", 23, appPalette.textPrimary, true)
	view.resultHint = canvasText("Noch keine Berechnung", 16, appPalette.textSecondary, false)
	view.costValue = detailValue("—")
	view.energyValue = detailValue("—")
	view.co2Value = detailValue("—")

	view.calculateButton = newActionButton("Jetzt berechnen", view.calculate)
	view.printButton = newCircleIconButton(printIconResource(), view.print)
	view.printButton.Disable()

	view.content = buildAppShell(view, window)
	view.quantityEntry.OnSubmitted = func(string) { view.calculate() }
	view.quantityEntry.OnChanged = func(string) {
		view.quantityControl.SetError(false)
		view.clearResult()
		setStatus(view.status, " ", appPalette.textSecondary)
		view.setHeaderStatus("Bereit", appPalette.success)
	}
	view.setHeaderStatus("Bereit", appPalette.success)
	return view
}

func buildAppShell(view *referenceView, window fyne.Window) fyne.CanvasObject {
	navigation := buildTopNavigation(window, view.mode)
	header := buildHeader(view)
	form := container.NewGridWrap(fyne.NewSize(ui.contentWidth, 560), buildForm(view))
	result := container.NewGridWrap(fyne.NewSize(ui.contentWidth, ui.resultHeight), buildResultCard(view))
	page := container.NewVBox(
		header,
		verticalGap(36),
		form,
		verticalGap(72),
		result,
		verticalGap(42),
	)
	pageFrame := container.NewGridWrap(fyne.NewSize(ui.contentWidth, ui.pageHeight), page)
	scroll := container.NewVScroll(container.NewCenter(pageFrame))
	view.scroll = scroll
	main := container.NewStack(canvas.NewRectangle(appPalette.background), scroll)
	return container.NewBorder(navigation, nil, nil, nil, main)
}

func buildTopNavigation(window fyne.Window, activeMode string) fyne.CanvasObject {
	logoBackground := canvas.NewImageFromResource(logoPanelResource())
	logoBackground.FillMode = canvas.ImageFillStretch
	logoText := canvasText("e°", 22, appPalette.textPrimary, true)
	logo := container.NewGridWrap(fyne.NewSize(64, 64), container.NewStack(logoBackground, container.NewCenter(logoText)))
	briquettes := newFuelNavigationButton("briquettes", "Briketts", activeMode == "briketts", func() {
		window.SetContent(buildView(window, "briketts"))
	})
	oil := newFuelNavigationButton("oil", "Heizöl", activeMode == "oil", func() {
		window.SetContent(buildView(window, "oil"))
	})
	navigationActions := container.NewHBox(
		container.NewGridWrap(fyne.NewSize(136, 62), briquettes),
		horizontalGap(10),
		container.NewGridWrap(fyne.NewSize(136, 62), oil),
	)
	barContent := container.NewBorder(nil, nil, logo, navigationActions, nil)
	barFrame := container.NewGridWrap(fyne.NewSize(ui.contentWidth, ui.topBarHeight), barContent)
	background := canvas.NewRectangle(appPalette.railBackground)
	background.SetMinSize(fyne.NewSize(0, ui.topBarHeight))
	return container.NewStack(background, container.NewCenter(barFrame))
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
	statusGroup := container.NewGridWrap(fyne.NewSize(170, 54), container.NewCenter(container.NewHBox(readyIndicator, view.headerStatus)))
	actions := container.NewHBox(statusGroup, horizontalGap(8), container.NewGridWrap(fyne.NewSize(54, 54), view.printButton))
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
		canvasText(titleForMode(view.mode)+".", 48, appPalette.textPrimary, true),
		canvasText("Einfach berechnet.", 48, appPalette.textPrimary, true),
	)
	description := container.NewVBox(
		canvasText("Menge eingeben und CO₂-Ausstoß, Energiegehalt sowie", 17, appPalette.textSecondary, false),
		canvasText("Kostenanteil auf einen Blick erhalten.", 17, appPalette.textSecondary, false),
	)

	view.quantityControl = newQuantityControl(view.quantityEntry, unitForMode(view.mode))
	statusFrame := container.NewGridWrap(fyne.NewSize(ui.contentWidth, 20), view.status)
	return container.NewVBox(
		step,
		verticalGap(30),
		title,
		verticalGap(20),
		description,
		verticalGap(38),
		canvasText("M E N G E  E I N G E B E N", 12, appPalette.textSecondary, true),
		verticalGap(12),
		view.quantityControl.content,
		verticalGap(8),
		statusFrame,
		verticalGap(16),
		container.NewHBox(view.calculateButton),
	)
}

func buildResultCard(view *referenceView) fyne.CanvasObject {
	label := canvasText("G E S A M T E M I S S I O N E N", 13, appPalette.textSecondary, true)
	view.resultBadge = canvasText("Live-Ergebnis", 14, appPalette.textSecondary, true)
	view.resultBadgeBackground = canvas.NewRectangle(color.Transparent)
	view.resultBadgeBackground.CornerRadius = 22
	view.resultBadgeBackground.StrokeColor = color.NRGBA{R: 0x62, G: 0x77, B: 0x2e, A: 90}
	view.resultBadgeBackground.StrokeWidth = 1
	badge := container.NewGridWrap(
		fyne.NewSize(136, 44),
		container.NewStack(view.resultBadgeBackground, container.NewCenter(view.resultBadge)),
	)
	metrics := container.NewGridWithColumns(3,
		resultMetric("CO₂-Kostenanteil", view.costValue),
		resultMetric("Energiegehalt", view.energyValue),
		resultMetric("CO₂ pro kWh", view.co2Value),
	)
	separator := canvas.NewRectangle(color.NRGBA{R: 0x18, G: 0x21, B: 0x3d, A: 36})
	separator.SetMinSize(fyne.NewSize(1, 1))
	view.resultValueRow = container.NewVBox(view.resultValue, view.resultUnit)
	cardContent := container.NewVBox(
		container.NewBorder(nil, nil, label, badge, nil),
		verticalGap(58),
		view.resultValueRow,
		verticalGap(18),
		view.resultHint,
		verticalGap(42),
		separator,
		verticalGap(28),
		metrics,
	)

	view.resultBackground = canvas.NewImageFromResource(resultPanelResource(false))
	view.resultBackground.FillMode = canvas.ImageFillStretch
	contentInset := container.NewBorder(verticalGap(42), verticalGap(42), horizontalGap(42), horizontalGap(42), cardContent)
	card := container.NewStack(view.resultBackground, contentInset)
	shadowImage := canvas.NewImageFromResource(resultShadowResource())
	shadowImage.FillMode = canvas.ImageFillStretch
	return container.NewStack(shadowImage, card)
}

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
	if view.mode == "oil" {
		view.result = calculation.CalculateOil(input)
	} else {
		view.result = calculation.CalculateBriquettes(input)
	}
	value, unit := splitEmissions(view.result.Emissions)
	value = formatGermanNumberString(value, 2)
	view.resultValue.Text = value
	view.resultValue.TextSize = resultTextSize(value)
	view.resultUnit.Text = unit
	view.resultHint.Text = fmt.Sprintf("Berechnet für %s %s %s", formatQuantityDisplay(input), unitForMode(view.mode), titleForMode(view.mode))
	view.costValue.Text = formatMeasurement(view.result.EmissionCost, 2)
	view.energyValue.Text = formatMeasurement(view.result.EnergyContent, 2)
	view.co2Value.Text = formatGermanNumberString(view.result.CO2PerKWh, 4) + " kg"
	fitDetailText(view.costValue)
	fitDetailText(view.energyValue)
	fitDetailText(view.co2Value)
	view.resultBackground.Resource = resultPanelResource(true)
	view.printButton.Enable()
	refreshTexts(view.resultValue, view.resultUnit, view.resultHint, view.costValue, view.energyValue, view.co2Value, view.resultBadge)
	view.resultValueRow.Refresh()
	view.resultBackground.Refresh()
	setStatus(view.status, "Ergebnis aktualisiert", appPalette.success)
	view.setHeaderStatus("Berechnet", appPalette.success)
}

func (view *referenceView) clearResult() {
	view.result = models.CalculationResult{}
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
	refreshTexts(view.resultValue, view.resultUnit, view.resultHint, view.costValue, view.energyValue, view.co2Value, view.resultBadge)
	view.resultValueRow.Refresh()
	view.resultBackground.Refresh()
}

func (view *referenceView) print() {
	if !view.result.Valid {
		setStatus(view.status, "Bitte zuerst eine gültige Berechnung durchführen.", appPalette.error)
		view.setHeaderStatus("Eingabe erforderlich", appPalette.error)
		return
	}
	view.setHeaderStatus("PDF wird erstellt", appPalette.accent)
	path, err := exportLabel(view.result.Emissions, view.result.EmissionCost, view.result.EnergyContent, view.result.CO2PerKWh)
	if err != nil {
		setStatus(view.status, "Fehler beim PDF-Export: "+err.Error(), appPalette.error)
		view.setHeaderStatus("PDF-Fehler", appPalette.error)
		return
	}
	if err := openExportedFile(path); err != nil {
		setStatus(view.status, "PDF erstellt: "+path, appPalette.success)
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

func splitEmissions(value string) (string, string) {
	const suffix = " kg CO2"
	if strings.HasSuffix(value, suffix) {
		return strings.TrimSuffix(value, suffix), "kg CO₂"
	}
	return value, ""
}

func resultTextSize(value string) float32 {
	switch {
	case len(value) <= 10:
		return 72
	case len(value) <= 13:
		return 62
	default:
		return 52
	}
}

func fitDetailText(text *canvas.Text) {
	switch {
	case len(text.Text) <= 18:
		text.TextSize = 18
	case len(text.Text) <= 20:
		text.TextSize = 16
	default:
		text.TextSize = 14
	}
}

func formatQuantityDisplay(value float64) string {
	whole := int64(math.Floor(value))
	digits := strconv.FormatInt(whole, 10)
	for index := len(digits) - 3; index > 0; index -= 3 {
		digits = digits[:index] + "." + digits[index:]
	}
	fraction := value - float64(whole)
	if math.Abs(fraction) < 0.005 {
		return digits
	}
	decimal := fmt.Sprintf("%.2f", fraction)
	decimal = strings.TrimPrefix(decimal, "0")
	decimal = strings.TrimRight(decimal, "0")
	return digits + strings.ReplaceAll(decimal, ".", ",")
}

func formatMeasurement(value string, decimals int) string {
	parts := strings.SplitN(strings.TrimSpace(value), " ", 2)
	formatted := formatGermanNumberString(parts[0], decimals)
	if len(parts) == 1 {
		return formatted
	}
	return formatted + " " + parts[1]
}

func formatGermanNumberString(value string, decimals int) string {
	normalized := strings.ReplaceAll(strings.TrimSpace(value), ".", "")
	normalized = strings.ReplaceAll(normalized, ",", ".")
	number, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return value
	}
	factor := math.Pow10(decimals)
	number = math.Round((number+math.Copysign(1e-9, number))*factor) / factor
	formatted := strconv.FormatFloat(number, 'f', decimals, 64)
	parts := strings.SplitN(formatted, ".", 2)
	integer := parts[0]
	sign := ""
	if strings.HasPrefix(integer, "-") {
		sign = "-"
		integer = strings.TrimPrefix(integer, "-")
	}
	for index := len(integer) - 3; index > 0; index -= 3 {
		integer = integer[:index] + "." + integer[index:]
	}
	if len(parts) == 1 || decimals == 0 {
		return sign + integer
	}
	return sign + integer + "," + parts[1]
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
	return canvasText(value, 18, appPalette.textPrimary, true)
}

func resultMetric(label string, value *canvas.Text) fyne.CanvasObject {
	return container.NewVBox(
		canvasText(label, 13, appPalette.textSecondary, false),
		verticalGap(16),
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
	return canvasText(" ", 11, appPalette.textSecondary, true)
}

func setStatus(status *canvas.Text, message string, colorValue color.Color) {
	status.Text = message
	status.Color = colorValue
	status.Refresh()
}

func refreshTexts(texts ...*canvas.Text) {
	for _, text := range texts {
		text.Refresh()
	}
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
