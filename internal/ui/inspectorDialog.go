package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/pcarion/cleric/internal/configuration"
)

func ShowInspectorDialog(window fyne.Window, mcpServer *configuration.McpServerDescription) {

	inspectorArgs := mcpServer.Configuration.GetMcpInspectorArgs(configuration.McpVersion030)
	fmt.Println(strings.Join(inspectorArgs, " "))

	labelMcp030 := widget.NewLabel("")
	labelMcp030.Alignment = fyne.TextAlignLeading

	if mcpServer.Configuration.HasEnvironmentVariables() {
		labelMcp030.SetText("copy and paste the following command in your terminal to start the MCP Inspector\n\nIMPORTANT: you need to set the environment variables in the MCP inspector web page")
	} else {
		labelMcp030.SetText("copy and paste the following command in your terminal to start the MCP Inspector")
	}
	// create a label with copy to hold the inspector args
	inspectorArgsLabel := NewTextWithCopy(strings.Join(inspectorArgs, " "), window)

	// Create dialog with output and kill button
	content := container.NewBorder(labelMcp030, nil, nil, nil,
		container.NewVScroll(inspectorArgsLabel))
	d := dialog.NewCustom("MCP Inspector", "Close", content, window)
	d.Resize(fyne.NewSize(600, 400))
	d.Show()
}
