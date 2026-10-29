package application

import (
	"fmt"
	"maps"
	"strconv"
	"strings"
)

// Los dto se construyen desde un map de primitivos para que el command pueda juntar todos los campos de la
// peticion en un solo map. Casts al estilo PHP sobre un valor que puede faltar, ser nil o venir de otro tipo.

// GetTrimmedString trim((string) ($primitives[$key] ?? ""))
func GetTrimmedString(primitives map[string]any, key string) string {
	return strings.TrimSpace(ToString(primitives[key]))
}

// MergePrimitives [...$a, ...$b]: con claves de texto, la ultima gana.
func MergePrimitives(primitivesList ...map[string]any) map[string]any {
	merged := map[string]any{}
	for _, primitives := range primitivesList {
		maps.Copy(merged, primitives)
	}
	return merged
}

// ToString (string) $value de PHP: true "1", false/nil "", 1.0 "1".
func ToString(value any) string {
	switch typedValue := value.(type) {
	case nil:
		return ""
	case string:
		return typedValue
	case bool:
		if typedValue {
			return "1"
		}
		return ""
	case int:
		return strconv.Itoa(typedValue)
	case int64:
		return strconv.FormatInt(typedValue, 10)
	case float64:
		return strconv.FormatFloat(typedValue, 'f', -1, 64)
	default:
		return fmt.Sprint(typedValue)
	}
}
