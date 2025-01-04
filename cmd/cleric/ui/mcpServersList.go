package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"github.com/pcarion/cleric/pkg/configuration"
)

func NewMcpServersList(mcpServers []*configuration.McpServerDescription) *widget.List {
	// Create a binding data list
	data := binding.NewUntypedList()
	for _, server := range mcpServers {
		data.Append(server)
	}

	list := widget.NewListWithData(
		data,
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

	return list
}
