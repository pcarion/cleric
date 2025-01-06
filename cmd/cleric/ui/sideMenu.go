package ui

import (
	"errors"

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
	config             *configuration.Configuration
	mcpServers         []*configuration.McpServerDescription
	list               *widget.List
	refreshMainContent func()
}

func NewSideMenu(refreshMainContent func()) *SideMenu {
	config := configuration.LoadConfiguration()
	mcpServers := config.LoadMcpServers()

	return &SideMenu{
		config:             config,
		mcpServers:         mcpServers,
		refreshMainContent: refreshMainContent,
	}
}

func (s *SideMenu) SaveMcpServers() {
	s.config.SaveMcpServers(s.mcpServers)
}

func (s *SideMenu) AsServerListActions() ServerListActions {
	return s
}

func (s *SideMenu) RefreshSideMenu() {
	if s.list != nil {
		s.list.Refresh()
	}
}

func (s *SideMenu) RefreshCurrentContent() {
	if s.refreshMainContent != nil {
		s.refreshMainContent()
	}
}

func (s *SideMenu) MakeNavigation(
	setMainContent setMainContentFunc,
	refreshMainContent func(),
	myApp fyne.App,
	window fyne.Window,
) fyne.CanvasObject {

	// use the mcp servers as the data
	data := make([]menuItem, 0, 1+len(s.mcpServers))
	// add the welcome item
	data = append(data, NewContentWelcome().menuItem())

	// create the content for each mcp server
	for _, server := range s.mcpServers {
		data = append(data, NewContentMcpServer(
			window,
			server,
			s.AsServerListActions(),
		).menuItem())
	}

	// add the action to add a new mcp server
	data = append(data, NewContentAddMcpServer().menuItem())

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

func (s *SideMenu) ValidateNewName(name string) error {
	for _, server := range s.mcpServers {
		if server.Name == name {
			return errors.New("a server with this name already exists")
		}
	}
	return nil
}

func (s *SideMenu) ValidateExistingName(uuid string) func(name string) error {
	return func(name string) error {
		// Check for empty string
		if len(name) == 0 {
			return errors.New("name cannot be empty")
		}

		// Check if name starts with a number
		if name[0] >= '0' && name[0] <= '9' {
			return errors.New("name cannot start with a number")
		}

		// Check if name contains only alphanumeric characters and underscore
		for _, char := range name {
			if !((char >= 'a' && char <= 'z') ||
				(char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') ||
				char == '_') {
				return errors.New("name can only contain letters, numbers, and underscores")
			}
		}

		// Check for duplicate names
		for _, server := range s.mcpServers {
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
