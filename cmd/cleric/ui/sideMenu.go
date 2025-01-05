package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pcarion/cleric/pkg/configuration"
)

type setMainContentFunc func(mainContent *MainContent)

type menuItem interface {
	label() string
	icon() fyne.Resource
	content() *MainContent
}

// definition of the side menu
type SideMenu struct {
	config     *configuration.Configuration
	mcpServers []*configuration.McpServerDescription
	list       *widget.List
}

func NewSideMenu() *SideMenu {
	config := configuration.LoadConfiguration()
	mcpServers := config.LoadMcpServers()

	return &SideMenu{
		config:     config,
		mcpServers: mcpServers,
	}
}

func (s *SideMenu) MakeNavigation(setMainContent setMainContentFunc, myApp fyne.App) fyne.CanvasObject {

	// use the mcp servers as the data
	data := make([]menuItem, 0, 1+len(s.mcpServers))
	// add the welcome item
	data = append(data, NewContentWelcome().menuItem())

	for _, server := range s.mcpServers {
		data = append(data, NewContentMcpServer(server).menuItem())
	}

	s.list = widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()), // 0
				widget.NewLabel("Template Object"),   // 1
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			// get the side menu item
			sideMenuItem := data[id]
			hbox := item.(*fyne.Container) // that's the Hbox

			hbox.Objects[0].(*widget.Icon).SetResource(sideMenuItem.icon())
			hbox.Objects[1].(*widget.Label).SetText(sideMenuItem.label())
		},
	)

	// add code if the item is clicked
	s.list.OnSelected = func(id widget.ListItemID) {
		sideMenuItem := data[id]
		setMainContent(sideMenuItem.content())
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			SetDarkTheme(myApp)
		}),
		widget.NewButton("Light", func() {
			SetLightTheme(myApp)
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, s.list)
}

func (s *SideMenu) SelectItem(id widget.ListItemID) {
	s.list.Select(id)
}
