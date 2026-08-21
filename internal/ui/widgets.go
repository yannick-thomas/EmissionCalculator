package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type actionButton struct {
	widget.BaseWidget
	text     string
	onTapped func()
	hovered  bool
	pressed  bool
	focused  bool
}

func newActionButton(text string, onTapped func()) *actionButton {
	button := &actionButton{text: text, onTapped: onTapped}
	button.ExtendBaseWidget(button)
	return button
}

func (button *actionButton) Tapped(*fyne.PointEvent) {
	if button.onTapped != nil {
		button.onTapped()
	}
}

func (button *actionButton) MouseIn(*desktop.MouseEvent) {
	button.hovered = true
	button.Refresh()
}

func (button *actionButton) MouseMoved(*desktop.MouseEvent) {}

func (button *actionButton) MouseOut() {
	button.hovered = false
	button.pressed = false
	button.Refresh()
}

func (button *actionButton) MouseDown(event *desktop.MouseEvent) {
	if event.Button == desktop.MouseButtonPrimary {
		button.pressed = true
		button.Refresh()
	}
}

func (button *actionButton) MouseUp(*desktop.MouseEvent) {
	button.pressed = false
	button.Refresh()
}

func (button *actionButton) FocusGained() {
	button.focused = true
	fyne.Do(button.Refresh)
}

func (button *actionButton) FocusLost() {
	button.focused = false
	fyne.Do(button.Refresh)
}

func (button *actionButton) TypedRune(rune) {}

func (button *actionButton) TypedKey(event *fyne.KeyEvent) {
	switch event.Name {
	case fyne.KeyEnter, fyne.KeyReturn, fyne.KeySpace:
		button.Tapped(nil)
	}
}

func (button *actionButton) CreateRenderer() fyne.WidgetRenderer {
	renderer := &actionButtonRenderer{
		button:     button,
		background: canvas.NewRectangle(appPalette.textPrimary),
		label:      canvasText(button.text, 13, appPalette.white, true),
		arrowDisk:  canvas.NewCircle(appPalette.resultSurface),
		arrow:      canvas.NewImageFromResource(arrowIconResource()),
	}
	renderer.background.CornerRadius = 12
	renderer.arrow.FillMode = canvas.ImageFillContain
	renderer.objects = []fyne.CanvasObject{renderer.background, renderer.label, renderer.arrowDisk, renderer.arrow}
	renderer.Refresh()
	return renderer
}

type actionButtonRenderer struct {
	button     *actionButton
	background *canvas.Rectangle
	label      *canvas.Text
	arrowDisk  *canvas.Circle
	arrow      *canvas.Image
	objects    []fyne.CanvasObject
}

func (renderer *actionButtonRenderer) Layout(size fyne.Size) {
	renderer.background.Resize(size)
	labelSize := renderer.label.MinSize()
	renderer.label.Move(fyne.NewPos(18, (size.Height-labelSize.Height)/2))
	renderer.label.Resize(labelSize)
	renderer.arrowDisk.Move(fyne.NewPos(size.Width-38, (size.Height-28)/2))
	renderer.arrowDisk.Resize(fyne.NewSize(28, 28))
	renderer.arrow.Move(fyne.NewPos(size.Width-31, (size.Height-14)/2))
	renderer.arrow.Resize(fyne.NewSize(14, 14))
}

func (renderer *actionButtonRenderer) MinSize() fyne.Size { return fyne.NewSize(178, 46) }

func (renderer *actionButtonRenderer) Refresh() {
	fill := appPalette.textPrimary
	if renderer.button.hovered {
		fill = appPalette.accent
	}
	if renderer.button.pressed {
		fill = appPalette.accentHover
	}
	renderer.background.FillColor = fill
	if renderer.button.focused {
		renderer.background.StrokeColor = appPalette.resultSurface
		renderer.background.StrokeWidth = 2
	} else {
		renderer.background.StrokeColor = color.Transparent
		renderer.background.StrokeWidth = 0
	}
	renderer.label.Text = renderer.button.text
	canvas.Refresh(renderer.background)
	canvas.Refresh(renderer.label)
}

func (renderer *actionButtonRenderer) Objects() []fyne.CanvasObject { return renderer.objects }
func (renderer *actionButtonRenderer) Destroy()                     {}

type fuelNavigationButton struct {
	widget.BaseWidget
	iconKind string
	label    string
	active   bool
	onTapped func()
	hovered  bool
	pressed  bool
	focused  bool
}

func newFuelNavigationButton(iconKind, label string, active bool, onTapped func()) *fuelNavigationButton {
	button := &fuelNavigationButton{iconKind: iconKind, label: label, active: active, onTapped: onTapped}
	button.ExtendBaseWidget(button)
	return button
}

func (button *fuelNavigationButton) Tapped(*fyne.PointEvent) {
	if button.onTapped != nil {
		button.onTapped()
	}
}

func (button *fuelNavigationButton) MouseIn(*desktop.MouseEvent) {
	button.hovered = true
	button.Refresh()
}

func (button *fuelNavigationButton) MouseMoved(*desktop.MouseEvent) {}

func (button *fuelNavigationButton) MouseOut() {
	button.hovered = false
	button.pressed = false
	button.Refresh()
}

func (button *fuelNavigationButton) MouseDown(event *desktop.MouseEvent) {
	if event.Button == desktop.MouseButtonPrimary {
		button.pressed = true
		button.Refresh()
	}
}

func (button *fuelNavigationButton) MouseUp(*desktop.MouseEvent) {
	button.pressed = false
	button.Refresh()
}

func (button *fuelNavigationButton) FocusGained() {
	button.focused = true
	fyne.Do(button.Refresh)
}

func (button *fuelNavigationButton) FocusLost() {
	button.focused = false
	fyne.Do(button.Refresh)
}

func (button *fuelNavigationButton) TypedRune(rune) {}

func (button *fuelNavigationButton) TypedKey(event *fyne.KeyEvent) {
	switch event.Name {
	case fyne.KeyEnter, fyne.KeyReturn, fyne.KeySpace:
		button.Tapped(nil)
	}
}

func (button *fuelNavigationButton) CreateRenderer() fyne.WidgetRenderer {
	renderer := &fuelNavigationRenderer{
		button:     button,
		background: canvas.NewRectangle(color.Transparent),
		icon:       canvas.NewImageFromResource(navigationIconResource(button.iconKind, button.active, false)),
		label:      canvasText(button.label, 10, appPalette.whiteMuted, true),
	}
	renderer.background.CornerRadius = 17
	renderer.icon.FillMode = canvas.ImageFillContain
	renderer.objects = []fyne.CanvasObject{renderer.background, renderer.icon, renderer.label}
	renderer.Refresh()
	return renderer
}

type fuelNavigationRenderer struct {
	button     *fuelNavigationButton
	background *canvas.Rectangle
	icon       *canvas.Image
	label      *canvas.Text
	objects    []fyne.CanvasObject
}

func (renderer *fuelNavigationRenderer) Layout(size fyne.Size) {
	renderer.background.Resize(size)
	renderer.icon.Move(fyne.NewPos((size.Width-22)/2, 11))
	renderer.icon.Resize(fyne.NewSize(22, 22))
	labelSize := renderer.label.MinSize()
	renderer.label.Move(fyne.NewPos((size.Width-labelSize.Width)/2, 42))
	renderer.label.Resize(labelSize)
}

func (renderer *fuelNavigationRenderer) MinSize() fyne.Size { return fyne.NewSize(66, 66) }

func (renderer *fuelNavigationRenderer) Refresh() {
	background := color.Color(color.Transparent)
	textColor := appPalette.whiteMuted
	if renderer.button.hovered && !renderer.button.active {
		background = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 28}
		textColor = appPalette.white
	}
	if renderer.button.pressed && !renderer.button.active {
		background = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 44}
	}
	if renderer.button.active {
		background = appPalette.white
		textColor = appPalette.accent
	}
	renderer.background.FillColor = background
	if renderer.button.focused {
		renderer.background.StrokeColor = appPalette.resultSurface
		renderer.background.StrokeWidth = 2
	} else {
		renderer.background.StrokeColor = color.Transparent
		renderer.background.StrokeWidth = 0
	}
	renderer.icon.Resource = navigationIconResource(renderer.button.iconKind, renderer.button.active, renderer.button.hovered)
	renderer.label.Color = textColor
	canvas.Refresh(renderer.background)
	renderer.icon.Refresh()
	canvas.Refresh(renderer.label)
}

func (renderer *fuelNavigationRenderer) Objects() []fyne.CanvasObject { return renderer.objects }
func (renderer *fuelNavigationRenderer) Destroy()                     {}

type focusEntry struct {
	widget.Entry
	onFocusChanged func(bool)
}

func newFocusEntry() *focusEntry {
	entry := &focusEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (entry *focusEntry) FocusGained() {
	entry.Entry.FocusGained()
	if entry.onFocusChanged != nil {
		fyne.Do(func() { entry.onFocusChanged(true) })
	}
}

func (entry *focusEntry) FocusLost() {
	entry.Entry.FocusLost()
	if entry.onFocusChanged != nil {
		fyne.Do(func() { entry.onFocusChanged(false) })
	}
}

type quantityControl struct {
	content    fyne.CanvasObject
	background *canvas.Rectangle
	unit       *canvas.Text
	focused    bool
	invalid    bool
}

func newQuantityControl(entry *focusEntry, unitValue string) *quantityControl {
	control := &quantityControl{}
	control.unit = canvasText(unitValue, 12, appPalette.accent, true)
	unitBackground := canvas.NewRectangle(appPalette.accentSoft)
	unitBackground.CornerRadius = 8
	unitPill := container.NewGridWrap(fyne.NewSize(76, 36), container.NewStack(unitBackground, container.NewCenter(control.unit)))
	control.background = canvas.NewRectangle(appPalette.surface)
	control.background.CornerRadius = 11
	inputContent := container.NewBorder(nil, nil, nil, unitPill, entry)
	control.content = container.NewStack(control.background, container.NewPadded(inputContent))
	control.refreshBorder()
	entry.onFocusChanged = control.SetFocused
	return control
}

func (control *quantityControl) SetFocused(focused bool) {
	control.focused = focused
	control.refreshBorder()
}

func (control *quantityControl) SetError(invalid bool) {
	control.invalid = invalid
	control.refreshBorder()
}

func (control *quantityControl) refreshBorder() {
	control.background.StrokeColor = appPalette.border
	control.background.StrokeWidth = 1
	if control.focused {
		control.background.StrokeColor = appPalette.accent
		control.background.StrokeWidth = 2
	}
	if control.invalid {
		control.background.StrokeColor = appPalette.error
		control.background.StrokeWidth = 2
	}
	canvas.Refresh(control.background)
}
