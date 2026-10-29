package commands

import (
	"errors"
	"fmt"
	"os"

	sharedcommands "yt-downloader-go/src/modules/shared/infrastructure/commands"
	"yt-downloader-go/src/modules/shared/infrastructure/exceptions"
	"yt-downloader-go/src/modules/videos/application/downloadplaylist"
	"yt-downloader-go/src/modules/videos/domain/enums"
	videosexceptions "yt-downloader-go/src/modules/videos/domain/exceptions"
)

// DownloadPlaylistCommand ./console playlist-download --url=<url o id de playlist> --format=video|audio
// --quality=best|hd720|medium --output=<carpeta base>
// Cada video se anuncia "[i/n] titulo" y su progreso va por stderr; el JSON resumen por stdout. Un video que falla
// se marca y se sigue; exit 1 solo si no se descargo ninguno.
type DownloadPlaylistCommand struct {
	sharedcommands.AbstractConsoleCommand
	progress *progressPrinter
}

func DownloadPlaylistCommandGetInstance(input []string) *DownloadPlaylistCommand {
	return &DownloadPlaylistCommand{
		AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input),
		progress:               newProgressPrinter(),
	}
}

func (c *DownloadPlaylistCommand) Run() {
	resultDto, err := downloadplaylist.DownloadPlaylistServiceGetInstance(c.progress.Print, c.printVideoStart, c.printVideoEnd).Invoke(
		downloadplaylist.DownloadPlaylistDtoFromPrimitives(map[string]any{
			enums.VideosKeyEnum.URL:     c.GetOptionValue(enums.VideosKeyEnum.URL),
			enums.VideosKeyEnum.FORMAT:  c.GetOptionValue(enums.VideosKeyEnum.FORMAT),
			enums.VideosKeyEnum.QUALITY: c.GetOptionValue(enums.VideosKeyEnum.QUALITY),
			enums.VideosKeyEnum.OUTPUT:  c.GetOptionValue(enums.VideosKeyEnum.OUTPUT),
		}),
	)
	if err == nil {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, c.GetBlue(fmt.Sprintf("%d descargados, %d fallidos", resultDto.CountDownloaded(), resultDto.CountFailed())))
		c.PrintResult(resultDto.ToArray())
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
	c.ExitWithUnexpectedError(err, "DownloadPlaylistCommand.Run")
}

func (c *DownloadPlaylistCommand) printVideoStart(index int, total int, title string) {
	c.progress.Reset()
	fmt.Fprintf(os.Stderr, "\n%s %s\n", c.GetYellow(fmt.Sprintf("[%d/%d]", index, total)), title)
}

func (c *DownloadPlaylistCommand) printVideoEnd(index int, total int, title string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "\r  %s %s\n", c.GetRed("ERROR"), err.Error())
	}
}
