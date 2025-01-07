package ui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ContentAddMcpServer struct {
	actions ServerListActions
}

func NewContentAddMcpServer(actions ServerListActions) *ContentAddMcpServer {
	return &ContentAddMcpServer{actions: actions}
}

func (c *ContentAddMcpServer) menuItem() menuItem {
	return c
}

func (c *ContentAddMcpServer) content() *MainContent {
	return &MainContent{
		View: func(window fyne.Window) fyne.CanvasObject {
			// Create input field
			createLabel := widget.NewLabel("Add a new MCP Server")
			createLabel.Alignment = fyne.TextAlignLeading
			createLabel.TextStyle = fyne.TextStyle{Bold: true}

			nameEntry := widget.NewEntry()
			nameEntry.SetPlaceHolder("Server Name")
			nameEntry.Validator = func(s string) error {
				if !isValidServerName(s) {
					return errors.New("name can only contain letters, numbers, and underscores")
				}
				if c.actions.ValidateNewMcpServerName(s) != nil {
					return errors.New("a server with this name already exists")
				}
				return nil
			}

			// add an error label
			hintLabel1 := widget.NewLabel("- must not already exist")
			hintLabel1.TextStyle = fyne.TextStyle{Italic: true}
			hintLabel2 := widget.NewLabel("- be an alphanumeric string")
			hintLabel2.TextStyle = fyne.TextStyle{Italic: true}

			// Create button (moved before OnChanged to reference it)
			createButton := widget.NewButton("Add server", func() {
				err := c.actions.AddMcpServer(nameEntry.Text)
				if err != nil {
					dialog.ShowError(err, window)
				}
				c.actions.RefreshSideMenu()
				c.actions.RefreshCurrentContent()
			})
			createButton.Disable() // Initially disabled

			// Add validation check on text change
			nameEntry.OnChanged = func(s string) {
				if nameEntry.Validate() == nil {
					createButton.Enable()
				} else {
					createButton.Disable()
				}
			}

			// Arrange widgets vertically and center them
			content := container.NewVBox(
				createLabel,
				layout.NewSpacer(),
				nameEntry,
				hintLabel1,
				hintLabel2,
				layout.NewSpacer(),
				createButton,
			)

			// Center the content in the window
			centered := container.NewCenter(content)

			return centered
		},
	}
}

func (c *ContentAddMcpServer) label() string {
	return "Add MCP Server"
}

func (c *ContentAddMcpServer) icon() fyne.Resource {
	return theme.ContentAddIcon()
}
