package ui

import "fyne.io/fyne/v2"

type workspaceLayout struct{}

func (layout *workspaceLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) < 2 {
		return
	}
	contentWidth := minFloat(size.Width-48, ui.contentWidth)
	contentWidth = maxFloat(contentWidth, ui.contentWidth)
	contentHeight := minFloat(size.Height-24, ui.workspaceHeight)
	contentHeight = maxFloat(contentHeight, ui.workspaceHeight)
	leftWidth := contentWidth - ui.resultWidth - 42
	startX := (size.Width - contentWidth) / 2
	startY := (size.Height - contentHeight) / 2

	objects[0].Move(fyne.NewPos(startX, startY))
	objects[0].Resize(fyne.NewSize(leftWidth, contentHeight))
	objects[1].Move(fyne.NewPos(startX+leftWidth+42, startY+42))
	objects[1].Resize(fyne.NewSize(ui.resultWidth, ui.resultHeight))
}

func (layout *workspaceLayout) MinSize([]fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(ui.contentWidth, ui.workspaceHeight)
}

type resultValueLayout struct{}

func (layout *resultValueLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) < 2 {
		return
	}
	valueSize := objects[0].MinSize()
	unitSize := objects[1].MinSize()
	objects[0].Move(fyne.NewPos(0, 0))
	objects[0].Resize(valueSize)
	unitX := valueSize.Width + 12
	if unitX+unitSize.Width > size.Width {
		unitX = size.Width - unitSize.Width
	}
	unitY := valueSize.Height - unitSize.Height - 7
	objects[1].Move(fyne.NewPos(unitX, maxFloat(unitY, 0)))
	objects[1].Resize(unitSize)
}

func (layout *resultValueLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) < 2 {
		return fyne.NewSize(1, 1)
	}
	valueSize := objects[0].MinSize()
	unitSize := objects[1].MinSize()
	return fyne.NewSize(valueSize.Width+12+unitSize.Width, maxFloat(valueSize.Height, unitSize.Height))
}

func minFloat(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func maxFloat(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
