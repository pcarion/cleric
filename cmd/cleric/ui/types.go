package ui

import "fyne.io/fyne/v2"

type MainContent struct {
	Title string
	View  func(w fyne.Window) fyne.CanvasObject
}

type listRefreshable interface {
	RefreshCurrentContent()
	RefreshSideMenu()
}
