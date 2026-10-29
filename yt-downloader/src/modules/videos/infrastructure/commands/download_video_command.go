package commands

import (
	"errors"
	"fmt"
	"os"

	sharedcommands "yt-downloader-go/src/modules/shared/infrastructure/commands"
	"yt-downloader-go/src/modules/shared/infrastructure/exceptions"
	"yt-downloader-go/src/modules/videos/application/downloadvideo"
	"yt-downloader-go/src/modules/videos/domain/enums"
	videosexceptions "yt-downloader-go/src/modules/videos/domain/exceptions"
)

// DownloadVideoCommand ./console video-download --url=<url> --format=video|audio --quality=best|hd720|medium
// --output=<carpeta> --filename=<nombre sin extension>
// El progreso va por stderr (stdout queda limpio para el JSON del resultado).
type DownloadVideoCommand struct {
	sharedcommands.AbstractConsoleCommand
	progress *progressPrinter
}

func DownloadVideoCommandGetInstance(input []string) *DownloadVideoCommand {
	return &DownloadVideoCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input), progress: newProgressPrinter()}
}

func (c *DownloadVideoCommand) Run() {
	downloadVideoResultDto, err := downloadvideo.DownloadVideoServiceGetInstance(c.progress.Print).Invoke(
		downloadvideo.DownloadVideoDtoFromPrimitives(map[string]any{
			enums.VideosKeyEnum.URL:      c.GetOptionValue(enums.VideosKeyEnum.URL),
			enums.VideosKeyEnum.FORMAT:   c.GetOptionValue(enums.VideosKeyEnum.FORMAT),
			enums.VideosKeyEnum.QUALITY:  c.GetOptionValue(enums.VideosKeyEnum.QUALITY),
			enums.VideosKeyEnum.OUTPUT:   c.GetOptionValue(enums.VideosKeyEnum.OUTPUT),
			enums.VideosKeyEnum.FILENAME: c.GetOptionValue(enums.VideosKeyEnum.FILENAME),
		}),
	)
	if err == nil {
		fmt.Fprintln(os.Stderr)
		c.PrintResult(downloadVideoResultDto.ToArray())
		return
	}

	var fieldsException *exceptions.FieldsException
	if errors.As(err, &fieldsException) {
		c.ExitWithDomainError(fieldsException)
	}
	var videosException *videosexceptions.VideosException
	if errors.As(err, &videosException) {
		c.ExitWithDomainError(videosException)
	}
	c.ExitWithUnexpectedError(err, "DownloadVideoCommand.Run")
}
