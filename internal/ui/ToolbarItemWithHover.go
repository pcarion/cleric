package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// ToolbarItemWithHover represents a toolbar item that can respond to hover events
type ToolbarItemWithHover struct {
	Icon        fyne.Resource // The content to display in the toolbar item
	OnActivated func()        // Function to call when mouse enters
	Hoverable   ToolbarHoverable
}

type ToolbarHoverable interface {
	onHoverIn()
	onHoverOut()
}

// NewToolbarItemWithHover creates a new toolbar item with hover functionality
func NewToolbarItemWithHover(icon fyne.Resource, onActivated func(), hoverable ToolbarHoverable) *ToolbarItemWithHover {
	return &ToolbarItemWithHover{
		Icon:        icon,
		OnActivated: onActivated,
		Hoverable:   hoverable,
	}
}

func (t *ToolbarItemWithHover) ToolbarObject() fyne.CanvasObject {
	button := NewHoverButton(t.Icon, "", t.OnActivated, t.Hoverable.onHoverIn, t.Hoverable.onHoverOut)
	return button
}

type ImplementedToolbarHoverable struct {
	doOnHoverIn  func()
	doOnHoverOut func()
}

func (h *ImplementedToolbarHoverable) onHoverIn()  { h.doOnHoverIn() }
func (h *ImplementedToolbarHoverable) onHoverOut() { h.doOnHoverOut() }

func mkHoverable(helpText string, label *widget.Label) ToolbarHoverable {
	imp := ImplementedToolbarHoverable{
		doOnHoverIn: func() {
			label.SetText(helpText)
		},
		doOnHoverOut: func() {
			label.SetText("")
		},
	}
	return &imp
}
