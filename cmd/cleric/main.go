package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pcarion/cleric/cmd/cleric/ui"
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
	myWindow := myApp.NewWindow("Claude Mcp Servers")

	// Get window size and calculate desired width
	windowWidth := float32(800) // default reasonable width

	// Set window size (using the calculated width and a reasonable height)
	myWindow.Resize(fyne.NewSize(windowWidth, 400))

	// Center the window on screen
	myWindow.CenterOnScreen()

	mcpServersList := ui.NewMcpServersList(myWindow)

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			myApp.Settings().SetTheme(&forcedVariantCustomTheme{Theme: theme.DefaultTheme(), variant: theme.VariantDark})
		}),
		widget.NewButton("Light", func() {
			myApp.Settings().SetTheme(&forcedVariantCustomTheme{Theme: theme.DefaultTheme(), variant: theme.VariantLight})
		}),
	)

	buttons := container.New(layout.NewGridLayout(4),
		widget.NewButton("Add new MCPserver", func() {
			mcpServersList.AddMcpServer()
		}),
		// separator
		widget.NewSeparator(),
		widget.NewSeparator(),
		themes,
	)

	content := container.NewBorder(nil, buttons, nil, nil, mcpServersList.GetList())

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
