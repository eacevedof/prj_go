package application

import (
	"fmt"

	"yt-downloader-go/src/modules/shared/infrastructure/commands"
	"yt-downloader-go/src/modules/shared/infrastructure/repositories/configuration"
)

// VersionCommand ./console version -> "eaf-yt-dl v0.2.0 (2026-10-30)". La version es APP_VERSION (unica fuente) y da
// nombre al instalable de make dist.
type VersionCommand struct {
	commands.AbstractConsoleCommand
}

func VersionCommandGetInstance(input []string) *VersionCommand {
	return &VersionCommand{AbstractConsoleCommand: commands.NewAbstractConsoleCommand(input)}
}

func (c *VersionCommand) Run() {
	environment := configuration.GetInstance()
	fmt.Printf("%s v%s (%s)\n", environment.GetAppName(), environment.GetAppVersion(), environment.GetAppVersionUpdate())
}
