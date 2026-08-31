package ui

import (
	"emissioncalculator/internal/calculation"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func buildSharedNavigation(window fyne.Window, settings *settingsStore, history *historyController, views []*referenceView) fyne.CanvasObject {
	logoBackground := canvas.NewImageFromResource(logoPanelResource())
	logoBackground.FillMode = canvas.ImageFillStretch
	logoText := canvasText("e°", 22, appPalette.textPrimary, true)
	logo := container.NewGridWrap(fyne.NewSize(64, 64), container.NewStack(logoBackground, container.NewCenter(logoText)))
	buttons := make([]*fuelNavigationButton, len(views))
	for index, view := range views {
		descriptor, _ := calculation.FuelByType(calculation.FuelType(view.mode))
		buttonIndex := index
		buttons[index] = newFuelNavigationButton(view.mode, navigationLabel(descriptor), index == 0, func() {
			for otherIndex := range views {
				buttons[otherIndex].active = otherIndex == buttonIndex
				buttons[otherIndex].Refresh()
				if otherIndex == buttonIndex {
					views[otherIndex].content.Show()
				} else {
					views[otherIndex].content.Hide()
				}
			}
		})
	}
	settingsButton := newCircleIconButton(settingsIconResource(), func() {
		showSettingsDialog(window, settings, func() {
			for _, view := range views {
				view.refreshForSettingsChange()
			}
		})
	})
	navigationDivider := canvas.NewRectangle(appPalette.border)
	navigationDivider.SetMinSize(fyne.NewSize(1, 34))
	navigationObjects := make([]fyne.CanvasObject, 0, len(buttons)*2+5)
	for index, button := range buttons {
		if index > 0 {
			navigationObjects = append(navigationObjects, horizontalGap(6))
		}
		navigationObjects = append(navigationObjects, container.NewGridWrap(fyne.NewSize(136, 62), button))
	}
	historyButton := newCircleIconButton(historyIconResource(), func() { showHistoryDialog(window, history) })
	navigationObjects = append(navigationObjects, horizontalGap(10), container.NewGridWrap(fyne.NewSize(54, 54), historyButton))
	navigationObjects = append(navigationObjects,
		horizontalGap(14),
		container.NewCenter(navigationDivider),
		horizontalGap(14),
		container.NewGridWrap(fyne.NewSize(54, 54), settingsButton),
	)
	navigationActions := container.NewHBox(navigationObjects...)
	barContent := container.NewBorder(nil, nil, logo, navigationActions, nil)
	barFrame := container.NewGridWrap(fyne.NewSize(ui.contentWidth, ui.topBarHeight), barContent)
	background := canvas.NewRectangle(appPalette.railBackground)
	background.SetMinSize(fyne.NewSize(0, ui.topBarHeight))
	separator := canvas.NewRectangle(appPalette.border)
	separator.SetMinSize(fyne.NewSize(0, 1))
	return container.NewBorder(nil, separator, nil, nil, container.NewStack(background, container.NewCenter(barFrame)))
}

func navigationLabel(descriptor calculation.FuelDescriptor) string {
	if descriptor.Fuel == calculation.FuelBriquettes {
		return "Briketts"
	}
	return descriptor.Label
}
