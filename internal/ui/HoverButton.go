package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// HoverButton represents a button that can respond to hover events
type HoverButton struct {
	widget.Button
	OnHoverIn  func() // Function to call when mouse enters
	OnHoverOut func() // Function to call when mouse leaves
}

// NewHoverButton creates a new button with hover functionality
func NewHoverButton(icon fyne.Resource, label string, onTapped func(), onHoverIn func(), onHoverOut func()) *HoverButton {
	button := &HoverButton{
		OnHoverIn:  onHoverIn,
		OnHoverOut: onHoverOut,
	}
	button.Icon = icon
	if label != "" {
		button.Text = label
	}
	button.Importance = widget.LowImportance
	button.OnTapped = onTapped
	button.ExtendBaseWidget(button) // Important: Required for custom widgets
	return button
}

// MouseIn handles the mouse enter event
func (b *HoverButton) MouseIn(e *desktop.MouseEvent) {
	b.Button.MouseIn(e) // Call parent implementation
	if b.OnHoverIn != nil {
		b.OnHoverIn()
	}
}

// MouseOut handles the mouse exit event
func (b *HoverButton) MouseOut() {
	b.Button.MouseOut() // Call parent implementation
	if b.OnHoverOut != nil {
		b.OnHoverOut()
	}
}

// MouseMoved handles mouse movement over the button
func (b *HoverButton) MouseMoved(e *desktop.MouseEvent) {
	b.Button.MouseMoved(e) // Call parent implementation
}
