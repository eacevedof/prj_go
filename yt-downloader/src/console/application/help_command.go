package application

import (
	"fmt"
	"strings"

	"yt-downloader-go/src/modules/shared/infrastructure/commands"
	"yt-downloader-go/src/modules/shared/infrastructure/repositories/configuration"
)

// HelpCommand ./console help (o sin alias): en verde el comando, en amarillo la llamada que se copia.
type HelpCommand struct {
	commands.AbstractConsoleCommand
	commandConfigs []CommandConfig
}

func HelpCommandGetInstance(input []string, commandConfigs []CommandConfig) *HelpCommand {
	return &HelpCommand{
		AbstractConsoleCommand: commands.NewAbstractConsoleCommand(input),
		commandConfigs:         commandConfigs,
	}
}

func (c *HelpCommand) Run() {
	var lines []string
	for i, commandConfig := range c.commandConfigs {
		info := ""
		if commandConfig.Info != "" {
			info = "\n\t" + commandConfig.Info
		}
		call := strings.TrimSpace("./console " + commandConfig.Alias + " " + commandConfig.Example)
		lines = append(lines, fmt.Sprintf("\n(%d) ", i+1)+c.GetGreen(commandConfig.Alias)+":"+info+"\n\t"+c.GetYellow(call))
	}
	environment := configuration.GetInstance()
	c.EchoBlue(fmt.Sprintf("%s v%s (%s) - help", environment.GetAppName(), environment.GetAppVersion(), environment.GetAppVersionUpdate()))
	fmt.Println(strings.Join(lines, ""))
}
