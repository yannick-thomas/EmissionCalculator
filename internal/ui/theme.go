package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/gobolditalic"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/font/gofont/gomono"
)

type designTokens struct {
	topBarHeight float32
	contentWidth float32
	pageHeight   float32
	resultHeight float32
}

type palette struct {
	background     color.Color
	railBackground color.Color
	surface        color.Color
	resultSurface  color.Color
	border         color.Color
	textPrimary    color.Color
	textSecondary  color.Color
	textMuted      color.Color
	textDisabled   color.Color
	accent         color.Color
	accentHover    color.Color
	accentSoft     color.Color
	coral          color.Color
	error          color.Color
	success        color.Color
	white          color.Color
	whiteMuted     color.Color
}

var ui = designTokens{
	topBarHeight: 96,
	contentWidth: 1000,
	pageHeight:   1410,
	resultHeight: 540,
}

var appPalette = palette{
	background:     color.NRGBA{R: 0xf4, G: 0xf1, B: 0xe8, A: 255},
	railBackground: color.NRGBA{R: 0x31, G: 0x55, B: 0xd7, A: 255},
	surface:        color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 255},
	resultSurface:  color.NRGBA{R: 0xe1, G: 0xf7, B: 0x7b, A: 255},
	border:         color.NRGBA{R: 0xd4, G: 0xd1, B: 0xc8, A: 255},
	textPrimary:    color.NRGBA{R: 0x18, G: 0x21, B: 0x3d, A: 255},
	textSecondary:  color.NRGBA{R: 0x61, G: 0x66, B: 0x73, A: 255},
	textMuted:      color.NRGBA{R: 0x82, G: 0x86, B: 0x90, A: 255},
	textDisabled:   color.NRGBA{R: 0xab, G: 0xae, B: 0xb5, A: 255},
	accent:         color.NRGBA{R: 0x31, G: 0x55, B: 0xd7, A: 255},
	accentHover:    color.NRGBA{R: 0x24, G: 0x45, B: 0xbd, A: 255},
	accentSoft:     color.NRGBA{R: 0xe7, G: 0xec, B: 0xff, A: 255},
	coral:          color.NRGBA{R: 0xff, G: 0x77, B: 0x5f, A: 255},
	error:          color.NRGBA{R: 0xc7, G: 0x3d, B: 0x46, A: 255},
	success:        color.NRGBA{R: 0x3d, G: 0x7e, B: 0x52, A: 255},
	white:          color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 255},
	whiteMuted:     color.NRGBA{R: 0xd6, G: 0xdd, B: 0xff, A: 255},
}

var (
	regularFont    = fyne.NewStaticResource("go-medium.ttf", gomedium.TTF)
	boldFont       = fyne.NewStaticResource("go-bold.ttf", gobold.TTF)
	italicFont     = fyne.NewStaticResource("go-italic.ttf", goitalic.TTF)
	boldItalicFont = fyne.NewStaticResource("go-bold-italic.ttf", gobolditalic.TTF)
	monoFont       = fyne.NewStaticResource("go-mono.ttf", gomono.TTF)
)

type emissionTheme struct{ fyne.Theme }

func (themeOverride emissionTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary, theme.ColorNameFocus, theme.ColorNameSelection:
		return appPalette.accent
	case theme.ColorNameHover:
		return appPalette.accentHover
	case theme.ColorNameBackground:
		return appPalette.background
	case theme.ColorNameInputBackground:
		return appPalette.background
	case theme.ColorNameInputBorder:
		return color.Transparent
	case theme.ColorNameButton:
		return appPalette.surface
	case theme.ColorNameDisabledButton:
		return appPalette.accentSoft
	case theme.ColorNameDisabled:
		return appPalette.textDisabled
	case theme.ColorNameForeground:
		return appPalette.textPrimary
	case theme.ColorNamePlaceHolder:
		return appPalette.textMuted
	}
	return themeOverride.Theme.Color(name, variant)
}

func (themeOverride emissionTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 20
	case theme.SizeNameInputBorder:
		return 0
	}
	return themeOverride.Theme.Size(name)
}

func (themeOverride emissionTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return monoFont
	}
	if style.Bold && style.Italic {
		return boldItalicFont
	}
	if style.Bold {
		return boldFont
	}
	if style.Italic {
		return italicFont
	}
	if style.Symbol {
		return themeOverride.Theme.Font(style)
	}
	return regularFont
}
