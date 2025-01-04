package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
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
			check := widget.NewCheck("", func(bool) {})
			return container.NewBorder(
				nil, nil,
				widget.NewLabel("template"),
				check,
			)
		},
		func(di binding.DataItem, o fyne.CanvasObject) {
			server, err := di.(binding.Untyped).Get()
			if err != nil {
				log.Fatalf("Failed to get server: %v", err)
			}
			mcpServer := server.(*configuration.McpServerDescription)
			cont := o.(*fyne.Container)
			label := cont.Objects[0].(*widget.Label)
			label.SetText(mcpServer.Name)
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
				// newServer := &configuration.McpServerDescription{
				// 	Name:        nameEntry.Text,
				// 	Description: descEntry.Text,
				// 	Command:     cmdEntry.Text,
				// 	Arguments:   argsEntry.Text,
				// }
				// l.mcpServers = append(l.mcpServers, newServer)
				// l.data.Append(newServer)
				// l.SaveMcpServers()
			}
		},
		window,
	)
	dialog.Resize(fyne.NewSize(400, 300))
	dialog.Show()
}
