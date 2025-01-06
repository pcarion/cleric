package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pcarion/cleric/pkg/configuration"
)

type ContentMcpServer struct {
	mcpServer       *configuration.McpServerDescription
	toolbar         *widget.Toolbar
	listRefreshable listRefreshable
	editMode        bool
}

func NewContentMcpServer(
	mcpServer *configuration.McpServerDescription,
	listRefreshable listRefreshable,
) *ContentMcpServer {
	return &ContentMcpServer{
		mcpServer:       mcpServer,
		listRefreshable: listRefreshable,
		editMode:        false,
	}
}

func (c *ContentMcpServer) menuItem() menuItem {
	return c
}

func (c *ContentMcpServer) claudeAction() ClaudeAction {
	return c
}

func (c *ContentMcpServer) IsServerInClaude() bool {
	return c.mcpServer.InConfiguration
}

func (c *ContentMcpServer) AddToClaude() {
	c.mcpServer.InConfiguration = true
	c.toolbar.Refresh()
	c.listRefreshable.RefreshSideMenu()
}

func (c *ContentMcpServer) RemoveFromClaude() {
	c.mcpServer.InConfiguration = false
	c.toolbar.Refresh()
	c.listRefreshable.RefreshSideMenu()
}

func (c *ContentMcpServer) editAction() ToolbarEditAction {
	return c
}

func (c *ContentMcpServer) IsEditMode() bool {
	return c.editMode
}

func (c *ContentMcpServer) CancelEditMode() {
	c.editMode = false
	c.listRefreshable.RefreshCurrentContent()
}

func (c *ContentMcpServer) EditMode() {
	c.editMode = true
	c.listRefreshable.RefreshCurrentContent()
}

func (c *ContentMcpServer) content() *MainContent {
	return &MainContent{
		View: func(window fyne.Window) fyne.CanvasObject {
			// create a toolbar with buttons
			t := widget.NewToolbar()

			if c.IsEditMode() {
				t.Append(NewEditToolbar(c.editAction()))
			} else {
				t.Append(NewToolbarClaudeAction(c.claudeAction()))
				t.Append(widget.NewToolbarSeparator())
				t.Append(widget.NewToolbarSpacer())
				t.Append(widget.NewToolbarAction(theme.ContentCutIcon(), func() { fmt.Println("Cut") }))
				t.Append(widget.NewToolbarAction(theme.ContentCopyIcon(), func() { fmt.Println("Copy") }))
				t.Append(NewEditToolbar(c.editAction()))
			}

			t.Refresh()
			c.toolbar = t

			// Create a vertical box to hold all rows
			vbox := container.NewVBox()

			// Add row for Name
			vbox.Add(newLabelTitle("Name"))
			vbox.Add(newLabelValue(c.mcpServer.Name))
			vbox.Add(widget.NewSeparator())

			// add row for description
			vbox.Add(newLabelTitle("Description"))
			vbox.Add(newLabelValue(c.mcpServer.Description))
			vbox.Add(widget.NewSeparator())

			// Add row for Command
			vbox.Add(newLabelTitle("Command"))
			vbox.Add(newLabelValue(c.mcpServer.Configuration.Command))
			vbox.Add(widget.NewSeparator())

			// Add row for Arguments
			argumentsVbox := container.NewVBox()
			for _, arg := range c.mcpServer.Configuration.Args {
				argumentsVbox.Add(widget.NewLabel(arg))
			}
			vbox.Add(newLabelTitle("Arguments"))
			vbox.Add(argumentsVbox)
			vbox.Add(widget.NewSeparator())

			// add row for environment variables
			envVarsVbox := container.NewVBox()
			for key, value := range c.mcpServer.Configuration.Env {
				envVarsVbox.Add(widget.NewLabel(key + "=" + value))
			}
			vbox.Add(newLabelTitle("Environment Variables"))
			vbox.Add(envVarsVbox)
			vbox.Add(widget.NewSeparator())

			return container.NewBorder(t, nil, nil, nil, container.NewVScroll(vbox))
		},
	}
}

func (c *ContentMcpServer) label() string {
	return c.mcpServer.Name
}

func (c *ContentMcpServer) icon() fyne.Resource {
	if c.mcpServer.InConfiguration {
		return theme.CheckButtonCheckedIcon()
	}
	return theme.CheckButtonIcon()
}

func newLabelTitle(title string) *widget.Label {
	label := widget.NewLabel(title)
	label.TextStyle = fyne.TextStyle{Bold: true}
	return label
}

func newLabelValue(value string) fyne.CanvasObject {
	label := widget.NewLabel(value)
	label.TextStyle = fyne.TextStyle{Monospace: true}
	return label
}
