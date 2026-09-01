package ui

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
)

func TestResponsivePageSwitchesFromColumnsToStack(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	header := canvas.NewRectangle(color.Transparent)
	form := canvas.NewRectangle(color.Transparent)
	form.SetMinSize(fyne.NewSize(460, 570))
	result := canvas.NewRectangle(color.Transparent)
	result.SetMinSize(fyne.NewSize(460, 570))
	page := newResponsivePage(header, form, result)
	renderer := page.CreateRenderer()

	page.Resize(fyne.NewSize(1200, 800))
	renderer.Layout(page.Size())
	if page.compactMode || form.Position().Y != result.Position().Y || result.Position().X <= form.Position().X {
		t.Fatalf("expected side-by-side desktop layout, form=%v result=%v", form.Position(), result.Position())
	}

	page.Resize(fyne.NewSize(700, 800))
	renderer.Layout(page.Size())
	if !page.compactMode || result.Position().Y <= form.Position().Y || result.Position().X != form.Position().X {
		t.Fatalf("expected vertically stacked compact layout, form=%v result=%v", form.Position(), result.Position())
	}
	if page.Size().Height <= 800 {
		t.Fatalf("expected compact content to grow vertically for scrolling, got %v", page.Size())
	}
}

func TestResponsiveNavigationUsesSecondRowWhenCompact(t *testing.T) {
	logo := canvas.NewRectangle(color.Transparent)
	fuels := canvas.NewRectangle(color.Transparent)
	tools := canvas.NewRectangle(color.Transparent)
	tools.SetMinSize(fyne.NewSize(140, 62))
	layout := &responsiveNavigationLayout{}
	objects := []fyne.CanvasObject{logo, fuels, tools}

	layout.Layout(objects, fyne.NewSize(700, ui.topBarHeight))
	if fuels.Position().Y <= logo.Position().Y || tools.Position().Y != 5 {
		t.Fatalf("expected compact two-row navigation: logo=%v fuels=%v tools=%v", logo.Position(), fuels.Position(), tools.Position())
	}
	if fuels.Size().Height < 44 || tools.Size().Height != 44 {
		t.Fatalf("expected compact navigation controls to remain fully visible: fuels=%v tools=%v", fuels.Size(), tools.Size())
	}
	layout.Layout(objects, fyne.NewSize(1200, ui.topBarHeight))
	if fuels.Position().Y != tools.Position().Y {
		t.Fatalf("expected one-row desktop navigation: fuels=%v tools=%v", fuels.Position(), tools.Position())
	}
}

func TestResponsiveDialogSizeNeverExceedsParent(t *testing.T) {
	actual := responsiveDialogSize(fyne.NewSize(700, 500), fyne.NewSize(920, 650))
	if actual != (fyne.NewSize(652, 428)) {
		t.Fatalf("unexpected fitted dialog size: %v", actual)
	}
	tiny := responsiveDialogSize(fyne.NewSize(300, 200), fyne.NewSize(760, 680))
	if tiny.Width > 300 || tiny.Height > 200 {
		t.Fatalf("dialog must not exceed a tiny parent: %v", tiny)
	}
}

func TestSettingsPanelCanShrinkAndScroll(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	panel := newSettingsPanel(newSettingsStore(app.Preferences()), nil)
	minimum := panel.content.MinSize()
	if minimum.Width > 360 || minimum.Height > 320 {
		t.Fatalf("settings overlay has an inflexible minimum: %v", minimum)
	}
}
