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
	fmt.Printf("@@ content() - editMode:%t\n", c.editMode)
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
			label := widget.NewLabel("Name")
			label.TextStyle = fyne.TextStyle{Bold: true}
			nameRow := container.NewGridWithColumns(3,
				label,
				widget.NewLabel(c.mcpServer.Name),
				widget.NewButton(fmt.Sprintf("Edit (%t)", c.editMode), func() {
					// TODO: Implement edit functionality
				}),
			)
			vbox.Add(nameRow)
			// add row for description
			descriptionRow := container.NewGridWithColumns(3,
				widget.NewLabel("Description"),
				widget.NewLabel(c.mcpServer.Description),
				widget.NewButton("Edit", func() {
					// TODO: Implement edit functionality
				}),
			)
			vbox.Add(descriptionRow)

			// Add row for Command
			commandRow := container.NewGridWithColumns(3,
				widget.NewLabel("Command"),
				widget.NewLabel(c.mcpServer.Configuration.Command),
				widget.NewButton("Edit", func() {
					// TODO: Implement edit functionality
				}),
			)
			vbox.Add(commandRow)

			// Add row for Arguments
			argumentsVbox := container.NewVBox()
			for _, arg := range c.mcpServer.Configuration.Args {
				argumentsVbox.Add(widget.NewLabel(arg))
			}
			argsRow := container.NewGridWithColumns(3,
				widget.NewLabel("Arguments"),
				argumentsVbox,
				widget.NewButton("Edit", func() {
					// TODO: Implement edit functionality
				}),
			)
			vbox.Add(argsRow)

			// add row for environment variables
			envVarsVbox := container.NewVBox()
			for key, value := range c.mcpServer.Configuration.Env {
				envVarsVbox.Add(widget.NewLabel(key + "=" + value))
			}
			envVarsRow := container.NewGridWithColumns(3,
				widget.NewLabel("Environment Variables"),
				envVarsVbox,
				widget.NewButton("Edit", func() {
					// TODO: Implement edit functionality
				}),
			)
			vbox.Add(envVarsRow)

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
