package form

import "encoding/json"

// Element represent an element in a form that will be json marshalled.
type Element interface {
	// Name returns the name of the element.
	Name() string
	json.Marshaler
}

// Button is a clickable button in the form, multiple buttons cannot be used in the CustomForm
// since there's already a 'Submit' button.
type Button struct {
	// Text ...
	Text string
	// Image is an image that will be present next to the button only in the simple form.
	Image string
	// Submit is when the player clicks the button.
	Submit func(submitter Submitter)
}

// Name ...
func (b Button) Name() string {
	return "button"
}

// MarshalJSON ...
func (b Button) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Text  string `json:"text"`
		Image string `json:"image,omitempty"`
	}{Text: b.Text, Image: b.Image})
}

// Label is an additional text in the form.
type Label struct {
	// Text ...
	Text string
}

// Name ...
func (l Label) Name() string {
	return "label"
}

// MarshalJSON ...
func (l Label) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": "label",
		"text": l.Text,
	})
}

// Toggle is an element property that allows only a true or false option.
type Toggle struct {
	// Text ...
	Text string
	// Default is the default value of the toggle, usually it is set to false.
	Default bool
	// Value is the final value of the toggle.
	Value func(bool)
}

// Name ...
func (t Toggle) Name() string {
	return "toggle"
}

// MarshalJSON ...
func (t Toggle) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "toggle",
		"text":    t.Text,
		"default": t.Default,
	})
}

// Input is an elements that holds an input function.
type Input struct {
	// Text ...
	Text string
	// Placeholder is the returned text if the input is empty.
	Placeholder string
	// Default ...
	Default string
	// Final is the final selected value.
	Final func(string)
}

// Name ...
func (i Input) Name() string {
	return "input"
}
func (i Input) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":        "input",
		"text":        i.Text,
		"placeholder": i.Placeholder,
		"default":     i.Default,
	})
}

// Slider is an element that holds a range of numbers.
type Slider struct {
	// Text ...
	Text string
	// Min ...
	Min float32
	// Max ...
	Max float32
	// Step ...
	Step float32
	// Default ...
	Default float32
	// Selected is the final value of the slider.
	Selected func(float32)
}

// Name ...
func (s Slider) Name() string {
	return "slider"
}

// MarshalJSON ...
func (s Slider) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "slider",
		"text":    s.Text,
		"min":     s.Min,
		"max":     s.Max,
		"step":    s.Step,
		"default": s.Default,
	})
}

// Dropdown is an element that holds a list of options.
type Dropdown struct {
	// Text ...
	Text string
	// Options ...
	Options []string
	// DefaultIndex ...
	DefaultIndex int32
	// Selected is the final value of the dropdown.
	Selected func(string)
}

// Name ...
func (d Dropdown) Name() string {
	return "dropdown"
}

// MarshalJSON ...
func (d Dropdown) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "dropdown",
		"text":    d.Text,
		"options": d.Options,
		"default": d.DefaultIndex,
	})
}
