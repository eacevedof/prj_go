package boot

var colorsEnabled = true

// EnableConsoleColors se llama una vez al arrancar la consola: activa los colores ANSI en la terminal si hace falta
// y, si no es posible, los desactiva para que no salgan los codigos en crudo.
func EnableConsoleColors() {
	colorsEnabled = enableVirtualTerminal()
}

func AreColorsEnabled() bool {
	return colorsEnabled
}
