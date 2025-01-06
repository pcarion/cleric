package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/pcarion/cleric/cmd/cleric/ui"
)

func main() {
	myApp := app.New()
	ui.SetDarkTheme(myApp)
	myWindow := myApp.NewWindow("Claude Mcp Servers")
	var currentMainContent *ui.MainContent = nil

	// the right side of the side menu
	mainContentStack := container.NewStack()
	mainContentContainer := container.NewBorder(nil, nil, nil, nil, mainContentStack)

	setMainContent := func(mainContent *ui.MainContent) {
		currentMainContent = mainContent
		mainContentStack.Objects = []fyne.CanvasObject{mainContent.View(myWindow)}
		mainContentStack.Refresh()
	}

	refreshMainContent := func() {
		if currentMainContent != nil {
			mainContentStack.Objects = []fyne.CanvasObject{currentMainContent.View(myWindow)}
			mainContentStack.Refresh()
		}
	}

	sideMenu := ui.NewSideMenu(refreshMainContent)

	split := container.NewHSplit(sideMenu.MakeNavigation(setMainContent, refreshMainContent, myApp), mainContentContainer)
	split.Offset = 0.2
	myWindow.SetContent(split)

	sideMenu.SelectItem(0)

	myWindow.Resize(fyne.NewSize(640, 460))
	// Center the window on screen
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}
