package cmd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/permission"
	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/position"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// Command defines the interface that all custom commands must implement.
type Command interface {
	// Name returns the primary name of the command.
	Name() string
	// Description returns a short explanation of the command's purpose.
	Description() string
	// Aliases returns alternative names that trigger the same command.
	Aliases() []string
	// PermissionLevel defines the minimum rank required to execute the command.
	PermissionLevel() permission.Permission
	// Run contains the logic to be executed when the command is called.
	Run(src Source)
}

// RegisteredCommand stores the reflection type of command for instantiation.
type RegisteredCommand struct {
	Type reflect.Type
}

// CustomCommands is a global registry of all commands available in the proxy.
var CustomCommands = make(map[string]RegisteredCommand)

// NewCommand translates a Go struct into Minecraft protocol command overloads.
func NewCommand(t reflect.Type, enums *[]protocol.CommandEnum, enumValues *[]string) []protocol.CommandOverload {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	var params []protocol.CommandParameter
	enumIndices := make(map[string]uint32)
	for i, e := range *enums {
		enumIndices[e.Type] = uint32(i)
	}

	foundOptional := false
	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		var tpe uint32
		isOptional := false
		fieldType := field.Type

		if strings.HasPrefix(fieldType.Name(), "Optional[") {
			isOptional = true
			foundOptional = true
		} else if foundOptional {
			panic(fmt.Sprintf("command structure: non-optional field '%s' in %s found after an optional field. optionals must be at the end", field.Name, t.Name()))
		}

		if isOptional {
			fieldType = fieldType.Field(0).Type
		}

		var enumType string
		var enumOptions []string

		val := reflect.New(fieldType)
		fieldValue := val.Interface()

		switch fv := fieldValue.(type) {
		case *SubCommand:
			enumType = "SubCommand" + field.Name
			enumOptions = []string{strings.ToLower(field.Name)}
		case *[]Target:
			tpe = protocol.CommandArgTypeTarget
		case *position.Position, *mgl32.Vec3, *mgl64.Vec3:
			tpe = protocol.CommandArgTypePosition
		case *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64:
			tpe = protocol.CommandArgTypeInt
		case *float32, *float64:
			tpe = protocol.CommandArgTypeFloat
		case *string:
			tpe = protocol.CommandArgTypeString
		case *bool:
			enumType = "boolean"
			enumOptions = []string{"true", "false"}
		case *Varargs:
			tpe = protocol.CommandArgTypeRawText
		default:
			if enum, ok := fv.(Enum); ok {
				enumType = enum.Type()
				enumOptions = enum.Options()
			} else {
				tpe = protocol.CommandArgTypeValue
			}
		}

		if enumOptions != nil {
			index, ok := enumIndices[enumType]
			if !ok {
				index = uint32(len(*enums))
				enumIndices[enumType] = index
				protoEnum := protocol.CommandEnum{Type: enumType}
				for _, opt := range enumOptions {
					var valIndex uint32
					foundVal := false
					for idx, v := range *enumValues {
						if v == opt {
							valIndex = uint32(idx)
							foundVal = true
							break
						}
					}
					if !foundVal {
						valIndex = uint32(len(*enumValues))
						*enumValues = append(*enumValues, opt)
					}
					protoEnum.ValueIndices = append(protoEnum.ValueIndices, valIndex)
				}
				*enums = append(*enums, protoEnum)
			}
			tpe = protocol.CommandArgEnum | index
		}

		params = append(params, protocol.CommandParameter{
			Name:     strings.ToLower(field.Name),
			Type:     tpe | protocol.CommandArgValid,
			Optional: isOptional,
		})
	}
	return []protocol.CommandOverload{{Parameters: params}}
}
