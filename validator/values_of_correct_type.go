package validator

import (
	"fmt"

	"github.com/vektah/gqlparser"
)

func init() {
	addRule("ValuesOfCorrectType", func(observers *Events, addError addErrFunc) {
		observers.OnValue(func(walker *Walker, valueType gqlparser.Type, def *gqlparser.Definition, value gqlparser.Value) {
			if def == nil || valueType == nil {
				fmt.Println("BADLANDS")
				return
			}

			switch def.Kind {
			case gqlparser.Enum:
				rawVal, err := value.Value(nil)
				if err != nil {
					addError(Message("Expected type %s, found %s %s.", valueType.String(), value, err.Error()))
				}

				var possible []string
				for _, val := range def.Values {
					possible = append(possible, val.Name)
				}

				ev, isEnum := value.(gqlparser.EnumValue)
				if !isEnum || def.EnumValue(string(ev)) == nil {
					rawValStr := fmt.Sprint(rawVal)

					addError(
						Message("Expected type %s, found %s.", valueType.String(), value),
						SuggestListUnquoted("Did you mean the enum value", rawValStr, possible),
					)
				}

			case gqlparser.Scalar:
				_, err := value.Value(nil)
				if err != nil {
					fmt.Println(err.Error())
					addError(Message("Expected type %s, found %s.", valueType.String(), value))
				}

				if !validateScalar(valueType, value) {
					addError(Message("Expected type %s, found %s.", valueType.String(), value))
				}
			}
		})
	})
}

func validateScalar(valueType gqlparser.Type, value gqlparser.Value) bool {
	switch value.(type) {
	case gqlparser.NullValue:
		return valueType.IsNullable()
	case gqlparser.Variable:
		return true
	case gqlparser.IntValue:
		return valueType.Name() == "Int" || valueType.Name() == "Float" || valueType.Name() == "ID"
	case gqlparser.FloatValue:
		return valueType.Name() == "Float"
	case gqlparser.StringValue, gqlparser.BlockValue:
		return valueType.Name() == "String" || valueType.Name() == "ID"
	case gqlparser.BooleanValue:
		return valueType.Name() == "Boolean"
	}
	return false
}
