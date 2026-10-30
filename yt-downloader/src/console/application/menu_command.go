package application

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"yt-downloader-go/src/modules/shared/infrastructure/commands"
	"yt-downloader-go/src/modules/shared/infrastructure/repositories/configuration"
	"yt-downloader-go/src/modules/videos/application/getvideoinfo"
	videosenums "yt-downloader-go/src/modules/videos/domain/enums"
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

// qualityOptions para playlists (no se puede listar video a video): valor de --quality y su explicacion.
var qualityOptions = []struct {
	value string
	label string
}{
	{value: "best", label: "best (the highest each video has; needs ffmpeg above 360p)"},
	{value: "hd1080", label: "hd1080 (up to 1080p)"},
	{value: "hd720", label: "hd720 (up to 720p)"},
	{value: "medium", label: "medium (360p, smaller files, no ffmpeg needed)"},
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
	// en el menu el resultado se pinta como tarjeta legible, no como JSON
	commands.SetHumanOutput(true)
	defer commands.SetHumanOutput(false)

	for {
		c.printMenu()
		choice := c.ask(fmt.Sprintf("select an option [1-%d]", len(menuOptions)+1), "")
		if choice == fmt.Sprint(len(menuOptions)+1) || choice == "" {
			return
		}
		optionIndex := -1
		if _, err := fmt.Sscanf(choice, "%d", &optionIndex); err != nil || optionIndex < 1 || optionIndex > len(menuOptions) {
			c.EchoRed("invalid option")
			continue
		}
		c.runOption(menuOptions[optionIndex-1])
		fmt.Println()
		c.ask("press enter to go back to the menu", "")
	}
}

func (c *MenuCommand) printMenu() {
	environment := configuration.GetInstance()
	fmt.Println()
	c.EchoBlue(fmt.Sprintf("%s v%s (%s)", environment.GetAppName(), environment.GetAppVersion(), environment.GetAppVersionUpdate()))
	fmt.Println(c.GetOrange("downloads folder: ") + c.GetYellow(environment.GetDownloadsPath()))
	fmt.Println()
	c.EchoOrange("what do you want to do?")
	for i, option := range menuOptions {
		fmt.Printf("  %s) %s\n", c.GetGreen(fmt.Sprint(i+1)), option.label)
	}
	fmt.Printf("  %s) Exit\n", c.GetGreen(fmt.Sprint(len(menuOptions)+1)))
	fmt.Println()
	fmt.Println("  (*) ffmpeg: in PowerShell run `winget install Gyan.FFmpeg`, or download it from " + ffmpegDownloadUrl + " and add its bin folder to the PATH")
	fmt.Println()
}

func (c *MenuCommand) runOption(option menuOption) {
	url := ""
	for url == "" {
		url = c.ask("paste the youtube url (or the video/playlist id)", "")
	}
	input := []string{"--url=" + url}
	if option.format != "" {
		input = append(input, "--format="+option.format)
	}
	if option.asksQuality {
		quality := ""
		if option.alias == "video-download" {
			quality = c.askVideoResolution(url)
		}
		if quality == "" {
			quality = c.askQuality()
		}
		input = append(input, "--quality="+quality)
	}
	if option.asksFilename {
		if filename := c.ask("file name (enter = the video title)", ""); filename != "" {
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

// askVideoResolution pide la info del video y lista sus resoluciones reales (1080p, 720p, 360p...) para elegir una;
// "" si la info falla (se cae a la lista generica) o el usuario no elige.
func (c *MenuCommand) askVideoResolution(url string) string {
	fmt.Println()
	c.EchoOrange("reading video info...")
	resultDto, err := getvideoinfo.GetVideoInfoServiceGetInstance().Invoke(
		getvideoinfo.GetVideoInfoDtoFromPrimitives(map[string]any{videosenums.VideosKeyEnum.URL: url}),
	)
	if err != nil {
		c.EchoRed(err.Error())
		return ""
	}
	heights := resultDto.GetAvailableHeights()
	if len(heights) == 0 {
		return ""
	}
	fmt.Println(c.GetOrange("title: ") + resultDto.GetTitle())
	c.EchoOrange("resolutions found (above 360p ffmpeg is needed to merge video and audio):")
	for i, height := range heights {
		fmt.Printf("  %s) %dp\n", c.GetGreen(fmt.Sprint(i+1)), height)
	}
	for {
		choice := c.ask("resolution", "1")
		index := -1
		if _, err := fmt.Sscanf(choice, "%d", &index); err == nil && index >= 1 && index <= len(heights) {
			return fmt.Sprintf("%dp", heights[index-1])
		}
		c.EchoRed("invalid option")
	}
}

func (c *MenuCommand) askQuality() string {
	fmt.Println()
	c.EchoOrange("quality?")
	for i, quality := range qualityOptions {
		fmt.Printf("  %s) %s\n", c.GetGreen(fmt.Sprint(i+1)), quality.label)
	}
	for {
		choice := c.ask("quality", "1")
		for i, quality := range qualityOptions {
			if choice == fmt.Sprint(i+1) {
				return quality.value
			}
		}
		c.EchoRed("invalid option")
	}
}

// ask pregunta (en naranja, en minusculas) por stdout y devuelve la linea sin espacios; Enter vacio = defaultValue. EOF (stdin
// cerrado) termina el programa.
func (c *MenuCommand) ask(question string, defaultValue string) string {
	if defaultValue != "" {
		fmt.Print(c.GetOrange(fmt.Sprintf("%s [%s]: ", question, defaultValue)))
	} else {
		fmt.Print(c.GetOrange(question + ": "))
	}
	line, err := c.reader.ReadString('\n')
	if err != nil && line == "" {
		// stdin cerrado (consola cerrada o entrada agotada): no hay nada mas que preguntar
		fmt.Println()
		os.Exit(0)
	}
	line = strings.TrimSpace(strings.Trim(strings.TrimSpace(line), `"'`))
	if line == "" {
		return defaultValue
	}
	return line
}
