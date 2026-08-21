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
		label:      canvasText(button.text, 17, appPalette.white, true),
		arrowDisk:  canvas.NewCircle(appPalette.resultSurface),
		arrow:      canvas.NewImageFromResource(arrowIconResource()),
	}
	renderer.background.CornerRadius = 34
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
	renderer.label.Move(fyne.NewPos(28, (size.Height-labelSize.Height)/2))
	renderer.label.Resize(labelSize)
	renderer.arrowDisk.Move(fyne.NewPos(size.Width-58, (size.Height-42)/2))
	renderer.arrowDisk.Resize(fyne.NewSize(42, 42))
	renderer.arrow.Move(fyne.NewPos(size.Width-48, (size.Height-22)/2))
	renderer.arrow.Resize(fyne.NewSize(22, 22))
}

func (renderer *actionButtonRenderer) MinSize() fyne.Size { return fyne.NewSize(248, 68) }

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
		label:      canvasText(button.label, 16, appPalette.textSecondary, true),
	}
	renderer.background.CornerRadius = 18
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
	renderer.icon.Move(fyne.NewPos(20, (size.Height-24)/2))
	renderer.icon.Resize(fyne.NewSize(24, 24))
	labelSize := renderer.label.MinSize()
	renderer.label.Move(fyne.NewPos(54, (size.Height-labelSize.Height)/2))
	renderer.label.Resize(labelSize)
}

func (renderer *fuelNavigationRenderer) MinSize() fyne.Size { return fyne.NewSize(136, 62) }

func (renderer *fuelNavigationRenderer) Refresh() {
	background := color.Color(color.Transparent)
	textColor := appPalette.textSecondary
	if renderer.button.hovered && !renderer.button.active {
		background = appPalette.accentSoft
		textColor = appPalette.accent
	}
	if renderer.button.pressed && !renderer.button.active {
		background = appPalette.accentSoft
	}
	if renderer.button.active {
		background = appPalette.resultSurface
		textColor = appPalette.textPrimary
	}
	renderer.background.FillColor = background
	if renderer.button.focused {
		renderer.background.StrokeColor = appPalette.accent
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

type circleIconButton struct {
	widget.BaseWidget
	icon     fyne.Resource
	onTapped func()
	hovered  bool
	pressed  bool
	focused  bool
	disabled bool
}

func newCircleIconButton(icon fyne.Resource, onTapped func()) *circleIconButton {
	button := &circleIconButton{icon: icon, onTapped: onTapped}
	button.ExtendBaseWidget(button)
	return button
}

func (button *circleIconButton) Tapped(*fyne.PointEvent) {
	if !button.disabled && button.onTapped != nil {
		button.onTapped()
	}
}

func (button *circleIconButton) MouseIn(*desktop.MouseEvent) {
	button.hovered = true
	button.Refresh()
}

func (button *circleIconButton) MouseMoved(*desktop.MouseEvent) {}

func (button *circleIconButton) MouseOut() {
	button.hovered = false
	button.pressed = false
	button.Refresh()
}

func (button *circleIconButton) MouseDown(event *desktop.MouseEvent) {
	if !button.disabled && event.Button == desktop.MouseButtonPrimary {
		button.pressed = true
		button.Refresh()
	}
}

func (button *circleIconButton) MouseUp(*desktop.MouseEvent) {
	button.pressed = false
	button.Refresh()
}

func (button *circleIconButton) FocusGained() {
	button.focused = true
	fyne.Do(button.Refresh)
}

func (button *circleIconButton) FocusLost() {
	button.focused = false
	fyne.Do(button.Refresh)
}

func (button *circleIconButton) TypedRune(rune) {}

func (button *circleIconButton) TypedKey(event *fyne.KeyEvent) {
	if event.Name == fyne.KeyEnter || event.Name == fyne.KeyReturn || event.Name == fyne.KeySpace {
		button.Tapped(nil)
	}
}

func (button *circleIconButton) Disable() {
	button.disabled = true
	button.Refresh()
}

func (button *circleIconButton) Enable() {
	button.disabled = false
	button.Refresh()
}

func (button *circleIconButton) Disabled() bool { return button.disabled }

func (button *circleIconButton) CreateRenderer() fyne.WidgetRenderer {
	background := canvas.NewCircle(color.Transparent)
	background.StrokeColor = appPalette.border
	background.StrokeWidth = 1
	icon := canvas.NewImageFromResource(button.icon)
	icon.FillMode = canvas.ImageFillContain
	renderer := &circleIconButtonRenderer{button: button, background: background, icon: icon}
	renderer.objects = []fyne.CanvasObject{background, icon}
	renderer.Refresh()
	return renderer
}

type circleIconButtonRenderer struct {
	button     *circleIconButton
	background *canvas.Circle
	icon       *canvas.Image
	objects    []fyne.CanvasObject
}

func (renderer *circleIconButtonRenderer) Layout(size fyne.Size) {
	renderer.background.Resize(size)
	renderer.icon.Move(fyne.NewPos((size.Width-24)/2, (size.Height-24)/2))
	renderer.icon.Resize(fyne.NewSize(24, 24))
}

func (renderer *circleIconButtonRenderer) MinSize() fyne.Size { return fyne.NewSize(54, 54) }

func (renderer *circleIconButtonRenderer) Refresh() {
	fill := color.Color(color.Transparent)
	if renderer.button.hovered && !renderer.button.disabled {
		fill = appPalette.surface
	}
	if renderer.button.pressed && !renderer.button.disabled {
		fill = appPalette.accentSoft
	}
	renderer.background.FillColor = fill
	renderer.background.StrokeColor = appPalette.border
	if renderer.button.focused {
		renderer.background.StrokeColor = appPalette.accent
		renderer.background.StrokeWidth = 2
	} else {
		renderer.background.StrokeWidth = 1
	}
	renderer.icon.Translucency = 0
	if renderer.button.disabled {
		renderer.icon.Translucency = 0.55
	}
	canvas.Refresh(renderer.background)
	renderer.icon.Refresh()
}

func (renderer *circleIconButtonRenderer) Objects() []fyne.CanvasObject { return renderer.objects }
func (renderer *circleIconButtonRenderer) Destroy()                     {}

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
	content   fyne.CanvasObject
	underline *canvas.Rectangle
	unit      *canvas.Text
	focused   bool
	invalid   bool
}

func newQuantityControl(entry *focusEntry, unitValue string) *quantityControl {
	control := &quantityControl{}
	entry.TextStyle = fyne.TextStyle{Bold: true}
	entry.Scroll = fyne.ScrollNone
	control.unit = canvasText(unitValue, 18, appPalette.accent, true)
	unitFrame := container.NewGridWrap(fyne.NewSize(110, 54), container.NewCenter(control.unit))
	inputContent := container.NewBorder(nil, nil, nil, unitFrame, entry)
	control.underline = canvas.NewRectangle(appPalette.textPrimary)
	control.underline.SetMinSize(fyne.NewSize(1, 3))
	control.content = container.NewVBox(
		container.NewGridWrap(fyne.NewSize(ui.formColWidth, 58), inputContent),
		control.underline,
	)
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
	control.underline.FillColor = appPalette.textPrimary
	if control.focused {
		control.underline.FillColor = appPalette.accent
	}
	if control.invalid {
		control.underline.FillColor = appPalette.error
	}
	canvas.Refresh(control.underline)
}
