package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pcarion/cleric/pkg/configuration"
)

type ContentMcpServer struct {
	window          fyne.Window
	mcpServer       *configuration.McpServerDescription
	toolbar         *widget.Toolbar
	listRefreshable listRefreshable
	editMode        bool
}

func NewContentMcpServer(
	window fyne.Window,
	mcpServer *configuration.McpServerDescription,
	listRefreshable listRefreshable,
) *ContentMcpServer {
	return &ContentMcpServer{
		window:          window,
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
			vbox.Add(c.newLabelTitle("Name"))
			vbox.Add(c.newLabelValue("edit name", "name", c.mcpServer.Name, func(value string) {
				c.mcpServer.Name = value
				c.listRefreshable.RefreshCurrentContent()
			}))
			vbox.Add(widget.NewSeparator())

			// add row for description
			vbox.Add(c.newLabelTitle("Description"))
			vbox.Add(c.newLabelValue("edit description", "description", c.mcpServer.Description, func(value string) {
				c.mcpServer.Description = value
				c.listRefreshable.RefreshCurrentContent()
			}))
			vbox.Add(widget.NewSeparator())

			// Add row for Command
			vbox.Add(c.newLabelTitle("Command"))
			vbox.Add(c.newLabelValue("edit command", "command", c.mcpServer.Configuration.Command, func(value string) {
				c.mcpServer.Configuration.Command = value
				c.listRefreshable.RefreshCurrentContent()
			}))
			vbox.Add(widget.NewSeparator())

			// Add row for Arguments
			argumentsVbox := container.NewVBox()
			for _, arg := range c.mcpServer.Configuration.Args {
				argumentsVbox.Add(c.newListValue(arg))
			}
			vbox.Add(c.newLabelTitle("Arguments"))
			vbox.Add(argumentsVbox)

			// if edit mode, add a button to add an argument
			if c.IsEditMode() {
				vbox.Add(widget.NewButton("Add Argument", func() {
					c.mcpServer.Configuration.Args = append(c.mcpServer.Configuration.Args, "")
					c.listRefreshable.RefreshCurrentContent()
				}))
			}
			vbox.Add(widget.NewSeparator())

			// add row for environment variables
			envVarsVbox := container.NewVBox()
			for key, value := range c.mcpServer.Configuration.Env {
				envVarsVbox.Add(c.newEnvValue(key, value))
			}
			vbox.Add(c.newLabelTitle("Environment Variables"))
			vbox.Add(envVarsVbox)

			// if edit mode, add a button to add an environment variable
			if c.IsEditMode() {
				vbox.Add(widget.NewButton("Add Environment Variable", func() {
				}))
			}
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

func (c *ContentMcpServer) newLabelTitle(title string) *widget.Label {
	label := widget.NewLabel(title)
	label.TextStyle = fyne.TextStyle{Bold: true}
	return label
}

func (c *ContentMcpServer) newLabelValue(title string, label string, value string, onSave func(string)) fyne.CanvasObject {
	hbox := container.NewHBox()

	if c.IsEditMode() {
		t := widget.NewToolbar()
		t.Append(widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			c.displayEditLabelValue(title, label, value, onSave)
		}))
		hbox.Add(t)
	}
	lbl := widget.NewLabel(value)
	lbl.TextStyle = fyne.TextStyle{Monospace: true}
	hbox.Add(lbl)
	return hbox
}

func (c *ContentMcpServer) newListValue(value string) fyne.CanvasObject {
	hbox := container.NewHBox()
	if c.IsEditMode() {
		t := widget.NewToolbar()
		t.Append(widget.NewToolbarAction(theme.DocumentCreateIcon(), func() { fmt.Println("edit") }))
		t.Append(widget.NewToolbarAction(theme.ContentCutIcon(), func() { fmt.Println("cut") }))
		hbox.Add(t)
	}
	lbl := widget.NewLabel(value)
	lbl.TextStyle = fyne.TextStyle{Monospace: true}
	hbox.Add(lbl)
	return hbox
}

func (c *ContentMcpServer) newEnvValue(key string, value string) fyne.CanvasObject {
	hbox := container.NewHBox()
	if c.IsEditMode() {
		t := widget.NewToolbar()
		t.Append(widget.NewToolbarAction(theme.DocumentCreateIcon(), func() { fmt.Println("edit") }))
		t.Append(widget.NewToolbarAction(theme.ContentCutIcon(), func() { fmt.Println("cut") }))
		hbox.Add(t)
	}
	lblKey := widget.NewLabel(key)
	lblKey.TextStyle = fyne.TextStyle{Monospace: true}
	hbox.Add(lblKey)

	hbox.Add(widget.NewIcon(theme.NewThemedResource(theme.MoreVerticalIcon())))

	lblValue := widget.NewLabel(value)
	lblValue.TextStyle = fyne.TextStyle{Monospace: true}
	hbox.Add(lblValue)
	return hbox
}

func (c *ContentMcpServer) displayEditLabelValue(title string, label string, value string, onSave func(string)) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(value)

	dialog := dialog.NewForm(
		title,
		"Save",
		"Cancel",
		[]*widget.FormItem{
			{Text: label, Widget: nameEntry},
		},
		func(confirm bool) {
			if confirm {
				value = nameEntry.Text
				onSave(value)
			}
		},
		c.window,
	)
	dialog.Resize(fyne.NewSize(600, 300))
	dialog.Show()
}
