package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func buildSharedNavigation(window fyne.Window, settings *settingsStore, oilView, briquettesView *referenceView) fyne.CanvasObject {
	logoBackground := canvas.NewImageFromResource(logoPanelResource())
	logoBackground.FillMode = canvas.ImageFillStretch
	logoText := canvasText("e°", 22, appPalette.textPrimary, true)
	logo := container.NewGridWrap(fyne.NewSize(64, 64), container.NewStack(logoBackground, container.NewCenter(logoText)))
	var oilBtn, briquettesBtn *fuelNavigationButton
	oilBtn = newFuelNavigationButton(modeOil, "Heizöl", true, func() {
		oilBtn.active = true
		briquettesBtn.active = false
		oilBtn.Refresh()
		briquettesBtn.Refresh()
		oilView.content.Show()
		briquettesView.content.Hide()
	})
	briquettesBtn = newFuelNavigationButton(modeBriquettes, "Briketts", false, func() {
		briquettesBtn.active = true
		oilBtn.active = false
		oilBtn.Refresh()
		briquettesBtn.Refresh()
		briquettesView.content.Show()
		oilView.content.Hide()
	})
	settingsButton := newCircleIconButton(settingsIconResource(), func() {
		showSettingsDialog(window, settings, func() {
			oilView.refreshForSettingsChange()
			briquettesView.refreshForSettingsChange()
		})
	})
	navigationDivider := canvas.NewRectangle(appPalette.border)
	navigationDivider.SetMinSize(fyne.NewSize(1, 34))
	navigationActions := container.NewHBox(
		container.NewGridWrap(fyne.NewSize(136, 62), briquettesBtn),
		horizontalGap(10),
		container.NewGridWrap(fyne.NewSize(136, 62), oilBtn),
		horizontalGap(14),
		container.NewCenter(navigationDivider),
		horizontalGap(14),
		container.NewGridWrap(fyne.NewSize(54, 54), settingsButton),
	)
	barContent := container.NewBorder(nil, nil, logo, navigationActions, nil)
	barFrame := container.NewGridWrap(fyne.NewSize(ui.contentWidth, ui.topBarHeight), barContent)
	background := canvas.NewRectangle(appPalette.railBackground)
	background.SetMinSize(fyne.NewSize(0, ui.topBarHeight))
	separator := canvas.NewRectangle(appPalette.border)
	separator.SetMinSize(fyne.NewSize(0, 1))
	return container.NewBorder(nil, separator, nil, nil, container.NewStack(background, container.NewCenter(barFrame)))
}
