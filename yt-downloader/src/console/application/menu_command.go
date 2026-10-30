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
// Todo el dialogo va en ingles; lo que pregunta la app se pinta en naranja.
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
	{label: "Download a video (mp4)", alias: "video-download", format: "video", asksQuality: true, asksFilename: true},
	{label: "Show video info", alias: "video-info"},
	{label: "Download a playlist as video (mp4)", alias: "playlist-download", format: "video", asksQuality: true},
	{label: "Download only the audio of a video (m4a) - ffmpeg must be installed first (*)", alias: "video-download", format: "audio", asksFilename: true},
	{label: "Download a playlist as audio (m4a) - ffmpeg must be installed first (*)", alias: "playlist-download", format: "audio"},
}

// qualityOptions valor del parametro --quality y su explicacion. Solo hay progresivo (ver README §3): 720p como mucho.
var qualityOptions = []struct {
	value string
	label string
}{
	{value: "best", label: "best (the highest available: 720p or 360p, depends on the video)"},
	{value: "hd720", label: "hd720 (720p if the video has it)"},
	{value: "medium", label: "medium (360p, smaller file)"},
}

const ffmpegDownloadUrl = "https://www.gyan.dev/ffmpeg/builds/"

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
		choice := c.ask("Option", "0")
		if choice == "0" || choice == "" {
			return
		}
		optionIndex := -1
		if _, err := fmt.Sscanf(choice, "%d", &optionIndex); err != nil || optionIndex < 1 || optionIndex > len(menuOptions) {
			c.EchoRed("Invalid option")
			continue
		}
		c.runOption(menuOptions[optionIndex-1])
		fmt.Println()
		c.ask("Press Enter to go back to the menu", "")
	}
}

func (c *MenuCommand) printMenu() {
	environment := configuration.GetInstance()
	fmt.Println()
	c.EchoBlue(fmt.Sprintf("%s v%s (%s)", environment.GetAppName(), environment.GetAppVersion(), environment.GetAppVersionUpdate()))
	fmt.Println(c.GetOrange("Downloads folder: ") + c.GetYellow(environment.GetDownloadsPath()))
	fmt.Println()
	c.EchoOrange("What do you want to do?")
	for i, option := range menuOptions {
		fmt.Printf("  %s) %s\n", c.GetGreen(fmt.Sprint(i+1)), option.label)
	}
	fmt.Printf("  %s) Exit\n", c.GetGreen("0"))
	fmt.Println()
	fmt.Println("  (*) ffmpeg: in PowerShell run `winget install Gyan.FFmpeg`, or download it from " + ffmpegDownloadUrl + " and add its bin folder to the PATH")
	fmt.Println()
}

func (c *MenuCommand) runOption(option menuOption) {
	url := ""
	for url == "" {
		url = c.ask("Paste the YouTube URL (or the video/playlist id)", "")
	}
	input := []string{"--url=" + url}
	if option.format != "" {
		input = append(input, "--format="+option.format)
	}
	if option.asksQuality {
		input = append(input, "--quality="+c.askQuality())
	}
	if option.asksFilename {
		if filename := c.ask("File name (Enter = the video title)", ""); filename != "" {
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
	fmt.Println()
	c.EchoOrange("Quality?")
	for i, quality := range qualityOptions {
		fmt.Printf("  %s) %s\n", c.GetGreen(fmt.Sprint(i+1)), quality.label)
	}
	for {
		choice := c.ask("Quality", "1")
		for i, quality := range qualityOptions {
			if choice == fmt.Sprint(i+1) {
				return quality.value
			}
		}
		c.EchoRed("Invalid option")
	}
}

// ask pregunta (en naranja) por stdout y devuelve la linea sin espacios; Enter vacio = defaultValue. EOF (consola
// cerrada) = "0".
func (c *MenuCommand) ask(question string, defaultValue string) string {
	if defaultValue != "" {
		fmt.Print(c.GetOrange(fmt.Sprintf("%s [%s]: ", question, defaultValue)))
	} else {
		fmt.Print(c.GetOrange(question + ": "))
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
