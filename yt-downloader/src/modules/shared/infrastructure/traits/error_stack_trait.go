package traits

import "strings"

// ErrorStackTrait da formato a la pila de errores de un validador para meterla en el mensaje de la excepcion.
type ErrorStackTrait struct{}

// GetErrorStackMessage "Fields validation errors: url: field is required; format: invalid format"
func (t ErrorStackTrait) GetErrorStackMessage(title string, errorStack []string) string {
	return title + ": " + strings.Join(errorStack, "; ")
}
