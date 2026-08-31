package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

const (
	pageBreakpoint = float32(900)
	pageMaxWidth   = float32(1160)
	pageMinWidth   = float32(500)
	pageInset      = float32(20)
	pageGap        = float32(24)
	headerHeight   = float32(126)
)

// responsivePage switches the calculator from two columns to a vertical flow when
// the available width can no longer carry both columns without clipping.
type responsivePage struct {
	widget.BaseWidget
	header      fyne.CanvasObject
	form        fyne.CanvasObject
	result      fyne.CanvasObject
	lastWidth   float32
	compactMode bool
}

func newResponsivePage(header, form, result fyne.CanvasObject) *responsivePage {
	page := &responsivePage{header: header, form: form, result: result, lastWidth: pageMinWidth}
	page.ExtendBaseWidget(page)
	return page
}

func (page *responsivePage) CreateRenderer() fyne.WidgetRenderer {
	return &responsivePageRenderer{
		page:    page,
		objects: []fyne.CanvasObject{page.header, page.form, page.result},
	}
}

func (page *responsivePage) MinSize() fyne.Size {
	return fyne.NewSize(pageMinWidth, page.desiredHeight(page.lastWidth))
}

func (page *responsivePage) Resize(size fyne.Size) {
	width := fyne.Max(size.Width, pageMinWidth)
	page.lastWidth = width
	page.compactMode = width < pageBreakpoint
	page.BaseWidget.Resize(fyne.NewSize(width, page.desiredHeight(width)))
}

func (page *responsivePage) desiredHeight(width float32) float32 {
	formHeight := fyne.Max(page.form.MinSize().Height, ui.colHeight)
	resultHeight := fyne.Max(page.result.MinSize().Height, ui.colHeight)
	if width >= pageBreakpoint {
		return headerHeight + 36 + fyne.Max(formHeight, resultHeight) + 36
	}
	return headerHeight + pageGap + formHeight + pageGap + resultHeight + 36
}

type responsivePageRenderer struct {
	page    *responsivePage
	objects []fyne.CanvasObject
}

func (renderer *responsivePageRenderer) Layout(size fyne.Size) {
	page := renderer.page
	contentWidth := fyne.Min(size.Width-pageInset*2, pageMaxWidth)
	contentWidth = fyne.Max(contentWidth, pageMinWidth-pageInset*2)
	left := (size.Width - contentWidth) / 2
	page.header.Move(fyne.NewPos(left, 0))
	page.header.Resize(fyne.NewSize(contentWidth, headerHeight))

	formHeight := fyne.Max(page.form.MinSize().Height, ui.colHeight)
	resultHeight := fyne.Max(page.result.MinSize().Height, ui.colHeight)
	y := headerHeight + 36
	if size.Width >= pageBreakpoint {
		formWidth := (contentWidth - pageGap) * 0.45
		resultWidth := contentWidth - pageGap - formWidth
		page.form.Move(fyne.NewPos(left, y))
		page.form.Resize(fyne.NewSize(formWidth, formHeight))
		page.result.Move(fyne.NewPos(left+formWidth+pageGap, y))
		page.result.Resize(fyne.NewSize(resultWidth, resultHeight))
		return
	}

	page.form.Move(fyne.NewPos(left, y))
	page.form.Resize(fyne.NewSize(contentWidth, formHeight))
	y += formHeight + pageGap
	page.result.Move(fyne.NewPos(left, y))
	page.result.Resize(fyne.NewSize(contentWidth, resultHeight))
}

func (renderer *responsivePageRenderer) MinSize() fyne.Size { return renderer.page.MinSize() }

func (renderer *responsivePageRenderer) Refresh() {
	for _, object := range renderer.objects {
		canvas.Refresh(object)
	}
}

func (renderer *responsivePageRenderer) Objects() []fyne.CanvasObject { return renderer.objects }
func (renderer *responsivePageRenderer) Destroy()                     {}

type responsiveHeaderLayout struct{}

func (layout *responsiveHeaderLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(360, headerHeight)
}

type responsiveNavigationLayout struct{}

func (layout *responsiveNavigationLayout) MinSize([]fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(360, ui.topBarHeight)
}

func (layout *responsiveNavigationLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) != 3 {
		return
	}
	logo, fuels, tools := objects[0], objects[1], objects[2]
	if size.Width >= pageBreakpoint {
		logo.Move(fyne.NewPos(0, (size.Height-64)/2))
		logo.Resize(fyne.NewSize(64, 64))
		toolsWidth := fyne.Max(tools.MinSize().Width, 140)
		tools.Move(fyne.NewPos(size.Width-toolsWidth, (size.Height-62)/2))
		tools.Resize(fyne.NewSize(toolsWidth, 62))
		fuels.Move(fyne.NewPos(80, (size.Height-62)/2))
		fuels.Resize(fyne.NewSize(size.Width-toolsWidth-96, 62))
		return
	}

	logo.Move(fyne.NewPos(0, 4))
	logo.Resize(fyne.NewSize(64, 64))
	toolsWidth := fyne.Min(size.Width-80, fyne.Max(tools.MinSize().Width, 140))
	tools.Move(fyne.NewPos(size.Width-toolsWidth, 5))
	tools.Resize(fyne.NewSize(toolsWidth, 62))
	fuels.Move(fyne.NewPos(0, 70))
	fuels.Resize(fyne.NewSize(size.Width, size.Height-74))
}

func (layout *responsiveHeaderLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) != 4 {
		return
	}
	brand, status, actions, separator := objects[0], objects[1], objects[2], objects[3]
	separator.Move(fyne.NewPos(0, size.Height-1))
	separator.Resize(fyne.NewSize(size.Width, 1))
	if size.Width >= pageBreakpoint {
		rowHeight := float32(62)
		y := (size.Height - rowHeight) / 2
		brandWidth := fyne.Max(brand.MinSize().Width, 180)
		statusWidth := float32(170)
		brand.Move(fyne.NewPos(0, y))
		brand.Resize(fyne.NewSize(brandWidth, rowHeight))
		status.Move(fyne.NewPos(brandWidth+12, y))
		status.Resize(fyne.NewSize(statusWidth, rowHeight))
		actions.Move(fyne.NewPos(brandWidth+statusWidth+24, y))
		actions.Resize(fyne.NewSize(size.Width-brandWidth-statusWidth-24, rowHeight))
		return
	}

	brandWidth := fyne.Min(size.Width*0.55, fyne.Max(brand.MinSize().Width, 180))
	brand.Move(fyne.NewPos(0, 4))
	brand.Resize(fyne.NewSize(brandWidth, 52))
	status.Move(fyne.NewPos(brandWidth, 4))
	status.Resize(fyne.NewSize(size.Width-brandWidth, 52))
	actions.Move(fyne.NewPos(0, 64))
	actions.Resize(fyne.NewSize(size.Width, 54))
}

// centeredCardLayout keeps a modal panel at a comfortable maximum size while
// allowing its scrollable child to shrink with the parent canvas.
type centeredCardLayout struct {
	preferred fyne.Size
	minimum   fyne.Size
	margin    float32
}

func (layout *centeredCardLayout) MinSize([]fyne.CanvasObject) fyne.Size { return layout.minimum }

func (layout *centeredCardLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) == 0 {
		return
	}
	available := fyne.NewSize(fyne.Max(1, size.Width-layout.margin*2), fyne.Max(1, size.Height-layout.margin*2))
	cardSize := layout.preferred.Min(available)
	objects[0].Move(fyne.NewPos((size.Width-cardSize.Width)/2, (size.Height-cardSize.Height)/2))
	objects[0].Resize(cardSize)
}
