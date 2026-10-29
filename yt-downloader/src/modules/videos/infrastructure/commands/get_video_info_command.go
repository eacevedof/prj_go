package commands

import (
	"errors"

	sharedcommands "yt-downloader-go/src/modules/shared/infrastructure/commands"
	"yt-downloader-go/src/modules/shared/infrastructure/exceptions"
	"yt-downloader-go/src/modules/videos/application/getvideoinfo"
	"yt-downloader-go/src/modules/videos/domain/enums"
	videosexceptions "yt-downloader-go/src/modules/videos/domain/exceptions"
)

// GetVideoInfoCommand ./console video-info --url=<url o id>
type GetVideoInfoCommand struct {
	sharedcommands.AbstractConsoleCommand
}

func GetVideoInfoCommandGetInstance(input []string) *GetVideoInfoCommand {
	return &GetVideoInfoCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *GetVideoInfoCommand) Run() {
	getVideoInfoResultDto, err := getvideoinfo.GetVideoInfoServiceGetInstance().Invoke(
		getvideoinfo.GetVideoInfoDtoFromPrimitives(map[string]any{
			enums.VideosKeyEnum.URL: c.GetOptionValue(enums.VideosKeyEnum.URL),
		}),
	)
	if err == nil {
		c.PrintResult(getVideoInfoResultDto.ToArray())
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
	c.ExitWithUnexpectedError(err, "GetVideoInfoCommand.Run")
}
