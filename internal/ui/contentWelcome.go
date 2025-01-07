package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ContentWelcome struct {
}

func NewContentWelcome() *ContentWelcome {
	return &ContentWelcome{}
}

func (c *ContentWelcome) menuItem() menuItem {
	return c
}

func (c *ContentWelcome) content() *MainContent {
	return &MainContent{
		ContentId: "welcome",
		View: func(window fyne.Window) fyne.CanvasObject {
			return widget.NewLabel("Welcome")
		},
	}
}

func (c *ContentWelcome) label() string {
	return "Welcome"
}

func (c *ContentWelcome) icon() fyne.Resource {
	return theme.HomeIcon()
}
