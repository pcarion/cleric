package ui

import (
	"fmt"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pcarion/cleric/internal/configuration"
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
		ContentId: c.mcpServer.Uuid,
		View: func(window fyne.Window) fyne.CanvasObject {
			statusLabel := widget.NewLabel("")

			// create a toolbar with buttons
			// 2 different toolbars for edit mode and visualization mode
			t := widget.NewToolbar()

			if c.IsEditMode() {
				t.Append(NewEditToolbar(c.editAction(), statusLabel))
			} else {
				t.Append(NewToolbarClaudeAction(c.claudeAction(), statusLabel))
				t.Append(widget.NewToolbarSpacer())
				t.Append(NewToolbarItemWithHover(theme.VisibilityIcon(), func() {
					inspectorArgs := []string{}
					inspectorArgs = append(inspectorArgs, "@modelcontextprotocol/inspector")
					inspectorArgs = c.mcpServer.Configuration.GetMcpInspectorArgs(inspectorArgs)
					fmt.Println(strings.Join(inspectorArgs, " "))
					go func() {
						cmd, err := runCommand("npx", inspectorArgs)
						if err != nil {
							fmt.Printf("Error running inspector: %v\n", err)
							return
						}
						cmd.Wait()
					}()
				}, mkHoverable("Start the MCP inspector", statusLabel)))
				t.Append(NewToolbarItemWithHover(theme.ContentCutIcon(), func() {
					c.DeleteMcpServer(c.mcpServer.Uuid)
				}, mkHoverable("Delete MCP server", statusLabel)))
				t.Append(NewEditToolbar(c.editAction(), statusLabel))
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
				nameWidgets.Add(c.newEditValueButtonWithValidator("Edit name", "name", c.mcpServer.Name, c.listActions.ValidateExistingMcpServerName(c.mcpServer.Uuid), func(value string) {
					c.mcpServer.Name = value
					c.listActions.RefreshSideMenu()
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

			if c.IsEditMode() {
				for index, value := range c.mcpServer.Configuration.Args {
					argLabel := widget.NewLabel(fmt.Sprintf("Argument %d", index))
					argLabel.TextStyle = fyne.TextStyle{Bold: true}
					argValue := widget.NewLabel(value)
					argValue.TextStyle = fyne.TextStyle{Monospace: true}
					argWidgets := container.NewHBox()
					argWidgets.Add(c.newDeleteValueButton(func() {
						c.mcpServer.Configuration.Args = append(c.mcpServer.Configuration.Args[:index], c.mcpServer.Configuration.Args[index+1:]...)
						c.listActions.RefreshCurrentContent()
						c.listActions.SaveMcpServers()
					}))
					argWidgets.Add(c.newEditValueButton("Edit argument", "argument", value, func(newValue string) {
						c.mcpServer.Configuration.Args[index] = newValue
						c.listActions.SaveMcpServers()
						c.listActions.RefreshCurrentContent()
					}))
					argControls := container.New(layout.NewBorderLayout(nil, nil, nil, argWidgets), argWidgets, argValue)
					formBuilder.AddField(argLabel, argControls)
				}
				addArgumentLabel := widget.NewLabel("")
				addArgumentButton := c.newAddValueButton("Add Argument", "argument", func(value string) {
					c.mcpServer.Configuration.Args = append(c.mcpServer.Configuration.Args, value)
					c.listActions.SaveMcpServers()
					c.listActions.RefreshCurrentContent()
				})
				addArgumentControls := container.NewHBox()
				addArgumentControls.Add(addArgumentButton)
				formBuilder.AddField(addArgumentLabel, addArgumentControls)
			} else {
				// Add rows for Arguments
				argumentsLabel := widget.NewLabel("Arguments")
				argumentsLabel.TextStyle = fyne.TextStyle{Bold: true}
				argumentsValue := widget.NewLabel(strings.Join(c.mcpServer.Configuration.Args, "\n"))
				argumentsValue.TextStyle = fyne.TextStyle{Monospace: true}
				argumentsWidgets := container.NewHBox()
				argumentsControls := container.New(layout.NewBorderLayout(nil, nil, nil, argumentsWidgets), argumentsWidgets, argumentsValue)
				formBuilder.AddField(argumentsLabel, argumentsControls)
			}

			// display environment variables
			envVars := orderEnvVars(c.mcpServer.Configuration.Env)
			if c.IsEditMode() {
				for _, value := range envVars {
					varName := value[0]
					varValue := value[1]
					argLabel := widget.NewLabel(fmt.Sprintf("variable: \"%s\"", varName))
					argLabel.TextStyle = fyne.TextStyle{Bold: true}
					argValue := widget.NewLabel(varValue)
					argValue.TextStyle = fyne.TextStyle{Monospace: true}
					argWidgets := container.NewHBox()
					argWidgets.Add(c.newDeleteValueButton(func() {
						delete(c.mcpServer.Configuration.Env, varName)
						c.listActions.RefreshCurrentContent()
						c.listActions.SaveMcpServers()
					}))
					argWidgets.Add(c.newEditValueButton(fmt.Sprintf("Edit environment variable: %s", varName), "value", varValue, func(newValue string) {
						c.mcpServer.Configuration.Env[varName] = newValue
						c.listActions.RefreshCurrentContent()
						c.listActions.SaveMcpServers()
					}))
					argControls := container.New(layout.NewBorderLayout(nil, nil, nil, argWidgets), argWidgets, argValue)
					formBuilder.AddField(argLabel, argControls)
				}
				addArgumentLabel := widget.NewLabel("")
				addArgumentButton := c.newAddValueButton("Add Environment variable", "name", func(value string) {
					c.mcpServer.Configuration.Env[value] = ""
					c.listActions.RefreshCurrentContent()
					c.listActions.SaveMcpServers()
				})
				addArgumentControls := container.NewHBox()
				addArgumentControls.Add(addArgumentButton)
				formBuilder.AddField(addArgumentLabel, addArgumentControls)
			} else {
				// prepare list of environment variables to be displayed
				lstEnvVars := []string{}
				for _, value := range envVars {
					lstEnvVars = append(lstEnvVars, fmt.Sprintf("%s=%s", value[0], value[1]))
				}

				// add row for environment variables
				envVarsLabel := widget.NewLabel("Environment Variables")
				envVarsLabel.TextStyle = fyne.TextStyle{Bold: true}
				envVarsValue := widget.NewLabel(strings.Join(lstEnvVars, "\n"))
				envVarsValue.TextStyle = fyne.TextStyle{Monospace: true}
				envVarsWidgets := container.NewHBox()
				if c.IsEditMode() {
					envVarsWidgets.Add(c.newEditValueButton("Edit environment variables", "environment variables", strings.Join(lstEnvVars, "\n"), func(value string) {
					}))
				}
				envVarsControls := container.New(layout.NewBorderLayout(nil, nil, nil, envVarsWidgets), envVarsWidgets, envVarsValue)
				formBuilder.AddField(envVarsLabel, envVarsControls)
			}

			// display the form
			pageContent := formBuilder.GetContainer()
			return container.NewBorder(t, statusLabel, nil, nil, container.NewVScroll(pageContent))
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

func (c *ContentMcpServer) newEditValueButtonWithValidator(title string, label string, value string, validator fyne.StringValidator, onSave func(string)) *widget.Button {
	button := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		c.displayEditValue(title, label, value, validator, onSave)
	})
	return button
}

func (c *ContentMcpServer) newEditValueButton(title string, label string, value string, onSave func(string)) *widget.Button {
	button := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		c.displayEditValue(title, label, value, nil, onSave)
	})
	return button
}

func (c *ContentMcpServer) newAddValueButton(title string, label string, onSave func(string)) *widget.Button {
	button := widget.NewButtonWithIcon(title, theme.ContentAddIcon(), func() {
		c.displayEditValue(title, label, "", nil, onSave)
	})
	return button
}

func (c *ContentMcpServer) newDeleteValueButton(onDelete func()) *widget.Button {
	button := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		onDelete()
	})
	return button
}

// shows a dialog to edit a value
func (c *ContentMcpServer) displayEditValue(title string, label string, value string, validator fyne.StringValidator, onSave func(string)) {
	nameEntry := widget.NewEntry()
	hintText := ""
	if validator != nil {
		nameEntry.Validator = validator
		hintText = "Must be unique and be only alphanumeric characters"
	}
	nameEntry.SetText(value)

	dialog := dialog.NewForm(
		title,
		"Save",
		"Cancel",
		[]*widget.FormItem{
			{Text: label, Widget: nameEntry, HintText: hintText},
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

func (c *ContentMcpServer) DeleteMcpServer(uuid string) {
	// ask for confirmation
	fmt.Println("delete MCP server", uuid)
	cnf := dialog.NewConfirm("Delete MCP Server", "Are you sure you want to delete this server?", func(confirm bool) {
		if confirm {
			// we remove the server from the list
			c.listActions.DeleteMcpServer(uuid)
			c.listActions.RefreshSideMenu()
			c.listActions.ResetListScroll()
		}
	}, c.window)
	cnf.SetDismissText("Cancel")
	cnf.SetConfirmText("Delete")
	cnf.Show()
}

func orderEnvVars(envVars map[string]string) [][2]string {
	lstEnvVars := [][2]string{}
	for key, value := range envVars {
		lstEnvVars = append(lstEnvVars, [2]string{key, value})
	}
	sort.Slice(lstEnvVars, func(i, j int) bool {
		return lstEnvVars[i][0] < lstEnvVars[j][0]
	})
	return lstEnvVars
}
