package main

import "fmt"

type Button interface {
	Render()
}

type Checkbox interface {
	Render()
}

type WindowsButton struct{}

func (w *WindowsButton) Render() {
	fmt.Println("Rendering Windows Button")
}

type WindowsCheckbox struct{}

func (w *WindowsCheckbox) Render() {
	fmt.Println("Rendering Windows Checkbox")
}

type MacButton struct{}

func (m *MacButton) Render() {
	fmt.Println("Rendering Mac Button")
}

type MacCheckbox struct{}

func (m *MacCheckbox) Render() {
	fmt.Println("Rendering Mac Checkbox")
}

type UIFactory interface {
	CreateButton() Button
	CreateCheckbox() Checkbox
}

type WindowsFactory struct{}

func (w *WindowsFactory) CreateButton() Button {
	return &WindowsButton{}
}

func (w *WindowsFactory) CreateCheckbox() Checkbox {
	return &WindowsCheckbox{}
}

type MacFactory struct{}

func (m *MacFactory) CreateButton() Button {
	return &MacButton{}
}

func (m *MacFactory) CreateCheckbox() Checkbox {
	return &MacCheckbox{}
}

func main() {
	var factory UIFactory

	os := "windows" // change to "mac"

	if os == "windows" {
		factory = &WindowsFactory{}
	} else {
		factory = &MacFactory{}
	}

	button := factory.CreateButton()
	checkbox := factory.CreateCheckbox()

	button.Render()
	checkbox.Render()
}
