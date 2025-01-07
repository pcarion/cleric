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
	ValidateNewMcpServerName(name string) error
	// can be used as validators for the name of an existing mcp server
	ValidateExistingMcpServerName(uuid string) func(name string) error
	// delete an existing mcp server
	DeleteMcpServer(uuid string)
	// add a new mcp server
	AddMcpServer(name string) error
}
