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

func resultPanelResource(active bool) fyne.Resource {
	fill := "#EDEEE7"
	if active {
		fill = "#E1F77B"
	}
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 420 380"><path d="M32 0H406Q420 0 420 14V348Q420 380 388 380H32Q0 380 0 348V32Q0 0 32 0Z" fill="%s"/></svg>`, fill)
	return fyne.NewStaticResource(fmt.Sprintf("result-%t.svg", active), []byte(svg))
}

func resultShadowResource() fyne.Resource {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 420 380"><path d="M32 0H406Q420 0 420 14V348Q420 380 388 380H32Q0 380 0 348V32Q0 0 32 0Z" fill="#18213D" fill-opacity="0.10"/></svg>`
	return fyne.NewStaticResource("result-shadow.svg", []byte(svg))
}
