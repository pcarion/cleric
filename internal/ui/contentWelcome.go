package ui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ContentWelcome struct {
	myApp   fyne.App
	version string
}

func NewContentWelcome(myApp fyne.App, version string) *ContentWelcome {
	return &ContentWelcome{myApp: myApp, version: version}
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
				container.NewCenter(container.NewHBox(githubbutton)),
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
