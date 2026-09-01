package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
)

func navigationIconResource(kind string, active, highlighted bool) fyne.Resource {
	stroke := "#616673"
	if highlighted && !active {
		stroke = "#3155D7"
	}
	if active {
		stroke = "#18213D"
	}
	var paths string
	switch kind {
	case modeOil:
		paths = `<path d="M12 3.5S6.8 9.5 6.8 14a5.2 5.2 0 0 0 10.4 0C17.2 9.5 12 3.5 12 3.5Z"/><path d="M9.5 15c.3 1.1 1.2 1.7 2.5 1.7"/>`
	case modeNaturalGas:
		paths = `<path d="M13.5 3.5c.4 3-2.8 4.7-2 7.4.4 1.4 1.7 2 2.5 3.2.9 1.3.3 3.5-1.9 4.4-3.2 1.3-6.6-.9-6.6-4.5 0-3.4 2.4-5.7 5.7-9.8-.2 2.5.9 3.7 1.9 4.2-.2-2.2 1-3.6 3.9-4.9Z"/>`
	case modeLPG:
		paths = `<path d="M8 5.5h8v15H8z"/><path d="M10 5.5V3.8h4v1.7M10.5 10h3M10.5 13h3"/>`
	default:
		paths = `<path d="M5 8.5 9 4.5l4 4-4 4-4-4Z"/><path d="m11 15.5 4-4 4 4-4 4-4-4Z"/>`
	}
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><g fill="none" stroke="%s" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">%s</g></svg>`, stroke, paths)
	return fyne.NewStaticResource(fmt.Sprintf("%s-%t-%t.svg", kind, active, highlighted), []byte(svg))
}

func arrowIconResource() fyne.Resource {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M5 12h13m-5-5 5 5-5 5" fill="none" stroke="#18213D" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>`
	return fyne.NewStaticResource("arrow-right.svg", []byte(svg))
}

func logoPanelResource() fyne.Resource {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64"><path d="M12 0h28c13.3 0 24 10.7 24 24v16c0 13.3-10.7 24-24 24H12C5.4 64 0 58.6 0 52V12C0 5.4 5.4 0 12 0Z" fill="#E1F77B"/></svg>`
	return fyne.NewStaticResource("logo-panel.svg", []byte(svg))
}

func printIconResource() fyne.Resource {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><g fill="none" stroke="#18213D" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round"><path d="M7 8V3.8h10V8"/><path d="M6 17H4.8A2.8 2.8 0 0 1 2 14.2v-3.4A2.8 2.8 0 0 1 4.8 8h14.4a2.8 2.8 0 0 1 2.8 2.8v3.4a2.8 2.8 0 0 1-2.8 2.8H18"/><path d="M6 13.5h12v6.7H6z"/><path d="M18.5 11.2h.1"/></g></svg>`
	return fyne.NewStaticResource("print.svg", []byte(svg))
}

func saveIconResource() fyne.Resource {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><g fill="none" stroke="#18213D" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round"><path d="M5 3.8h11.2L20 7.6V19a1.2 1.2 0 0 1-1.2 1.2H5A1.2 1.2 0 0 1 3.8 19V5A1.2 1.2 0 0 1 5 3.8Z"/><path d="M7.5 3.8v5h8v-5"/><path d="M7.5 14h9v6.2h-9z"/></g></svg>`
	return fyne.NewStaticResource("save.svg", []byte(svg))
}

func scenarioIconResource() fyne.Resource {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><g fill="none" stroke="#18213D" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round"><path d="M4 18V9"/><path d="M10 18V5"/><path d="M16 18v-7"/><path d="M20 18V11"/><path d="M3 20h18"/></g></svg>`
	return fyne.NewStaticResource("scenario.svg", []byte(svg))
}

func traceIconResource() fyne.Resource {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><g fill="none" stroke="#18213D" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round"><path d="M5 4h14v16H5z"/><path d="M8 8h8M8 12h8M8 16h5"/></g></svg>`
	return fyne.NewStaticResource("trace.svg", []byte(svg))
}

func historyIconResource() fyne.Resource {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><g fill="none" stroke="#18213D" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round"><path d="M4 6h6l2 2h8v11H4z"/><path d="M8 12h8M8 15h5"/></g></svg>`
	return fyne.NewStaticResource("history.svg", []byte(svg))
}

func settingsIconResource() fyne.Resource {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><g fill="none" stroke="#18213D" stroke-width="1.9" stroke-linecap="round"><line x1="3" y1="8" x2="6.2" y2="8"/><circle cx="9" cy="8" r="2.8"/><line x1="11.8" y1="8" x2="21" y2="8"/><line x1="3" y1="16" x2="12.2" y2="16"/><circle cx="15" cy="16" r="2.8"/><line x1="17.8" y1="16" x2="21" y2="16"/></g></svg>`
	return fyne.NewStaticResource("settings.svg", []byte(svg))
}

func resultPanelResource(active bool) fyne.Resource {
	fill := "#FFFFFF"
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1000 520"><rect width="990" height="510" rx="6" fill="%s"/></svg>`, fill)
	return fyne.NewStaticResource(fmt.Sprintf("result-%t.svg", active), []byte(svg))
}

func resultShadowResource() fyne.Resource {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1000 520"><rect x="20" y="16" width="980" height="504" rx="34" fill="#18213D" fill-opacity="0.11"/></svg>`
	return fyne.NewStaticResource("result-shadow.svg", []byte(svg))
}
