package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pcarion/cleric/pkg/configuration"
)

type ContentMcpServer struct {
	window      fyne.Window
	mcpServer   *configuration.McpServerDescription
	toolbar     *widget.Toolbar
	listActions ServerListActions
	editMode    bool
}

func NewContentMcpServer(
	window fyne.Window,
	mcpServer *configuration.McpServerDescription,
	listActions ServerListActions,
) *ContentMcpServer {
	return &ContentMcpServer{
		window:      window,
		mcpServer:   mcpServer,
		listActions: listActions,
		editMode:    false,
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
	c.listActions.RefreshSideMenu()
	c.listActions.SaveMcpServers()
}

func (c *ContentMcpServer) RemoveFromClaude() {
	c.mcpServer.InConfiguration = false
	c.toolbar.Refresh()
	c.listActions.RefreshSideMenu()
	c.listActions.SaveMcpServers()
}

func (c *ContentMcpServer) editAction() ToolbarEditAction {
	return c
}

func (c *ContentMcpServer) IsEditMode() bool {
	return c.editMode
}

func (c *ContentMcpServer) CancelEditMode() {
	c.editMode = false
	c.listActions.RefreshCurrentContent()
}

func (c *ContentMcpServer) EditMode() {
	c.editMode = true
	c.listActions.RefreshCurrentContent()
}

func (c *ContentMcpServer) content() *MainContent {
	return &MainContent{
		View: func(window fyne.Window) fyne.CanvasObject {
			// create a toolbar with buttons
			// 2 different toolbar for edit mode and normal mode
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

			// Create form elements
			formBuilder := NewFormBuilder()

			nameLabel := widget.NewLabel("Name")
			nameLabel.TextStyle = fyne.TextStyle{Bold: true}
			nameValue := widget.NewLabel(c.mcpServer.Name)
			nameWidgets := container.NewHBox()
			if c.IsEditMode() {
				nameWidgets.Add(c.newEditValueButton("Edit name", "name", c.mcpServer.Name, func(value string) {
					c.mcpServer.Name = value
					c.listActions.RefreshCurrentContent()
					c.listActions.SaveMcpServers()
				}))
			}
			nameControls := container.New(layout.NewBorderLayout(nil, nil, nil, nameWidgets), nameWidgets, nameValue)
			formBuilder.AddField(nameLabel, nameControls)

			// add row for description
			descriptionLabel := widget.NewLabel("Description")
			descriptionLabel.TextStyle = fyne.TextStyle{Bold: true}
			descriptionValue := widget.NewLabel(c.mcpServer.Description)
			descriptionValue.TextStyle = fyne.TextStyle{Italic: true}
			descriptionWidgets := container.NewHBox()
			if c.IsEditMode() {
				descriptionWidgets.Add(c.newEditValueButton("Edit description", "description", c.mcpServer.Description, func(value string) {
					c.mcpServer.Description = value
					c.listActions.RefreshCurrentContent()
					c.listActions.SaveMcpServers()
				}))
			}
			descriptionControls := container.New(layout.NewBorderLayout(nil, nil, nil, descriptionWidgets), descriptionWidgets, descriptionValue)
			formBuilder.AddField(descriptionLabel, descriptionControls)

			// Add row for Command
			commandLabel := widget.NewLabel("Command")
			commandLabel.TextStyle = fyne.TextStyle{Bold: true}
			commandValue := widget.NewLabel(c.mcpServer.Configuration.Command)
			commandValue.TextStyle = fyne.TextStyle{Monospace: true}
			commandWidgets := container.NewHBox()
			if c.IsEditMode() {
				commandWidgets.Add(c.newEditValueButton("Edit command", "command", c.mcpServer.Configuration.Command, func(value string) {
					c.mcpServer.Configuration.Command = value
					c.listActions.RefreshCurrentContent()
					c.listActions.SaveMcpServers()
				}))
			}
			commandControls := container.New(layout.NewBorderLayout(nil, nil, nil, commandWidgets), commandWidgets, commandValue)
			formBuilder.AddField(commandLabel, commandControls)

			// Add rows for Arguments
			argumentsLabel := widget.NewLabel("Arguments")
			argumentsLabel.TextStyle = fyne.TextStyle{Bold: true}
			argumentsValue := widget.NewLabel(strings.Join(c.mcpServer.Configuration.Args, "\n"))
			argumentsValue.TextStyle = fyne.TextStyle{Monospace: true}
			argumentsWidgets := container.NewHBox()
			if c.IsEditMode() {
				argumentsWidgets.Add(c.newEditValueButton("Edit arguments", "arguments", strings.Join(c.mcpServer.Configuration.Args, " "), func(value string) {
					c.mcpServer.Configuration.Args = strings.Split(value, " ")
					c.listActions.RefreshCurrentContent()
					c.listActions.SaveMcpServers()
				}))
			}
			argumentsControls := container.New(layout.NewBorderLayout(nil, nil, nil, argumentsWidgets), argumentsWidgets, argumentsValue)
			formBuilder.AddField(argumentsLabel, argumentsControls)

			// for index, arg := range c.mcpServer.Configuration.Args {
			// 	argumentsVbox.Add(c.newListValue(arg, func(value string) {
			// 		c.mcpServer.Configuration.Args[index] = value
			// 		c.listActions.RefreshCurrentContent()
			// 		c.listActions.SaveMcpServers()
			// 	}))
			// }

			// if edit mode, add a button to add an argument
			// if c.IsEditMode() {
			// 	vbox.Add(widget.NewButton("Add Argument", func() {
			// 		c.displayEditValue("Add Argument", "new argument", "", func(value string) {
			// 			c.mcpServer.Configuration.Args = append(c.mcpServer.Configuration.Args, value)
			// 			c.listActions.RefreshCurrentContent()
			// 			c.listActions.SaveMcpServers()
			// 		})
			// 	}))
			// }
			// vbox.Add(widget.NewSeparator())

			lstEnvVars := []string{}
			for key, value := range c.mcpServer.Configuration.Env {
				lstEnvVars = append(lstEnvVars, fmt.Sprintf("%s=%s", key, value))
			}

			// add row for environment variables
			envVarsLabel := widget.NewLabel("Environment Variables")
			envVarsLabel.TextStyle = fyne.TextStyle{Bold: true}
			envVarsValue := widget.NewLabel(strings.Join(lstEnvVars, "\n"))
			envVarsValue.TextStyle = fyne.TextStyle{Monospace: true}
			envVarsWidgets := container.NewHBox()
			if c.IsEditMode() {
				envVarsWidgets.Add(c.newEditValueButton("Edit environment variables", "environment variables", strings.Join(lstEnvVars, "\n"), func(value string) {
					//					c.mcpServer.Configuration.Env = strings.Split(value, "\n")
					//					c.listActions.RefreshCurrentContent()
					//					c.listActions.SaveMcpServers()
				}))
			}
			envVarsControls := container.New(layout.NewBorderLayout(nil, nil, nil, envVarsWidgets), envVarsWidgets, envVarsValue)
			formBuilder.AddField(envVarsLabel, envVarsControls)
			// for key, value := range c.mcpServer.Configuration.Env {
			// 	envVarsWidgets.Add(c.newEnvValue(key, value))
			// }

			// // if edit mode, add a button to add an environment variable
			// if c.IsEditMode() {
			// 	vbox.Add(widget.NewButton("Add Environment Variable", func() {
			// 	}))
			// }
			// vbox.Add(widget.NewSeparator())

			// v2
			pageContent := formBuilder.GetContainer()
			return container.NewBorder(t, nil, nil, nil, container.NewVScroll(pageContent))
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

func (c *ContentMcpServer) newEditValueButton(title string, label string, value string, onSave func(string)) *widget.Button {
	button := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		c.displayEditValue(title, label, value, onSave)
	})
	return button
}

func (c *ContentMcpServer) newLabelValue(title string, label string, value string, onSave func(string)) fyne.CanvasObject {
	hbox := container.NewHBox()

	if c.IsEditMode() {
		t := widget.NewToolbar()
		t.Append(widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			c.displayEditValue(title, label, value, onSave)
		}))
		hbox.Add(t)
	}
	lbl := widget.NewLabel(value)
	lbl.TextStyle = fyne.TextStyle{Monospace: true}
	hbox.Add(lbl)
	return hbox
}

func (c *ContentMcpServer) newListValue(value string, onSave func(string)) fyne.CanvasObject {
	hbox := container.NewHBox()
	if c.IsEditMode() {
		t := widget.NewToolbar()
		t.Append(widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			c.displayEditValue("Edit Argument", "argument", value, onSave)
		}))
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

// shows a dialog to edit a value
func (c *ContentMcpServer) displayEditValue(title string, label string, value string, onSave func(string)) {
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
