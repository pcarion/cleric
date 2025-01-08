package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func AddServerDialog(window fyne.Window, validator func(name string) error, actions ServerListActions) {
	nameEntry := widget.NewEntry()
	hintText := "must be an alphanumeric string, not already used for another server"
	nameEntry.Validator = validator
	nameEntry.SetText("")

	dialog := dialog.NewForm(
		"Add a new MCP Server",
		"Save",
		"Cancel",
		[]*widget.FormItem{
			{Text: "Server Name", Widget: nameEntry, HintText: hintText},
		},
		func(confirm bool) {
			if confirm {
				name := nameEntry.Text
				err := validator(name)
				if err != nil {
					dialog.ShowError(err, window)
				} else {
					serverUuid, err := actions.AddMcpServer(nameEntry.Text)
					if err != nil {
						dialog.ShowError(err, window)
					} else {
						actions.RefreshSideMenu()
						// we select the newly added server in the list
						actions.ResetListToContentId(serverUuid)
					}

				}
			}
		},
		window,
	)
	dialog.Resize(fyne.NewSize(600, 300))
	dialog.Show()

}
