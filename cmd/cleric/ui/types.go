package ui

import (
	"fyne.io/fyne/v2"
)

type MainContent struct {
	Title string
	View  func(w fyne.Window) fyne.CanvasObject
}

type ServerListActions interface {
	RefreshCurrentContent()
	RefreshSideMenu()
	SaveMcpServers()
	// can be used as validators for the name of a new mcp server
	ValidateNewName(name string) error
	// can be used as validators for the name of an existing mcp server
	ValidateExistingName(uuid string) func(name string) error
	DeleteMcpServer(uuid string)
}
