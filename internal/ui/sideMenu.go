package ui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"github.com/pcarion/cleric/internal/configuration"
)

type setMainContentFunc func(mainContent *MainContent)

type menuItem interface {
	label() string
	icon() fyne.Resource
	content() *MainContent
}

// definition of the side menu
type SideMenu struct {
	window             fyne.Window
	config             *configuration.Configuration
	sideMenuData       []menuItem
	mcpServers         []*configuration.McpServerDescription
	list               *widget.List
	refreshMainContent func()
}

func NewSideMenu(window fyne.Window, refreshMainContent func()) *SideMenu {
	config := configuration.LoadConfiguration()
	mcpServers := config.LoadMcpServers()

	return &SideMenu{
		config:             config,
		mcpServers:         mcpServers,
		refreshMainContent: refreshMainContent,
		window:             window,
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
		s.refreshSideMenuData()
		s.list.Refresh()
	}
}

func (s *SideMenu) refreshSideMenuData() {
	// use the mcp servers as the data
	data := make([]menuItem, 0, 1+len(s.mcpServers))
	// add the welcome item
	data = append(data, NewContentWelcome().menuItem())

	// create the content for each mcp server
	for _, server := range s.mcpServers {
		data = append(data, NewContentMcpServer(
			s.window,
			server,
			s.AsServerListActions(),
		).menuItem())
	}

	// add the action to add a new mcp server
	data = append(data, NewContentAddMcpServer(s.AsServerListActions()).menuItem())

	s.sideMenuData = data
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
) fyne.CanvasObject {

	// refresh the side menu data
	s.refreshSideMenuData()
	s.list = widget.NewList(
		func() int {
			return len(s.sideMenuData)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()), // 0
				widget.NewLabel("Template Object"),   // 1
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			// get the side menu item
			sideMenuItem := s.sideMenuData[id]
			hbox := item.(*fyne.Container) // that's the Hbox

			hbox.Objects[0].(*widget.Icon).SetResource(sideMenuItem.icon())
			hbox.Objects[1].(*widget.Label).SetText(sideMenuItem.label())
		},
	)

	// add code if the item is clicked
	s.list.OnSelected = func(id widget.ListItemID) {
		sideMenuItem := s.sideMenuData[id]
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

func (s *SideMenu) ValidateNewMcpServerName(name string) error {
	for _, server := range s.mcpServers {
		if server.Name == name {
			return errors.New("a server with this name already exists")
		}
	}
	return nil
}

func (s *SideMenu) ValidateExistingMcpServerName(uuid string) func(name string) error {
	return func(name string) error {
		if !isValidServerName(name) {
			return errors.New("name cannot be empty")
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

func (s *SideMenu) DeleteMcpServer(uuid string) {
	for i, server := range s.mcpServers {
		if server.Uuid == uuid {
			s.mcpServers = append(s.mcpServers[:i], s.mcpServers[i+1:]...)
			break
		}
	}
	s.SaveMcpServers()
}

func (s *SideMenu) AddMcpServer(name string) (string, error) {
	if !isValidServerName(name) {
		return "", errors.New("name must be an alphanumeric string")
	}
	uuid := uuid.New().String()

	s.mcpServers = append(s.mcpServers, &configuration.McpServerDescription{
		Name:            name,
		Uuid:            uuid,
		Description:     "",
		InConfiguration: false,
		Configuration: configuration.McpServerConfiguration{
			Command: "",
			Args:    []string{},
			Env:     map[string]string{},
		},
	})
	s.SaveMcpServers()
	return uuid, nil
}

func (s *SideMenu) ResetListScroll() {
	if s.list != nil && len(s.sideMenuData) > 0 {
		s.list.Select(0)
		s.list.ScrollToTop()
	}
}

func (s *SideMenu) ResetListToServer(uuid string) {
	if s.list != nil && len(s.sideMenuData) > 0 {
		for i, server := range s.sideMenuData {
			if server.content().ContentId == uuid {
				s.list.Select(i)
				break
			}
		}
	}
}
