package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/pcarion/cleric/pkg/claude"
)

func NewMcpServersList(mcpServers []*claude.McpServerDescription) *widget.List {
	list := widget.NewList(
		func() int {
			return len(mcpServers)
		},
		func() fyne.CanvasObject {
			check := widget.NewCheck("", func(bool) {})
			return container.NewBorder(
				nil, nil,
				widget.NewLabel("template"),
				check,
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			cont := o.(*fyne.Container)
			label := cont.Objects[0].(*widget.Label)

			label.SetText(mcpServers[i].Name)
		})

	return list
}
