package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type forcedVariantCustomTheme struct {
	fyne.Theme
	variant fyne.ThemeVariant
}

func (t *forcedVariantCustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameInputBorder {
		if t.variant == theme.VariantLight {
			return color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black for light theme
		}
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // White for dark theme
	}
	return t.Theme.Color(name, t.variant)
}

func (t *forcedVariantCustomTheme) Variant() fyne.ThemeVariant {
	return t.variant
}

func SetDarkTheme(myApp fyne.App) {
	myApp.Settings().SetTheme(
		&forcedVariantCustomTheme{
			Theme:   theme.DefaultTheme(),
			variant: theme.VariantDark,
		},
	)
}

func SetLightTheme(myApp fyne.App) {
	myApp.Settings().SetTheme(
		&forcedVariantCustomTheme{
			Theme:   theme.DefaultTheme(),
			variant: theme.VariantLight,
		},
	)
}
