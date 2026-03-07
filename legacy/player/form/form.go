package form

// SimpleForm represent an easy and quick to use form, some elements here such as Slider and Dropdown cannot be added here otherwise an error will
// occur.
type SimpleForm interface {
	// Title ...
	Title() string
	// Elements ...
	Elements() []Element
}

// CustomForm is a way more customizable version of SimpleForm, here every element can be added except Button.
type CustomForm interface {
	// Title ...
	Title() string
	// Elements ...
	Elements() []Element
	// Submit is when the player click the 'submit' button.
	Submit(submitter Submitter)
}

// ModalForm is a form with only two buttons and a message.
type ModalForm interface {
	// Title ...
	Title() string
	// Message ...
	Message() string
	// Buttons ...
	Buttons() [2]Button
}

// ElementsToButtons converts a slice of Element interfaces into a slice of Button structs, it filters the input by checking if each
// element is a Button or a pointer to one.
func ElementsToButtons(elements []Element) []Button {
	var buttons []Button
	for _, e := range elements {
		if b, ok := e.(Button); ok {
			buttons = append(buttons, b)
		} else if b, ok := e.(*Button); ok {
			buttons = append(buttons, *b)
		}
	}
	return buttons
}
