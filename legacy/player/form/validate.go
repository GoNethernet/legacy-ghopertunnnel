package form

import (
	"fmt"
)

// Validate checks if the form is valid.
func Validate(f interface{}) error {
	switch ft := f.(type) {
	case CustomForm:
		for _, e := range ft.Elements() {
			switch e.(type) {
			case Button, *Button:
				return fmt.Errorf("form validate: button not allowed in custom form")
			}
		}
		return nil
	case SimpleForm:
		for _, e := range ft.Elements() {
			switch e.(type) {
			case Button, *Button, Label, *Label:
				continue
			default:
				return fmt.Errorf("form validate: '%v' element not allowed in simple form", e.Name())
			}
		}
		return nil
	}
	return nil
}
