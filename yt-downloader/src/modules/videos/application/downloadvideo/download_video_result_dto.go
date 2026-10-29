package downloadvideo

import (
	"yt-downloader-go/src/modules/shared/application"
	"yt-downloader-go/src/modules/videos/domain/entities"
	"yt-downloader-go/src/modules/videos/domain/enums"
)

type DownloadVideoResultDto struct {
	videoId  string
	title    string
	filePath string
	bytes    int64
	format   entities.VideoFormatPrimitives
}

func DownloadVideoResultDtoFromPrimitives(primitives map[string]any) *DownloadVideoResultDto {
	bytes, _ := primitives[enums.VideosKeyEnum.BYTES].(int64)
	format, _ := primitives[enums.VideosKeyEnum.FORMAT].(entities.VideoFormatPrimitives)
	return &DownloadVideoResultDto{
		videoId:  application.GetTrimmedString(primitives, enums.VideosKeyEnum.VIDEO_ID),
		title:    application.GetTrimmedString(primitives, enums.VideosKeyEnum.TITLE),
		filePath: application.GetTrimmedString(primitives, enums.VideosKeyEnum.FILE_PATH),
		bytes:    bytes,
		format:   format,
	}
}

func (d *DownloadVideoResultDto) GetFilePath() string {
	return d.filePath
}

func (d *DownloadVideoResultDto) ToArray() map[string]any {
	return map[string]any{
		enums.VideosKeyEnum.VIDEO_ID:  d.videoId,
		enums.VideosKeyEnum.TITLE:     d.title,
		enums.VideosKeyEnum.FILE_PATH: d.filePath,
		enums.VideosKeyEnum.BYTES:     d.bytes,
		enums.VideosKeyEnum.FORMAT:    d.format,
	}
}
