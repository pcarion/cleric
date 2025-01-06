package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ContentAddMcpServer struct {
}

func NewContentAddMcpServer() *ContentAddMcpServer {
	return &ContentAddMcpServer{}
}

func (c *ContentAddMcpServer) menuItem() menuItem {
	return c
}

func (c *ContentAddMcpServer) content() *MainContent {
	return &MainContent{
		View: func(window fyne.Window) fyne.CanvasObject {
			return widget.NewLabel("Add MCP Server")
		},
	}
}

func (c *ContentAddMcpServer) label() string {
	return "Add MCP Server"
}

func (c *ContentAddMcpServer) icon() fyne.Resource {
	return theme.ContentAddIcon()
}
