package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// TextWithCopy is a widget that displays text with a copy button
// NewLabelWithCopy creates a new label with copy button
func NewTextWithCopy(text string, window fyne.Window) *fyne.Container {
	label := widget.NewMultiLineEntry()
	label.SetText(text)
	label.Wrapping = fyne.TextWrapWord

	// button := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
	// 	window.Clipboard().SetContent(label.Text)
	// })
	content := container.NewHBox(
		label,
		// button,
	)

	return content
}
