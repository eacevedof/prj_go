package downloadplaylist

import (
	"yt-downloader-go/src/modules/shared/application"
	"yt-downloader-go/src/modules/videos/domain/enums"
)

type DownloadPlaylistDto struct {
	url     string
	format  string
	quality string
	// carpeta base; vacio = APP_DOWNLOADS_PATH. Dentro se crea una subcarpeta con el titulo de la playlist.
	output string
}

func DownloadPlaylistDtoFromPrimitives(primitives map[string]any) *DownloadPlaylistDto {
	format := application.GetTrimmedString(primitives, enums.VideosKeyEnum.FORMAT)
	if format == "" {
		format = enums.VideoFormatEnum.VIDEO
	}
	quality := application.GetTrimmedString(primitives, enums.VideosKeyEnum.QUALITY)
	if quality == "" {
		quality = enums.VideoQualityEnum.BEST
	}
	return &DownloadPlaylistDto{
		url:     application.GetTrimmedString(primitives, enums.VideosKeyEnum.URL),
		format:  format,
		quality: quality,
		output:  application.GetTrimmedString(primitives, enums.VideosKeyEnum.OUTPUT),
	}
}

func (d *DownloadPlaylistDto) GetUrl() string     { return d.url }
func (d *DownloadPlaylistDto) GetFormat() string  { return d.format }
func (d *DownloadPlaylistDto) GetQuality() string { return d.quality }
func (d *DownloadPlaylistDto) GetOutput() string  { return d.output }
