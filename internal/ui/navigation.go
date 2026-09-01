package ui

import (
	"emissioncalculator/internal/calculation"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func buildSharedNavigation(window fyne.Window, settings *settingsStore, history *historyController, fuelMenu *fuelMenuStore, views []*referenceView, onFuelAdded func()) fyne.CanvasObject {
	logoBackground := canvas.NewImageFromResource(logoPanelResource())
	logoBackground.FillMode = canvas.ImageFillStretch
	logoText := canvasText("e°", 22, appPalette.textPrimary, true)
	logo := container.NewGridWrap(fyne.NewSize(44, 44), container.NewStack(logoBackground, container.NewCenter(logoText)))
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
	fuelObjects := make([]fyne.CanvasObject, 0, len(buttons)*2)
	for index, button := range buttons {
		if index > 0 {
			fuelObjects = append(fuelObjects, horizontalGap(6))
		}
		fuelObjects = append(fuelObjects, container.NewGridWrap(fyne.NewSize(124, 44), button))
	}
	addFuelButton := widget.NewButton("+", func() {
		showAddFuelDialog(window, fuelMenu, onFuelAdded)
	})
	if len(fuelMenu.availableDescriptors()) == 0 {
		addFuelButton.Disable()
	}
	fuelObjects = append(fuelObjects, horizontalGap(6), container.NewGridWrap(fyne.NewSize(44, 44), addFuelButton))
	fuelsScroll := container.NewHScroll(container.NewHBox(fuelObjects...))
	historyButton := newCircleIconButton(historyIconResource(), func() { showHistoryDialog(window, history) })
	toolActions := container.NewHBox(
		container.NewGridWrap(fyne.NewSize(42, 42), historyButton),
		horizontalGap(8),
		container.NewCenter(navigationDivider),
		horizontalGap(8),
		container.NewGridWrap(fyne.NewSize(42, 42), settingsButton),
	)
	barContent := container.New(&responsiveNavigationLayout{}, logo, fuelsScroll, toolActions)
	background := canvas.NewRectangle(appPalette.railBackground)
	background.SetMinSize(fyne.NewSize(0, ui.topBarHeight))
	separator := canvas.NewRectangle(appPalette.border)
	separator.SetMinSize(fyne.NewSize(0, 1))
	return container.NewBorder(nil, separator, nil, nil, container.NewStack(background, container.NewPadded(barContent)))
}

func showAddFuelDialog(window fyne.Window, fuelMenu *fuelMenuStore, onFuelAdded func()) {
	available := fuelMenu.availableDescriptors()
	if len(available) == 0 {
		return
	}

	options := make([]string, len(available))
	byLabel := make(map[string]calculation.FuelType, len(available))
	for index, descriptor := range available {
		options[index] = navigationLabel(descriptor)
		byLabel[options[index]] = descriptor.Fuel
	}
	selector := widget.NewSelect(options, nil)
	selector.SetSelected(options[0])
	content := container.NewVBox(
		canvasText("Brennstoff auswählen", 15, appPalette.textPrimary, true),
		verticalGap(8),
		canvasText("Der Brennstoff wird dem Menü hinzugefügt.", 14, appPalette.textSecondary, false),
		verticalGap(16),
		selector,
	)
	dialog.ShowCustomConfirm("Brennstoff hinzufügen", "Hinzufügen", "Abbrechen", content, func(confirmed bool) {
		if !confirmed {
			return
		}
		fuelMenu.enable(byLabel[selector.Selected])
		onFuelAdded()
	}, window)
}

func navigationLabel(descriptor calculation.FuelDescriptor) string {
	if descriptor.Fuel == calculation.FuelBriquettes {
		return "Briketts"
	}
	return descriptor.Label
}
