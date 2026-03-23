package cmd

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	MessageSyntax           = "Syntax error: Unexpected \"%s\": at \"%s >>%s<<\""
	MessageBooleanInvalid   = "Syntax error: \"%s\" is not a valid boolean"
	MessageNumberInvalid    = "Syntax error: \"%s\" is not a valid number"
	MessageParameterInvalid = "Syntax error: Unexpected \"%s\""
	MessagePlayerNotFound   = "Syntax error: Player \"%s\" not found"
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
	if len(line.args) == 0 {
		if optional {
			v.FieldByName("Has").SetBool(false)
			return nil
		}
		return fmt.Errorf(MessageSyntax, "", strings.Join(line.seen, " "), name)
	}

	arg, _ := line.Next()
	var err error

	target := v
	if _, ok := v.Interface().(optionalT); ok {
		field := v.FieldByName("Value")
		target = reflect.New(field.Type()).Elem()
	}

	i := target.Addr().Interface()
	switch res := i.(type) {
	case *SubCommand:
		if !strings.EqualFold(arg, name) {
			return fmt.Errorf(MessageParameterInvalid, arg)
		}
	case *int, *int8, *int16, *int32, *int64:
		val, parseErr := strconv.ParseInt(arg, 10, target.Type().Bits())
		if parseErr != nil {
			return fmt.Errorf(MessageNumberInvalid, arg)
		}
		target.SetInt(val)
	case *uint, *uint8, *uint16, *uint32, *uint64:
		val, parseErr := strconv.ParseUint(arg, 10, target.Type().Bits())
		if parseErr != nil {
			return fmt.Errorf(MessageNumberInvalid, arg)
		}
		target.SetUint(val)
	case *float32, *float64:
		val, parseErr := strconv.ParseFloat(arg, target.Type().Bits())
		if parseErr != nil {
			return fmt.Errorf(MessageNumberInvalid, arg)
		}
		target.SetFloat(val)
	case *bool:
		val, parseErr := strconv.ParseBool(arg)
		if parseErr != nil {
			return fmt.Errorf(MessageBooleanInvalid, arg)
		}
		target.SetBool(val)
	case *string:
		target.SetString(arg)
	case *Varargs:
		target.SetString(strings.Join(line.Leftover(), " "))
		p.finalizeOptional(v, target)
		return nil
	case *[]Target:
		targets, targetErr := p.parseTargets(line, arg)
		if targetErr != nil {
			return targetErr
		}
		target.Set(reflect.ValueOf(targets))
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
				return fmt.Errorf(MessageParameterInvalid, arg)
			}
		}
	}

	if err == nil {
		line.RemoveNext()
		p.finalizeOptional(v, target)
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

// finalizeOptional is a helper function that wraps the parsed target value into an
// Optional[T] structure, setting the 'Has' flag to true.
func (p Parser) finalizeOptional(v reflect.Value, target reflect.Value) {
	if o, ok := v.Interface().(optionalT); ok {
		v.Set(reflect.ValueOf(o.with(target.Interface())))
	}
}

// parseTargets parses one or more targets from the provided argument, it supports Minecraft selectors like
// @s, @p, @a, @r, @e and falls back to player name lookup, returning compliant error messages if no targets are found.
func (p Parser) parseTargets(line *Line, arg string) ([]Target, error) {
	var targets []Target
	switch arg {
	case "@s":
		if t, ok := line.src.(Target); ok {
			targets = append(targets, t)
		}
	case "@p", "@a", "@r", "@e":
		for _, entry := range line.players {
			targets = append(targets, PlayerTarget{NameValue: entry.Username})
		}
		if arg == "@p" || arg == "@r" {
			if len(targets) > 0 {
				targets = targets[:1]
			}
		}
	default:
		t, found := p.lookupPlayer(line, arg)
		if !found {
			return nil, fmt.Errorf(MessagePlayerNotFound, arg)
		}
		targets = append(targets, t)
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf(MessagePlayerNotFound, arg)
	}
	return targets, nil
}
