package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ClaudeAction interface {
	IsServerInClaude() bool
	AddToClaude()
	RemoveFromClaude()
}

// is a toolbar item that allows to send a message to claude
type ToolbarClaudeAction struct {
	claudeAction ClaudeAction
	statusLabel  *widget.Label
}

func NewToolbarClaudeAction(claudeAction ClaudeAction, statusLabel *widget.Label) *ToolbarClaudeAction {
	return &ToolbarClaudeAction{claudeAction: claudeAction, statusLabel: statusLabel}
}

func (t *ToolbarClaudeAction) ToolbarObject() fyne.CanvasObject {
	if t.claudeAction.IsServerInClaude() {
		return NewHoverButton(theme.CheckButtonCheckedIcon(), "Remove from claude", func() {
			t.claudeAction.RemoveFromClaude()
		}, func() {
			t.statusLabel.SetText("Remove from the list of claude servers")
		}, func() {
			t.statusLabel.SetText("")
		})
	} else {
		return NewHoverButton(theme.CheckButtonIcon(), "Add to claude", func() {
			t.claudeAction.AddToClaude()
		}, func() {
			t.statusLabel.SetText("Add to the list of claude servers")
		}, func() {
			t.statusLabel.SetText("")
		})
	}
}

type ToolbarEditAction interface {
	IsEditMode() bool
	EditMode()
	CancelEditMode()
}

type ToolbarEdit struct {
	action      ToolbarEditAction
	statusLabel *widget.Label
}

func NewEditToolbar(action ToolbarEditAction, statusLabel *widget.Label) *ToolbarEdit {
	return &ToolbarEdit{action: action, statusLabel: statusLabel}
}

func (t *ToolbarEdit) ToolbarObject() fyne.CanvasObject {
	if t.action.IsEditMode() {
		return NewHoverButton(theme.CancelIcon(), "Exit edit mode", func() {
			t.action.CancelEditMode()
		}, func() {
			t.statusLabel.SetText("Exit edit mode")
		}, func() {
			t.statusLabel.SetText("")
		})
	} else {
		return NewHoverButton(theme.DocumentCreateIcon(), "", func() {
			t.action.EditMode()
		}, func() {
			t.statusLabel.SetText("Enter edit mode")
		}, func() {
			t.statusLabel.SetText("")
		})
	}
}
