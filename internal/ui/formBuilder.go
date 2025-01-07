package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// utility class to build a form

type FormBuilder struct {
	elements []fyne.CanvasObject
}

func NewFormBuilder() *FormBuilder {
	return &FormBuilder{}
}

func (f *FormBuilder) AddField(label *widget.Label, controls *fyne.Container) {
	f.elements = append(f.elements, label, controls)
}

func (f *FormBuilder) GetContainer() *fyne.Container {
	return container.New(
		layout.NewFormLayout(),
		f.elements...,
	)
}
