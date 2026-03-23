package cmd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/permission"
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
	Types []reflect.Type
}

// CustomCommands is a global registry of all commands available in the proxy.
var CustomCommands = make(map[string]RegisteredCommand)

// NewCommand translates a Go struct into Minecraft protocol command overloads.
func NewCommand(types []reflect.Type, enums *[]protocol.CommandEnum, enumValues *[]string, dynamicEnums *[]protocol.DynamicEnum) []protocol.CommandOverload {
	var overloads []protocol.CommandOverload
	enumIndices := make(map[string]uint32)
	for i, e := range *enums {
		enumIndices[e.Type] = uint32(i)
	}
	dynamicEnumIndices := make(map[string]uint32)
	for i, e := range *dynamicEnums {
		dynamicEnumIndices[e.Type] = uint32(i)
	}

	for _, t := range types {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		var params []protocol.CommandParameter
		foundOptional := false

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			var tpe uint32
			isOptional := strings.HasPrefix(field.Type.Name(), "Optional[")
			if isOptional {
				foundOptional = true
			} else if foundOptional {
				panic(fmt.Sprintf("non-optional field %s after optional", field.Name))
			}

			fieldType := field.Type
			if isOptional {
				fieldType = fieldType.Field(0).Type
			}

			var enumType string
			var enumOptions []string
			var isSoft bool

			val := reflect.New(fieldType)
			fv := val.Interface()

			switch fv.(type) {
			case *SubCommand:
				enumType = "SubCommand" + field.Name
				enumOptions = []string{strings.ToLower(field.Name)}
			case *[]Target:
				tpe = protocol.CommandArgTypeTarget
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
					if se, ok := fv.(SoftEnum); ok && se.Soft() {
						isSoft = true
					}
				} else {
					tpe = protocol.CommandArgTypeValue
				}
			}

			if enumOptions != nil {
				if isSoft {
					index, ok := dynamicEnumIndices[enumType]
					if !ok {
						index = uint32(len(*dynamicEnums))
						dynamicEnumIndices[enumType] = index
						*dynamicEnums = append(*dynamicEnums, protocol.DynamicEnum{
							Type:   enumType,
							Values: enumOptions,
						})
					}
					tpe = protocol.CommandArgSoftEnum | index
				} else {
					index, ok := enumIndices[enumType]
					if !ok {
						index = uint32(len(*enums))
						enumIndices[enumType] = index
						protoEnum := protocol.CommandEnum{Type: enumType}
						for _, opt := range enumOptions {
							valIndex := uint32(len(*enumValues))
							*enumValues = append(*enumValues, opt)
							protoEnum.ValueIndices = append(protoEnum.ValueIndices, valIndex)
						}
						*enums = append(*enums, protoEnum)
					}
					tpe = protocol.CommandArgEnum | index
				}
			}

			params = append(params, protocol.CommandParameter{
				Name:     strings.ToLower(field.Name),
				Type:     tpe | protocol.CommandArgValid,
				Optional: isOptional,
			})
		}
		overloads = append(overloads, protocol.CommandOverload{Parameters: params})
	}
	return overloads
}
