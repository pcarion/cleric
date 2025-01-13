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
}

func NewToolbarClaudeAction(claudeAction ClaudeAction) *ToolbarClaudeAction {
	return &ToolbarClaudeAction{claudeAction: claudeAction}
}

func (t *ToolbarClaudeAction) ToolbarObject() fyne.CanvasObject {
	if t.claudeAction.IsServerInClaude() {
		return widget.NewButtonWithIcon("Remove from claude", theme.CheckButtonCheckedIcon(), func() {
			t.claudeAction.RemoveFromClaude()
		})
	} else {
		return widget.NewButtonWithIcon("Add to claude", theme.CheckButtonIcon(), func() {
			t.claudeAction.AddToClaude()
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
