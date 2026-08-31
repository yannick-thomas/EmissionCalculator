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
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const co2PricePreferenceKey = "calculation.co2_price_per_tonne"

// The persisted key keeps its historic name so existing user settings survive the terminology fix.
const calculationYearPreferenceKey = "calculation.factor_year"

type settingsStore struct {
	preferences fyne.Preferences
	config      calculation.Config
}

func newSettingsStore(preferences fyne.Preferences) *settingsStore {
	defaults := calculation.DefaultConfig()
	price := defaults.CO2PricePerTonne
	year := defaults.CalculationYear
	if preferences != nil {
		price = preferences.FloatWithFallback(co2PricePreferenceKey, price)
		year = preferences.IntWithFallback(calculationYearPreferenceKey, year)
	}
	if price <= 0 {
		price = defaults.CO2PricePerTonne
	}
	if !isAvailableYear(year) {
		year = defaults.CalculationYear
	}
	config := defaults
	config.CalculationYear = year
	config.Year = year
	if price != defaults.CO2PricePerTonne {
		config = config.WithCO2Price(price, "Manuelle Preisannahme")
	}
	return &settingsStore{
		preferences: preferences,
		config:      config,
	}
}

func (store *settingsStore) Config() calculation.Config {
	return store.config
}

func (store *settingsStore) SetCO2Price(price float64) {
	store.config = store.config.WithCO2Price(price, "Manuelle Preisannahme")
	if store.preferences != nil {
		store.preferences.SetFloat(co2PricePreferenceKey, price)
	}
}

func (store *settingsStore) SetCalculationYear(year int) {
	store.config.CalculationYear = year
	store.config.Year = year
	if store.config.PriceReference.Source == "Manuelle Preisannahme" {
		store.config = store.config.WithCO2Price(store.config.CO2PricePerTonne, "Manuelle Preisannahme")
	}
	if store.preferences != nil {
		store.preferences.SetInt(calculationYearPreferenceKey, year)
	}
}

func (store *settingsStore) Reset() {
	store.config = calculation.DefaultConfig()
	if store.preferences != nil {
		store.preferences.RemoveValue(co2PricePreferenceKey)
		store.preferences.RemoveValue(calculationYearPreferenceKey)
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
	factorHint   *widget.Label
	status       *canvas.Text
	content      fyne.CanvasObject
	popup        dialog.Dialog
	onSaved      func()
	resetDefault bool
}

func showSettingsDialog(window fyne.Window, store *settingsStore, onSaved func()) *settingsPanel {
	panel := newSettingsPanel(store, onSaved)
	panel.popup = showResponsiveDialogWithoutButtons("Einstellungen", panel.content, window, fyne.NewSize(640, 740))
	window.Canvas().Focus(panel.priceEntry)
	return panel
}

func newSettingsPanel(store *settingsStore, onSaved func()) *settingsPanel {
	panel := &settingsPanel{store: store, onSaved: onSaved}
	panel.priceEntry = newFocusEntry()
	panel.priceEntry.SetText(formatSettingsValue(store.Config().CO2PricePerTonne))
	panel.priceControl = newQuantityControl(panel.priceEntry, "€/t")
	panel.status = canvasText(" ", 12, appPalette.textSecondary, true)
	panel.factorHint = widget.NewLabel(factorHintText(store.Config().CalculationYear))
	panel.factorHint.Wrapping = fyne.TextWrapWord

	yearOptions := yearSelectOptions()
	panel.yearSelect = widget.NewSelect(yearOptions, func(selected string) {
		panel.resetDefault = false
		if year, err := strconv.Atoi(selected); err == nil {
			panel.factorHint.SetText(factorHintText(year))
		}
	})
	panel.yearSelect.SetSelected(strconv.Itoa(store.Config().CalculationYear))
	yearField := container.NewBorder(nil, nil, nil, nil, panel.yearSelect)

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
		canvasText("Lege CO₂-Preis und Berechnungsjahr getrennt fest.", 14, appPalette.textSecondary, false),
		canvasText("Manuelle Preise werden als Annahme dokumentiert.", 14, appPalette.textSecondary, false),
	)

	resetButton := widget.NewButton("Zurücksetzen", panel.restoreDefault)
	cancelButton := widget.NewButton("Abbrechen", panel.dismiss)
	applyButton := widget.NewButton("Übernehmen", panel.apply)
	applyButton.Importance = widget.HighImportance
	actions := container.NewGridWithColumns(3, resetButton, cancelButton, applyButton)

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
		canvasText("B E R E C H N U N G S J A H R", 12, appPalette.textSecondary, true),
		verticalGap(10),
		yearField,
		verticalGap(10),
		panel.factorHint,
		verticalGap(8),
		panel.status,
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
	card := container.NewStack(background, inset)
	cardScroll := container.NewScroll(card)
	panel.content = container.New(
		&centeredCardLayout{preferred: fyne.NewSize(560, 650), minimum: fyne.NewSize(360, 320), margin: 20},
		cardScroll,
	)
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
		setStatus(panel.status, "Bitte ein gültiges Berechnungsjahr wählen.", appPalette.error)
		return
	}
	defaults := calculation.DefaultConfig()
	if price == defaults.CO2PricePerTonne {
		panel.store.Reset()
		if year != defaults.CalculationYear {
			panel.store.SetCalculationYear(year)
		}
	} else {
		panel.store.SetCalculationYear(year)
		panel.store.SetCO2Price(price)
	}
	panel.dismiss()
	if panel.onSaved != nil {
		panel.onSaved()
	}
}

func (panel *settingsPanel) restoreDefault() {
	defaults := calculation.DefaultConfig()
	panel.priceEntry.SetText(formatSettingsValue(defaults.CO2PricePerTonne))
	panel.yearSelect.SetSelected(strconv.Itoa(defaults.CalculationYear))
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

// factorHintText distinguishes the calculation year from factor validity and source year.
func factorHintText(year int) string {
	asOf := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
	parts := make([]string, 0, len(calculation.Catalog))
	for _, descriptor := range calculation.Catalog {
		factor, err := calculation.FactorFor(descriptor.Fuel, asOf)
		if err != nil {
			return "Für dieses Jahr sind nicht alle Faktoren hinterlegt."
		}
		parts = append(parts, descriptor.Label+" "+formatFloat(factor.EmissionFactor, 4))
	}
	return "Paket " + calculation.CurrentFactorPack().ID + " · EF kg CO₂/MJ: " + strings.Join(parts, " / ")
}
