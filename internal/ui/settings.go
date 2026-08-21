package ui

import (
	"emissioncalculator/internal/calculation"
	"emissioncalculator/internal/validation"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const co2PricePreferenceKey = "calculation.co2_price_per_tonne"

type settingsStore struct {
	preferences fyne.Preferences
	config      calculation.Config
}

func newSettingsStore(preferences fyne.Preferences) *settingsStore {
	defaults := calculation.DefaultConfig()
	price := defaults.CO2PricePerTonne
	if preferences != nil {
		price = preferences.FloatWithFallback(co2PricePreferenceKey, price)
	}
	if price <= 0 {
		price = defaults.CO2PricePerTonne
	}
	return &settingsStore{
		preferences: preferences,
		config:      calculation.Config{CO2PricePerTonne: price},
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

func (store *settingsStore) Reset() {
	store.config = calculation.DefaultConfig()
	if store.preferences != nil {
		store.preferences.RemoveValue(co2PricePreferenceKey)
	}
}

type settingsPanel struct {
	store        *settingsStore
	priceEntry   *focusEntry
	priceControl *quantityControl
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

	iconBackground := canvas.NewCircle(appPalette.accentSoft)
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
		verticalGap(10),
		container.NewBorder(nil, nil, nil, container.NewHBox(cancelButton, applyButton), nil),
	)

	separatorTop := canvas.NewRectangle(appPalette.border)
	separatorTop.SetMinSize(fyne.NewSize(1, 1))
	separatorBottom := canvas.NewRectangle(appPalette.border)
	separatorBottom.SetMinSize(fyne.NewSize(1, 1))
	form := container.NewVBox(
		header,
		verticalGap(20),
		description,
		verticalGap(24),
		separatorTop,
		verticalGap(22),
		canvasText("C O₂ - P R E I S", 12, appPalette.textSecondary, true),
		verticalGap(10),
		panel.priceControl.content,
		verticalGap(10),
		canvasText("Dieser Wert gilt für Heizöl und Briketts.", 12, appPalette.textSecondary, false),
		verticalGap(8),
		container.NewGridWrap(fyne.NewSize(ui.formColWidth, 20), panel.status),
		verticalGap(18),
		separatorBottom,
		verticalGap(18),
		actions,
	)

	background := canvas.NewRectangle(appPalette.surface)
	background.CornerRadius = 28
	background.StrokeColor = appPalette.border
	background.StrokeWidth = 1
	inset := container.NewBorder(verticalGap(32), verticalGap(28), horizontalGap(34), horizontalGap(34), form)
	panel.content = container.NewGridWrap(fyne.NewSize(548, 548), container.NewStack(background, inset))
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
	if panel.resetDefault && price == calculation.DefaultConfig().CO2PricePerTonne {
		panel.store.Reset()
	} else {
		panel.store.SetCO2Price(price)
	}
	if panel.onSaved != nil {
		panel.onSaved()
	}
	panel.dismiss()
}

func (panel *settingsPanel) restoreDefault() {
	defaultPrice := calculation.DefaultConfig().CO2PricePerTonne
	panel.priceEntry.SetText(formatSettingsValue(defaultPrice))
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
