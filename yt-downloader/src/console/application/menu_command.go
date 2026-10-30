package application

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"yt-downloader-go/src/modules/shared/infrastructure/commands"
	"yt-downloader-go/src/modules/shared/infrastructure/repositories/configuration"
)

// MenuCommand ./console menu (o sin argumentos, o doble clic en el .exe): pregunta que hacer y los datos minimos
// (url, calidad) y lanza el comando de consola correspondiente. Pensado para quien no monta comandos a mano.
// Los textos van en castellano; los valores de los parametros (video|audio, best|hd720|medium) siempre en ingles.
type MenuCommand struct {
	commands.AbstractConsoleCommand
	commandConfigs []CommandConfig
	reader         *bufio.Reader
}

// menuOption una entrada del menu: que alias lanza y con que --format fijo.
type menuOption struct {
	label        string
	alias        string
	format       string
	asksQuality  bool
	asksFilename bool
}

// commandExit se lanza (panic) en vez de os.Exit cuando un comando falla dentro del menu; Run() lo recupera.
type commandExit struct{ code int }

var menuOptions = []menuOption{
	{label: "Descargar un video (mp4)", alias: "video-download", format: "video", asksQuality: true, asksFilename: true},
	{label: "Ver la informacion de un video", alias: "video-info"},
	{label: "Descargar una playlist en video (mp4)", alias: "playlist-download", format: "video", asksQuality: true},
	{label: "Descargar solo el audio de un video (m4a, necesita ffmpeg)", alias: "video-download", format: "audio", asksFilename: true},
	{label: "Descargar una playlist en audio (m4a, necesita ffmpeg)", alias: "playlist-download", format: "audio"},
}

var qualityOptions = []string{"best", "hd720", "medium"}

func MenuCommandGetInstance(input []string, commandConfigs []CommandConfig) *MenuCommand {
	return &MenuCommand{
		AbstractConsoleCommand: commands.NewAbstractConsoleCommand(input),
		commandConfigs:         commandConfigs,
		reader:                 bufio.NewReader(os.Stdin),
	}
}

func (c *MenuCommand) Run() {
	commands.SetExitFunc(func(code int) { panic(commandExit{code: code}) })
	defer commands.SetExitFunc(os.Exit)

	for {
		c.printMenu()
		choice := c.ask("Opcion", "0")
		if choice == "0" || choice == "" {
			return
		}
		optionIndex := -1
		if _, err := fmt.Sscanf(choice, "%d", &optionIndex); err != nil || optionIndex < 1 || optionIndex > len(menuOptions) {
			c.EchoRed("Opcion no valida")
			continue
		}
		c.runOption(menuOptions[optionIndex-1])
		fmt.Println()
		c.ask("Pulsa Enter para volver al menu", "")
	}
}

func (c *MenuCommand) printMenu() {
	environment := configuration.GetInstance()
	fmt.Println()
	c.EchoBlue(fmt.Sprintf("%s v%s", environment.GetAppName(), environment.GetAppVersion()))
	fmt.Println("Las descargas van a: " + c.GetYellow(environment.GetDownloadsPath()))
	fmt.Println()
	for i, option := range menuOptions {
		fmt.Printf("  %s) %s\n", c.GetGreen(fmt.Sprint(i+1)), option.label)
	}
	fmt.Printf("  %s) Salir\n", c.GetGreen("0"))
	fmt.Println()
}

func (c *MenuCommand) runOption(option menuOption) {
	url := ""
	for url == "" {
		url = c.ask("Pega la URL de YouTube (o el id)", "")
	}
	input := []string{"--url=" + url}
	if option.format != "" {
		input = append(input, "--format="+option.format)
	}
	if option.asksQuality {
		input = append(input, "--quality="+c.askQuality())
	}
	if option.asksFilename {
		if filename := c.ask("Nombre del fichero (Enter = el titulo del video)", ""); filename != "" {
			input = append(input, "--filename="+filename)
		}
	}
	fmt.Println()

	defer func() {
		if recovered := recover(); recovered != nil {
			if _, isExit := recovered.(commandExit); !isExit {
				panic(recovered)
			}
		}
	}()
	for _, commandConfig := range c.commandConfigs {
		if commandConfig.Alias == option.alias {
			commandConfig.Service(input).Run()
			return
		}
	}
	c.EchoRed("no command found for alias " + option.alias)
}

func (c *MenuCommand) askQuality() string {
	fmt.Println("Calidad: 1) best (la mejor disponible)  2) hd720  3) medium (360p, ocupa menos)")
	for {
		choice := c.ask("Calidad", "1")
		switch choice {
		case "1", "":
			return qualityOptions[0]
		case "2":
			return qualityOptions[1]
		case "3":
			return qualityOptions[2]
		}
		c.EchoRed("Opcion no valida")
	}
}

// ask pregunta por stdout y devuelve la linea sin espacios; Enter vacio = defaultValue. EOF (consola cerrada) = "0".
func (c *MenuCommand) ask(question string, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", question, defaultValue)
	} else {
		fmt.Printf("%s: ", question)
	}
	line, err := c.reader.ReadString('\n')
	if err != nil && line == "" {
		return "0"
	}
	line = strings.TrimSpace(strings.Trim(strings.TrimSpace(line), `"'`))
	if line == "" {
		return defaultValue
	}
	return line
}
