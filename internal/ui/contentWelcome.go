package ui

import (
	"errors"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/hashicorp/go-version"
)

type ContentWelcome struct {
	window             fyne.Window
	myApp              fyne.App
	version            string
	buttonCheckVersion *widget.Button
}

func NewContentWelcome(myApp fyne.App, window fyne.Window, version string) *ContentWelcome {
	return &ContentWelcome{myApp: myApp, window: window, version: version}
}

func (c *ContentWelcome) menuItem() menuItem {
	return c
}

func (c *ContentWelcome) content() *MainContent {
	return &MainContent{
		ContentId: "welcome",
		View: func(window fyne.Window) fyne.CanvasObject {
			img := canvas.NewImageFromResource(resourceIconPng)
			img.SetMinSize(fyne.NewSize(100, 100))
			img.FillMode = canvas.ImageFillContain

			mdContent := `
# Cleric 

Configure your Claude Desktop settings for MCP servers through a simple GUI

---
## Author ` + `
Pierre Carion - pcarion@gmail.com


## License
MIT


## Version` + `
` + c.version
			richhead := widget.NewRichTextFromMarkdown(mdContent)

			for i := range richhead.Segments {
				if seg, ok := richhead.Segments[i].(*widget.TextSegment); ok {
					seg.Style.Alignment = fyne.TextAlignCenter
				}
				if seg, ok := richhead.Segments[i].(*widget.HyperlinkSegment); ok {
					seg.Alignment = fyne.TextAlignCenter
				}
			}
			githubbutton := widget.NewButton("Github page", func() {
				go func() {
					u, _ := url.Parse("https://github.com/pcarion/cleric")
					_ = fyne.CurrentApp().OpenURL(u)
				}()
			})

			checkversion := widget.NewButton("Check version", func() {
				go c.checkVersion()
			})
			c.buttonCheckVersion = checkversion

			settings := widget.NewLabel("Appearance")
			settings.Alignment = fyne.TextAlignCenter
			settings.TextStyle.Bold = true

			themes := container.NewGridWithColumns(2,
				widget.NewButton("Dark", func() {
					SetDarkTheme(c.myApp)
				}),
				widget.NewButton("Light", func() {
					SetLightTheme(c.myApp)
				}),
			)

			bottom := container.NewHBox(
				layout.NewSpacer(),
				container.NewCenter(settings),
				container.NewCenter(themes),
			)

			cont := container.NewVBox(
				container.NewCenter(img),
				container.NewCenter(richhead),
				container.NewCenter(container.NewHBox(githubbutton, checkversion)),
			)
			mainContentContainer := container.NewBorder(cont, bottom, nil, nil)

			return mainContentContainer
		},
	}
}

func (c *ContentWelcome) label() string {
	return "Welcome"
}

func (c *ContentWelcome) icon() fyne.Resource {
	return theme.HomeIcon()
}

func (c *ContentWelcome) checkVersion() {
	c.buttonCheckVersion.Disable()
	defer c.buttonCheckVersion.Enable()
	errInvalidVersion := errors.New("invalid version " + c.version)

	versionCurrent, err := version.NewVersion(c.version)
	if err != nil {
		dialog.ShowError(errInvalidVersion, c.window)
		return
	}
	errRedirectChecker := errors.New("redirect")
	errVersionGet := errors.New("failed to get version info - check your internet connection")

	req, err := http.NewRequest("GET", "https://github.com/pcarion/cleric/releases/latest", nil)
	if err != nil {
		dialog.ShowError(errVersionGet, c.window)
		return
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errRedirectChecker
		},
	}

	response, err := client.Do(req)
	if err != nil && !errors.Is(err, errRedirectChecker) {
		dialog.ShowError(errVersionGet, c.window)
		return
	}

	defer response.Body.Close()

	if errors.Is(err, errRedirectChecker) {
		responceUrl, err := response.Location()
		if err != nil {
			dialog.ShowError(errVersionGet, c.window)
			return
		}
		str := strings.Trim(filepath.Base(responceUrl.Path), "v")
		versionRelease, err := version.NewVersion(str)
		if err != nil {
			dialog.ShowError(errVersionGet, c.window)
			return
		}

		switch {
		case versionRelease.GreaterThan(versionCurrent):
			dialog.ShowInformation("Version checker", "There is a new version available on Github: "+versionRelease.String(), c.window)
			return
		default:
			dialog.ShowInformation("Version checker", "No new version", c.window)
			return
		}
	}

	dialog.ShowError(errVersionGet, c.window)
}
