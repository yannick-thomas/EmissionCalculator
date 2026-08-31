package ui

import (
	"emissioncalculator/internal/calculation"
	"emissioncalculator/internal/validation"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const co2PricePreferenceKey = "calculation.co2_price_per_tonne"
const factorYearPreferenceKey = "calculation.factor_year"

type settingsStore struct {
	preferences fyne.Preferences
	config      calculation.Config
}

func newSettingsStore(preferences fyne.Preferences) *settingsStore {
	defaults := calculation.DefaultConfig()
	price := defaults.CO2PricePerTonne
	year := defaults.Year
	if preferences != nil {
		price = preferences.FloatWithFallback(co2PricePreferenceKey, price)
		year = preferences.IntWithFallback(factorYearPreferenceKey, year)
	}
	if price <= 0 {
		price = defaults.CO2PricePerTonne
	}
	if !isAvailableYear(year) {
		year = defaults.Year
	}
	return &settingsStore{
		preferences: preferences,
		config:      calculation.Config{CO2PricePerTonne: price, Year: year},
	}
}

func (store *settingsStore) Config() calculation.Config {
	return store.config
}

func (store *settingsStore) SetCO2Price(price float64) {
	store.config.CO2PricePerTonne = price
	if store.preferences != nil {
		store.preferences.SetFloat(co2PricePreferenceKey, price)
	}
}

func (store *settingsStore) SetFactorYear(year int) {
	store.config.Year = year
	if store.preferences != nil {
		store.preferences.SetInt(factorYearPreferenceKey, year)
	}
}

func (store *settingsStore) Reset() {
	store.config = calculation.DefaultConfig()
	if store.preferences != nil {
		store.preferences.RemoveValue(co2PricePreferenceKey)
		store.preferences.RemoveValue(factorYearPreferenceKey)
	}
}

func isAvailableYear(year int) bool {
	for _, available := range calculation.AvailableYears() {
		if available == year {
			return true
		}
	}
	return false
}

type settingsPanel struct {
	store        *settingsStore
	priceEntry   *focusEntry
	priceControl *quantityControl
	yearSelect   *widget.Select
	factorHint   *canvas.Text
	status       *canvas.Text
	content      fyne.CanvasObject
	popup        *widget.PopUp
	onSaved      func()
	resetDefault bool
}

func showSettingsDialog(window fyne.Window, store *settingsStore, onSaved func()) *settingsPanel {
	panel := newSettingsPanel(store, onSaved)
	panel.popup = widget.NewModalPopUp(panel.content, window.Canvas())
	panel.popup.Show()
	window.Canvas().Focus(panel.priceEntry)
	return panel
}

func newSettingsPanel(store *settingsStore, onSaved func()) *settingsPanel {
	panel := &settingsPanel{store: store, onSaved: onSaved}
	panel.priceEntry = newFocusEntry()
	panel.priceEntry.SetText(formatSettingsValue(store.Config().CO2PricePerTonne))
	panel.priceControl = newQuantityControl(panel.priceEntry, "€/t")
	panel.status = canvasText(" ", 12, appPalette.textSecondary, true)
	panel.factorHint = canvasText(factorHintText(store.Config().Year), 11, appPalette.textMuted, false)

	yearOptions := yearSelectOptions()
	panel.yearSelect = widget.NewSelect(yearOptions, func(selected string) {
		panel.resetDefault = false
		if year, err := strconv.Atoi(selected); err == nil {
			panel.factorHint.Text = factorHintText(year)
			panel.factorHint.Refresh()
		}
	})
	panel.yearSelect.SetSelected(strconv.Itoa(store.Config().Year))
	yearField := container.NewGridWrap(fyne.NewSize(ui.formColWidth, 42), panel.yearSelect)

	iconBackground := canvas.NewCircle(appPalette.resultSurface)
	icon := canvas.NewImageFromResource(settingsIconResource())
	icon.FillMode = canvas.ImageFillContain
	iconFrame := container.NewGridWrap(
		fyne.NewSize(54, 54),
		container.NewStack(iconBackground, container.NewPadded(icon)),
	)
	heading := container.NewVBox(
		canvasText("B E R E C H N U N G S G R U N D L A G E N", 11, appPalette.accent, true),
		verticalGap(6),
		canvasText("Einstellungen", 28, appPalette.textPrimary, true),
	)
	header := container.NewHBox(iconFrame, horizontalGap(16), heading)
	description := container.NewVBox(
		canvasText("Lege den gemeinsamen CO₂-Preis für alle", 14, appPalette.textSecondary, false),
		canvasText("Brennstoffarten fest.", 14, appPalette.textSecondary, false),
	)

	resetButton := widget.NewButton("Standard wiederherstellen", panel.restoreDefault)
	resetButton.Importance = widget.LowImportance
	cancelButton := widget.NewButton("Abbrechen", panel.dismiss)
	applyButton := widget.NewButton("Übernehmen", panel.apply)
	applyButton.Importance = widget.HighImportance
	actions := container.NewVBox(
		container.NewHBox(resetButton),
		verticalGap(12),
		container.NewBorder(nil, nil, cancelButton, applyButton, nil),
	)

	separatorTop := canvas.NewRectangle(appPalette.border)
	separatorTop.SetMinSize(fyne.NewSize(1, 1))
	separatorBottom := canvas.NewRectangle(appPalette.border)
	separatorBottom.SetMinSize(fyne.NewSize(1, 1))
	form := container.NewVBox(
		header,
		verticalGap(14),
		description,
		verticalGap(16),
		separatorTop,
		verticalGap(14),
		canvasText("C O₂ - P R E I S", 12, appPalette.textSecondary, true),
		verticalGap(10),
		panel.priceControl.content,
		verticalGap(16),
		canvasText("F A K T O R J A H R", 12, appPalette.textSecondary, true),
		verticalGap(10),
		yearField,
		verticalGap(10),
		panel.factorHint,
		verticalGap(8),
		container.NewGridWrap(fyne.NewSize(ui.formColWidth, 20), panel.status),
		verticalGap(12),
		separatorBottom,
		verticalGap(12),
		actions,
	)

	background := canvas.NewRectangle(appPalette.surface)
	background.CornerRadius = 28
	background.StrokeColor = appPalette.border
	background.StrokeWidth = 1
	inset := container.NewBorder(verticalGap(28), verticalGap(24), horizontalGap(30), horizontalGap(30), form)
	panel.content = container.NewGridWrap(fyne.NewSize(540, 590), container.NewStack(background, inset))
	panel.priceEntry.OnSubmitted = func(string) { panel.apply() }
	panel.priceEntry.OnChanged = func(string) {
		panel.resetDefault = false
		panel.priceControl.SetError(false)
		setStatus(panel.status, " ", appPalette.textSecondary)
	}
	return panel
}

func (panel *settingsPanel) apply() {
	price, err := validation.ParseQuantity(panel.priceEntry.Text)
	if err != nil {
		panel.priceControl.SetError(true)
		setStatus(panel.status, "Bitte einen gültigen CO₂-Preis eingeben.", appPalette.error)
		return
	}
	year, err := strconv.Atoi(panel.yearSelect.Selected)
	if err != nil {
		setStatus(panel.status, "Bitte ein gültiges Faktorjahr wählen.", appPalette.error)
		return
	}
	defaults := calculation.DefaultConfig()
	if panel.resetDefault && price == defaults.CO2PricePerTonne && year == defaults.Year {
		panel.store.Reset()
	} else {
		panel.store.SetCO2Price(price)
		panel.store.SetFactorYear(year)
	}
	panel.dismiss()
	if panel.onSaved != nil {
		panel.onSaved()
	}
}

func (panel *settingsPanel) restoreDefault() {
	defaults := calculation.DefaultConfig()
	panel.priceEntry.SetText(formatSettingsValue(defaults.CO2PricePerTonne))
	panel.yearSelect.SetSelected(strconv.Itoa(defaults.Year))
	panel.resetDefault = true
	setStatus(panel.status, "Standardwert eingesetzt – mit Übernehmen speichern.", appPalette.success)
}

func (panel *settingsPanel) dismiss() {
	if panel.popup != nil {
		panel.popup.Hide()
	}
}

func formatSettingsValue(value float64) string {
	formatted := strconv.FormatFloat(value, 'f', 2, 64)
	formatted = strings.TrimRight(formatted, "0")
	formatted = strings.TrimRight(formatted, ".")
	return strings.ReplaceAll(formatted, ".", ",")
}

func yearSelectOptions() []string {
	years := calculation.AvailableYears()
	options := make([]string, len(years))
	for i, year := range years {
		options[i] = strconv.Itoa(year)
	}
	return options
}

// factorHintText shows the emission factors that actually apply for the given year.
func factorHintText(year int) string {
	asOf := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
	oil, oilErr := calculation.FactorFor(calculation.FuelOil, asOf)
	briquettes, briquettesErr := calculation.FactorFor(calculation.FuelBriquettes, asOf)
	if oilErr != nil || briquettesErr != nil {
		return "Für dieses Jahr sind keine Faktoren hinterlegt."
	}
	return "Heizöl: EF " + formatFloat(oil.EmissionFactor, 4) + " · Briketts: EF " + formatFloat(briquettes.EmissionFactor, 4) + " kg CO₂/MJ"
}
