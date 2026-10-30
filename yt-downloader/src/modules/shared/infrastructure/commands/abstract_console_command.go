package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"yt-downloader-go/src/core/config"
	"yt-downloader-go/src/modules/shared/domain/enums"
	"yt-downloader-go/src/modules/shared/infrastructure/exceptions"
	"yt-downloader-go/src/modules/shared/infrastructure/traits"
)

const (
	EXIT_CODE_DOMAIN_ERROR         = 1
	EXIT_CODE_INFRASTRUCTURE_ERROR = 2
)

// exitFunc como termina un comando que falla. Por defecto os.Exit; el menu interactivo lo sustituye (SetExitFunc)
// para volver al menu en vez de cerrar la ventana.
var exitFunc = os.Exit

func SetExitFunc(fn func(code int)) {
	exitFunc = fn
}

// AbstractConsoleCommand base de todo comando: lee --clave=valor, imprime el ResultDto en JSON por stdout (sin
// color: se puede pasar a jq) y los errores en rojo por stderr con su codigo de salida.
// 0 ok · 1 error de dominio (lo que en HTTP es un 4xx) · 2 error de infraestructura (el 500, se loguea).
type AbstractConsoleCommand struct {
	traits.CliColorTrait
	traits.LogTrait
	// los argumentos que van tras el alias del comando
	input []string
}

func NewAbstractConsoleCommand(input []string) AbstractConsoleCommand {
	return AbstractConsoleCommand{LogTrait: traits.NewLogTrait(), input: input}
}

// GetOptionValue --url="https://..." -> "https://...". Si no se pasa, "" (el servicio decide si es obligatorio).
func (c *AbstractConsoleCommand) GetOptionValue(optionName string) string {
	optionPrefix := "--" + optionName + "="
	for _, argument := range c.input {
		if strings.HasPrefix(argument, optionPrefix) {
			return strings.Trim(argument[len(optionPrefix):], `"'`)
		}
	}
	return ""
}

// PrintResult JSON con 4 espacios sin escapar unicode ni barras.
func (c *AbstractConsoleCommand) PrintResult(result any) {
	var jsonBuffer bytes.Buffer
	jsonEncoder := json.NewEncoder(&jsonBuffer)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent("", "    ")
	_ = jsonEncoder.Encode(result)
	fmt.Print(jsonBuffer.String())
}

func (c *AbstractConsoleCommand) ExitWithDomainError(ex exceptions.InterfaceException) {
	c.printErrors(ex.GetCode(), []string{ex.Error()})
	exitFunc(EXIT_CODE_DOMAIN_ERROR)
}

// ExitWithUnexpectedError como el 500 de una API: se loguea con el metodo que fallo; quien tiene la consola ve el
// mensaje real (LogDebugCli).
func (c *AbstractConsoleCommand) ExitWithUnexpectedError(err error, logTitle string) {
	c.LogException(err, logTitle)
	c.LogDebugCli(err.Error(), logTitle)
	c.printErrors(
		enums.ExceptionEnum.CODE_INTERNAL_SERVER_ERROR,
		[]string{config.Tr(enums.SharedMessageEnum.AN_UNEXPECTED_ERROR_HAS_OCCURRED)},
	)
	exitFunc(EXIT_CODE_INFRASTRUCTURE_ERROR)
}

func (c *AbstractConsoleCommand) printErrors(code int, errors []string) {
	fmt.Fprintln(os.Stderr, c.GetRed(fmt.Sprintf("Error (%d):", code)))
	for _, errorMessage := range errors {
		fmt.Fprintln(os.Stderr, c.GetRed("  - "+errorMessage))
	}
}
