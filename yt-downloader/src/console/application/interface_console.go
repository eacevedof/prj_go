package application

// InterfaceConsole lo que tiene que hacer todo comando de consola.
type InterfaceConsole interface {
	Run()
}

// CommandConfig una entrada de console/commands.go: alias => el command, que hace y los argumentos de un ejemplo.
type CommandConfig struct {
	Alias   string
	Service func(input []string) InterfaceConsole
	Info    string
	Example string
}
