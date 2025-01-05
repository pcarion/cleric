package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pcarion/cleric/pkg/configuration"
)

type ContentMcpServer struct {
	mcpServer *configuration.McpServerDescription
}

func NewContentMcpServer(mcpServer *configuration.McpServerDescription) *ContentMcpServer {
	return &ContentMcpServer{mcpServer: mcpServer}
}

func (c *ContentMcpServer) menuItem() menuItem {
	return c
}

func (c *ContentMcpServer) content() *MainContent {
	return &MainContent{
		View: func(window fyne.Window) fyne.CanvasObject {
			// Create a vertical box to hold all rows
			vbox := container.NewVBox()

			// Add row for Name
			label := widget.NewLabel("Name")
			label.TextStyle = fyne.TextStyle{Bold: true}
			nameRow := container.NewGridWithColumns(3,
				label,
				widget.NewLabel(c.mcpServer.Name),
				widget.NewButton("Edit", func() {
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

			return container.NewVScroll(vbox)
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
