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
	action ToolbarEditAction
}

func NewEditToolbar(action ToolbarEditAction) *ToolbarEdit {
	return &ToolbarEdit{action: action}
}

func (t *ToolbarEdit) ToolbarObject() fyne.CanvasObject {
	if t.action.IsEditMode() {
		return widget.NewButtonWithIcon("Exit edit mode", theme.CancelIcon(), func() {
			t.action.CancelEditMode()
		})
	} else {
		return widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			t.action.EditMode()
		}).ToolbarObject()
	}
}
