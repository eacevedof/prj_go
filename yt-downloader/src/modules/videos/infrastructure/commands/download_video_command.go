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
	lastPercent int
}

func DownloadVideoCommandGetInstance(input []string) *DownloadVideoCommand {
	return &DownloadVideoCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input), lastPercent: -1}
}

func (c *DownloadVideoCommand) Run() {
	downloadVideoResultDto, err := downloadvideo.DownloadVideoServiceGetInstance(c.printProgress).Invoke(
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

// printProgress una sola linea que se sobrescribe: "  45% 12.3 MB / 27.1 MB"
func (c *DownloadVideoCommand) printProgress(read int64, total int64) {
	if total <= 0 {
		fmt.Fprintf(os.Stderr, "\r  %s", getHumanBytes(read))
		return
	}
	percent := int(read * 100 / total)
	if percent == c.lastPercent {
		return
	}
	c.lastPercent = percent
	fmt.Fprintf(os.Stderr, "\r  %s", c.GetGreen(fmt.Sprintf("%3d%% %s / %s", percent, getHumanBytes(read), getHumanBytes(total))))
}

func getHumanBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	value := float64(bytes)
	units := []string{"KB", "MB", "GB"}
	i := -1
	for value >= unit && i < len(units)-1 {
		value /= unit
		i++
	}
	return fmt.Sprintf("%.1f %s", value, units[i])
}
