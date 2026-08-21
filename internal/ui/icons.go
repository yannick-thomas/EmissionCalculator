package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
)

func navigationIconResource(kind string, active, highlighted bool) fyne.Resource {
	stroke := "#D6DDFF"
	if highlighted && !active {
		stroke = "#FFFFFF"
	}
	if active {
		stroke = "#3155D7"
	}
	var paths string
	if kind == "oil" {
		paths = `<path d="M12 3.5S6.8 9.5 6.8 14a5.2 5.2 0 0 0 10.4 0C17.2 9.5 12 3.5 12 3.5Z"/><path d="M9.5 15c.3 1.1 1.2 1.7 2.5 1.7"/>`
	} else {
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

func resultPanelResource(active bool) fyne.Resource {
	fill := "#EDEEE7"
	if active {
		fill = "#E1F77B"
	}
	coral := ""
	if active {
		coral = `<path d="M0 430A90 90 0 0 1 90 520H0Z" fill="#FF775F"/>`
	}
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1000 520"><defs><clipPath id="card"><rect width="990" height="510" rx="34"/></clipPath></defs><g clip-path="url(#card)"><rect width="990" height="510" rx="34" fill="%s"/><g fill="none" stroke="#62772E" stroke-opacity="0.13"><circle cx="930" cy="5" r="205" stroke-width="2"/><circle cx="930" cy="5" r="250" stroke-width="18" stroke-opacity="0.07"/></g>%s</g></svg>`, fill, coral)
	return fyne.NewStaticResource(fmt.Sprintf("result-%t.svg", active), []byte(svg))
}

func resultShadowResource() fyne.Resource {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1000 520"><rect x="20" y="16" width="980" height="504" rx="34" fill="#18213D" fill-opacity="0.11"/></svg>`
	return fyne.NewStaticResource("result-shadow.svg", []byte(svg))
}
