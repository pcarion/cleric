package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var data = []string{"a", "string", "list"}

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
	// // Change checkbox background color only when selected in dark mode
	// if name == theme.ColorNameInputBackground && t.variant == theme.VariantDark {
	// 	if variant == theme.VariantDark {
	// 		return color.RGBA{R: 255, G: 255, B: 0, A: 255} // Bright yellow for selected state
	// 	}
	// 	return t.Theme.Color(name, t.variant) // Default color for unselected state
	// }
	return t.Theme.Color(name, t.variant)
}

func (t *forcedVariantCustomTheme) Variant() fyne.ThemeVariant {
	return t.variant
}

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(&forcedVariantCustomTheme{Theme: theme.DefaultTheme(), variant: theme.VariantDark})
	myWindow := myApp.NewWindow("List Widget")

	// Get window size and calculate desired width
	windowWidth := float32(800) // default reasonable width

	// Set window size (using the calculated width and a reasonable height)
	myWindow.Resize(fyne.NewSize(windowWidth, 400))

	// Center the window on screen
	myWindow.CenterOnScreen()

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			check := widget.NewCheck("", func(bool) {})
			return container.NewBorder(
				nil, nil,
				widget.NewLabel("template"),
				check,
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			cont := o.(*fyne.Container)
			label := cont.Objects[0].(*widget.Label)

			label.SetText(data[i])
		})

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			myApp.Settings().SetTheme(&forcedVariantCustomTheme{Theme: theme.DefaultTheme(), variant: theme.VariantDark})
		}),
		widget.NewButton("Light", func() {
			myApp.Settings().SetTheme(&forcedVariantCustomTheme{Theme: theme.DefaultTheme(), variant: theme.VariantLight})
		}),
	)

	buttons := container.New(layout.NewGridLayout(3),
		widget.NewButton("Save", func() {}),
		widget.NewButton("Cancel", func() {}),
		themes,
	)

	content := container.NewBorder(nil, buttons, nil, nil, list)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
