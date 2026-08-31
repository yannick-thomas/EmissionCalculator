package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

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
	return canvasText(" ", 14, appPalette.textSecondary, true)
}

func setStatus(status *canvas.Text, message string, colorValue color.Color) {
	status.Text = message
	status.Color = colorValue
	status.Refresh()
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

func resultTextSize(value string) float32 {
	switch {
	case len(value) <= 10:
		return 72
	case len(value) <= 13:
		return 54
	default:
		return 44
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

func refreshTexts(texts ...*canvas.Text) {
	for _, text := range texts {
		text.Refresh()
	}
}
