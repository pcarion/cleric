package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

var data = []string{"a", "string", "list"}

func main() {
	myApp := app.New()
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
			return container.NewBorder(
				nil, nil, 
				widget.NewLabel("template"),
				widget.NewCheck("", func(bool) {}),
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