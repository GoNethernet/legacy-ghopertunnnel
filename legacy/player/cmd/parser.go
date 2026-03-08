package cmd

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// Line represents a command line holding arguments passed during execution.
type Line struct {
	args    []string
	seen    []string
	src     Source
	players []protocol.PlayerListEntry
}

// NewLine creates a new Line instance with the provided arguments and source.
func NewLine(args []string, src Source, players []protocol.PlayerListEntry) *Line {
	return &Line{args: args, src: src, players: players}
}

// Next returns the next argument from the line without removing it.
func (line *Line) Next() (string, bool) {
	if len(line.args) == 0 {
		return "", false
	}
	return line.args[0], true
}

// RemoveNext consumes the next argument from the command line.
func (line *Line) RemoveNext() {
	if len(line.args) > 0 {
		line.seen = append(line.seen, line.args[0])
		line.args = line.args[1:]
	}
}

// Leftover returns all remaining arguments and clears the line.
func (line *Line) Leftover() []string {
	v := line.args
	line.args = nil
	return v
}

// Parser manages the conversion of raw arguments into structured data.
type Parser struct{}

// ParseArgument handles the conversion of a single command line argument into a Go reflect.Value.
func (p Parser) ParseArgument(line *Line, v reflect.Value, optional bool, name string) error {
	if len(line.args) == 0 && optional {
		v.FieldByName("Has").SetBool(false)
		return nil
	}

	arg, ok := line.Next()
	if !ok {
		return fmt.Errorf("Syntax error: Unexpected \"\": at >>%s<<", name)
	}

	var err error
	target := v
	if _, ok := v.Interface().(optionalT); ok {
		field := v.FieldByName("Value")
		if !field.IsValid() {
			panic("field Value not found in Optional struct")
		}
		target = reflect.New(field.Type()).Elem()
	}

	if !target.IsValid() || !target.CanAddr() {
		panic("field " + name + " is not addressable.")
	}

	i := target.Addr().Interface()
	switch res := i.(type) {
	case *SubCommand:
		if !strings.EqualFold(arg, name) {
			err = fmt.Errorf("Syntax error: Unexpected \"%s\": at \"%s\"", arg, arg)
		}
	case *[]Target:
		var targets []Target
		switch arg {
		case "@s", "@p", "@r", "@a", "@e":
			if t, ok := line.src.(Target); ok {
				targets = append(targets, t)
			} else {
				err = fmt.Errorf("Syntax error: Source is not a valid target for selector %s", arg)
			}
		default:
			t, found := p.lookupPlayer(line, arg)
			if !found {
				return fmt.Errorf("Syntax error: Player \"%s\" not found", arg)
			}
			targets = append(targets, t)
		}
		if err == nil {
			target.Set(reflect.ValueOf(targets))
		}
	case *int, *int8, *int16, *int32, *int64:
		var val int64
		val, err = strconv.ParseInt(arg, 10, target.Type().Bits())
		if err == nil {
			target.SetInt(val)
		}
	case *uint, *uint8, *uint16, *uint32, *uint64:
		var val uint64
		val, err = strconv.ParseUint(arg, 10, target.Type().Bits())
		if err == nil {
			target.SetUint(val)
		}
	case *float32, *float64:
		var val float64
		val, err = strconv.ParseFloat(arg, target.Type().Bits())
		if err == nil {
			target.SetFloat(val)
		}
	case *string:
		target.SetString(arg)
	case *bool:
		var val bool
		val, err = strconv.ParseBool(arg)
		if err == nil {
			target.SetBool(val)
		}
	case *Varargs:
		target.SetString(strings.Join(line.Leftover(), " "))
		if o, ok := v.Interface().(optionalT); ok {
			v.Set(reflect.ValueOf(o.with(target.Interface())))
		}
		return nil
	default:
		if enum, ok := res.(Enum); ok {
			opts := enum.Options()
			found := false
			for _, opt := range opts {
				if strings.EqualFold(opt, arg) {
					target.SetString(opt)
					found = true
					break
				}
			}
			if !found {
				err = fmt.Errorf("Syntax error: Unexpected \"%s\": at \"%s\"", arg, arg)
			}
		}
	}

	if err == nil {
		line.RemoveNext()
		if o, ok := v.Interface().(optionalT); ok {
			v.Set(reflect.ValueOf(o.with(target.Interface())))
		}
	}
	return err
}

// lookupPlayer search up for a players by his name and returns a target.
func (p Parser) lookupPlayer(line *Line, name string) (Target, bool) {
	for _, entry := range line.players {
		if strings.EqualFold(entry.Username, name) {
			return PlayerTarget{NameValue: entry.Username}, true
		}
	}
	return nil, false
}
