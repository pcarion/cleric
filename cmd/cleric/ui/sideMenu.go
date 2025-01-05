package ui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pcarion/cleric/pkg/configuration"
)

// definition of the side menu
type SideMenu struct {
	config     *configuration.Configuration
	mcpServers []*configuration.McpServerDescription
}

func NewSideMenu() *SideMenu {
	config := configuration.LoadConfiguration()
	mcpServers := config.LoadMcpServers()

	return &SideMenu{
		config:     config,
		mcpServers: mcpServers,
	}
}

func (s *SideMenu) MakeNavigation(setMainContent func(mainContent MainContent), myApp fyne.App) fyne.CanvasObject {

	data := make([]string, 20)
	for i := range data {
		data[i] = "Test Item " + strconv.Itoa(i)
	}

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id])
		},
	)

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			SetDarkTheme(myApp)
		}),
		widget.NewButton("Light", func() {
			SetLightTheme(myApp)
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, list)
}
