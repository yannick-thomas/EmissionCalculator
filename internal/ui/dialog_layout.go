package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

var dialogMargin = fyne.NewSize(48, 72)

// showResponsiveDialog gives dialogs a useful preferred size, caps them to the
// parent window and keeps all content reachable in both directions.
func showResponsiveDialog(title, dismiss string, content fyne.CanvasObject, parent fyne.Window, preferred fyne.Size) *dialog.CustomDialog {
	scroll := container.NewScroll(content)
	preferred = responsiveDialogSize(parent.Canvas().Size(), preferred)
	scroll.SetMinSize(fyne.NewSize(fyne.Min(320, preferred.Width), fyne.Min(240, preferred.Height)))
	custom := dialog.NewCustom(title, dismiss, scroll, parent)
	custom.Resize(preferred)
	custom.Show()
	return custom
}

func showResponsiveDialogWithoutButtons(title string, content fyne.CanvasObject, parent fyne.Window, preferred fyne.Size) *dialog.CustomDialog {
	scroll := container.NewScroll(content)
	preferred = responsiveDialogSize(parent.Canvas().Size(), preferred)
	scroll.SetMinSize(fyne.NewSize(fyne.Min(320, preferred.Width), fyne.Min(240, preferred.Height)))
	custom := dialog.NewCustomWithoutButtons(title, scroll, parent)
	custom.Resize(preferred)
	custom.Show()
	return custom
}

func responsiveDialogSize(canvasSize, preferred fyne.Size) fyne.Size {
	if canvasSize.IsZero() {
		return preferred
	}
	available := fyne.NewSize(
		fyne.Max(1, canvasSize.Width-dialogMargin.Width),
		fyne.Max(1, canvasSize.Height-dialogMargin.Height),
	)
	return preferred.Min(available)
}
