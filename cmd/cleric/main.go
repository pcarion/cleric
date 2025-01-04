package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

var data = []string{"a", "string", "list"}

type customTheme struct {
	fyne.Theme
}

func (t *customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameInputBorder {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	return theme.DefaultTheme().Color(name, variant)
}

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(&customTheme{theme.DefaultTheme()})
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

	buttons := container.New(layout.NewGridLayout(3),
		widget.NewButton("Button 1", func() {}),
		widget.NewButton("Button 2", func() {}),
		widget.NewButton("Button 3", func() {}),
	)

	content := container.NewBorder(nil, buttons, nil, nil, list)
	
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}