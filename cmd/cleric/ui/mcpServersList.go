package ui

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/pcarion/cleric/pkg/configuration"
)

type MCPServersList struct {
	config     *configuration.Configuration
	data       binding.UntypedList
	list       *widget.List
	mcpServers []*configuration.McpServerDescription
	window     fyne.Window
}

func NewMcpServersList(window fyne.Window) *MCPServersList {
	config := configuration.LoadConfiguration()
	mcpServers := config.LoadMcpServers()

	// Create a binding data list
	data := binding.NewUntypedList()
	for _, server := range mcpServers {
		data.Append(server)
	}

	return &MCPServersList{
		config:     config,
		data:       data,
		list:       nil,
		mcpServers: mcpServers,
		window:     window,
	}
}

func (l *MCPServersList) GetList() *widget.List {
	if l.list != nil {
		return l.list
	}
	l.list = widget.NewListWithData(
		l.data,
		func() fyne.CanvasObject {
			// Create a container with background
			background := canvas.NewRectangle(color.RGBA{R: 205, G: 92, B: 92, A: 180})

			// create a horizontal box with the name and description
			name := widget.NewLabel("")
			name.TextStyle.Bold = true

			description := widget.NewLabel("")
			description.TextStyle.Italic = true

			hbox := container.NewHBox(
				name,
				layout.NewSpacer(),
				description,
			)

			command := widget.NewLabel("")
			command.TextStyle.Monospace = true

			// we create a vbox with the name and the description
			vbox1 := container.NewVBox(
				hbox,
				command,
			)
			vbox2 := container.NewVBox(
				container.NewStack(background, widget.NewCheck("in Claude Desktop", nil)),
				widget.NewSeparator(),
				widget.NewButtonWithIcon(
					"Edit",
					theme.DocumentIcon(),
					func() {
					},
				),
				widget.NewButtonWithIcon(
					"Delete",
					theme.DeleteIcon(),
					func() {
					},
				),
			)

			content := container.NewBorder(
				nil, nil,
				vbox1,
				vbox2,
			)

			return content
		},
		func(di binding.DataItem, o fyne.CanvasObject) {
			server, err := di.(binding.Untyped).Get()
			if err != nil {
				log.Fatalf("Failed to get server: %v", err)
			}
			mcpServer := server.(*configuration.McpServerDescription)

			// Get the background container
			cont := o.(*fyne.Container) // that's the Border component
			leftColumn := cont.Objects[0].(*fyne.Container)
			rightColumn := cont.Objects[1].(*fyne.Container)

			// Update the rest of the content
			hbox := leftColumn.Objects[0].(*fyne.Container)
			name := hbox.Objects[0].(*widget.Label)
			name.SetText(mcpServer.Name)
			description := hbox.Objects[2].(*widget.Label)
			description.SetText(mcpServer.Description)
			command := leftColumn.Objects[1].(*widget.Label)
			command.SetText(fmt.Sprintf("%s %s", mcpServer.Configuration.Command, strings.Join(mcpServer.Configuration.Args, " ")))

			// right column
			stack := rightColumn.Objects[0].(*fyne.Container)
			background := stack.Objects[0].(*canvas.Rectangle)
			check := stack.Objects[1].(*widget.Check)
			check.SetChecked(mcpServer.InConfiguration)
			editButton := rightColumn.Objects[2].(*widget.Button)
			deleteButton := rightColumn.Objects[3].(*widget.Button)

			check.OnChanged = func(checked bool) {
				mcpServer.InConfiguration = checked
				l.list.Refresh() // Refresh to update styling
			}
			editButton.OnTapped = func() {
				// we need the index of the server
				l.EditMcpServer(mcpServer)
			}
			deleteButton.OnTapped = func() {
				fmt.Println("delete button tapped")
				l.DeleteMcpServer(mcpServer)
			}
			// Set background color based on some condition
			if mcpServer.InConfiguration {
				// #DA7857
				background.FillColor = color.RGBA{R: 218, G: 120, B: 87, A: 200}
			} else {
				// transparent
				background.FillColor = color.RGBA{R: 87, G: 166, B: 218, A: 0}
			}

		})

	l.list.HideSeparators = false
	return l.list
}

func (l *MCPServersList) saveMcpServers() {
	l.config.SaveMcpServers(l.mcpServers)
}

func (l *MCPServersList) DeleteMcpServer(server *configuration.McpServerDescription) {
	// ask for confirmation
	fmt.Println("delete MCP server", server.Name, server.Uuid)
	cnf := dialog.NewConfirm("Delete MCP Server", "Are you sure you want to delete this server?", func(confirm bool) {
		if confirm {
			// we remove the server from the list
			for i, s := range l.mcpServers {
				if s.Uuid == server.Uuid {
					l.mcpServers = append(l.mcpServers[:i], l.mcpServers[i+1:]...)
					break
				}
			}
			l.data.Remove(server)
			l.list.Refresh()
			l.saveMcpServers()
		}
	}, l.window)
	cnf.SetDismissText("Cancel")
	cnf.SetConfirmText("Delete")
	cnf.Show()
}

func (l *MCPServersList) AddMcpServer() {
	nameEntry := widget.NewEntry()
	nameEntry.Validator = l.ValidateNewName

	descEntry := widget.NewEntry()
	cmdEntry := widget.NewEntry()
	argsEntry := widget.NewEntry()

	dialog := dialog.NewForm(
		"Add MCP Server",
		"Add",
		"Cancel",
		[]*widget.FormItem{
			{Text: "Name", Widget: nameEntry, HintText: "Must be unique"},
			{Text: "Description", Widget: descEntry},
			{Text: "Command", Widget: cmdEntry},
			{Text: "Arguments", Widget: argsEntry},
		},
		func(confirm bool) {
			if confirm {
				// we need to check that there is no other server with the same name
				for _, otherServer := range l.mcpServers {
					if otherServer.Name == nameEntry.Text {
						dialog.ShowError(errors.New("a server with this name already exists"), l.window)
						return
					}
				}
				newServer := &configuration.McpServerDescription{
					Name:        nameEntry.Text,
					Description: descEntry.Text,
					Configuration: configuration.McpServerConfiguration{
						Command: cmdEntry.Text,
						Args:    splitArgs(argsEntry.Text),
						Env:     map[string]string{},
					},
				}
				l.mcpServers = append(l.mcpServers, newServer)
				l.data.Append(newServer)
				l.saveMcpServers()
			}
		},
		l.window,
	)
	dialog.Resize(fyne.NewSize(600, 300))
	dialog.Show()
}

func (l *MCPServersList) EditMcpServer(server *configuration.McpServerDescription) {
	nameEntry := widget.NewEntry()
	nameEntry.Validator = l.ValidateExistingName(server.Uuid)
	descEntry := widget.NewEntry()
	cmdEntry := widget.NewEntry()
	argsEntry := widget.NewEntry()

	validation.NewRegexp("^[a-zA-Z0-9_-]+$", "Invalid name")

	nameEntry.SetText(server.Name)
	descEntry.SetText(server.Description)
	cmdEntry.SetText(server.Configuration.Command)
	argsEntry.SetText(strings.Join(server.Configuration.Args, " "))

	dialog := dialog.NewForm(
		"Edit MCP Server",
		"Save",
		"Cancel",
		[]*widget.FormItem{
			{Text: "Name", Widget: nameEntry, HintText: "Must be unique"},
			{Text: "Description", Widget: descEntry},
			{Text: "Command", Widget: cmdEntry},
			{Text: "Arguments", Widget: argsEntry},
		},
		func(confirm bool) {
			if confirm {
				server.Name = nameEntry.Text
				var count = 0
				// we need to check that there is no other server with the same name
				for _, otherServer := range l.mcpServers {
					if otherServer.Name == server.Name {
						count++
					}
				}
				if count > 1 {
					dialog.ShowError(errors.New("a server with this name already exists"), l.window)
					return
				}
				server.Description = descEntry.Text
				server.Configuration.Command = cmdEntry.Text
				server.Configuration.Args = splitArgs(argsEntry.Text)
				l.list.Refresh()
				l.saveMcpServers()
			}
		},
		l.window,
	)
	dialog.Resize(fyne.NewSize(600, 300))
	dialog.Show()
}

func splitArgs(args string) []string {
	return strings.Split(args, " ")
}

func (l *MCPServersList) ValidateNewName(name string) error {
	for _, server := range l.mcpServers {
		if server.Name == name {
			return errors.New("a server with this name already exists")
		}
	}
	return nil
}

func (l *MCPServersList) ValidateExistingName(uuid string) func(name string) error {
	return func(name string) error {
		for _, server := range l.mcpServers {
			if server.Uuid == uuid {
				continue
			}
			if server.Name == name {
				return errors.New("a server with this name already exists")
			}
		}
		return nil
	}
}
