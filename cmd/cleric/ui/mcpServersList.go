package ui

import (
	"fmt"
	"image/color"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/pcarion/cleric/pkg/configuration"
)

type MCPServersList struct {
	data       binding.UntypedList
	list       *widget.List
	mcpServers []*configuration.McpServerDescription
}

func NewMcpServersList() *MCPServersList {
	config := configuration.LoadConfiguration()
	mcpServers := config.LoadMcpServers()

	// Create a binding data list
	data := binding.NewUntypedList()
	for _, server := range mcpServers {
		data.Append(server)
	}

	return &MCPServersList{
		data:       data,
		list:       nil,
		mcpServers: mcpServers,
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

			// we create a vbox with the name and the description
			vbox1 := container.NewVBox(
				widget.NewLabel("name"),
				widget.NewLabel("description"),
			)
			vbox2 := container.NewVBox(
				container.NewStack(background, widget.NewCheck("in Claude Desktop", nil)),
				widget.NewSeparator(),
				widget.NewButtonWithIcon(
					"Edit",
					theme.DocumentIcon(),
					func() {
						fmt.Println("@@ delete")
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
			label := leftColumn.Objects[0].(*widget.Label)
			label.SetText(mcpServer.Name)
			label = leftColumn.Objects[1].(*widget.Label)
			label.SetText(mcpServer.Description)

			stack := rightColumn.Objects[0].(*fyne.Container)
			background := stack.Objects[0].(*canvas.Rectangle)
			check := stack.Objects[1].(*widget.Check)
			check.SetChecked(mcpServer.InConfiguration)

			check.OnChanged = func(checked bool) {
				mcpServer.InConfiguration = checked
				l.list.Refresh() // Refresh to update styling
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

	return l.list
}

func (l *MCPServersList) SaveMcpServers() {
	config := configuration.LoadConfiguration()
	config.SaveMcpServers(l.mcpServers)
}

func (l *MCPServersList) RevertMcpServers() {
	config := configuration.LoadConfiguration()
	l.mcpServers = config.LoadMcpServers()
	l.data.Set([]interface{}{})
	for _, server := range l.mcpServers {
		l.data.Append(server)
	}
	l.list.Refresh()
}

func (l *MCPServersList) AddMcpServer(window fyne.Window) {
	nameEntry := widget.NewEntry()
	descEntry := widget.NewMultiLineEntry()
	cmdEntry := widget.NewEntry()
	argsEntry := widget.NewEntry()

	dialog := dialog.NewForm(
		"Add MCP Server",
		"Add",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Name", nameEntry),
			widget.NewFormItem("Description", descEntry),
			widget.NewFormItem("Command", cmdEntry),
			widget.NewFormItem("Arguments", argsEntry),
		},
		func(confirm bool) {
			if confirm {
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
			}
		},
		window,
	)
	dialog.Resize(fyne.NewSize(400, 300))
	dialog.Show()
}

func splitArgs(args string) []string {
	return strings.Split(args, " ")
}
