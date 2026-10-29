package downloadvideo

import (
	"yt-downloader-go/src/modules/shared/application"
	"yt-downloader-go/src/modules/videos/domain/enums"
)

type DownloadVideoDto struct {
	url     string
	format  string
	quality string
	// carpeta destino; vacio = APP_DOWNLOADS_PATH
	output string
	// nombre de fichero sin extension; vacio = titulo del video
	fileName string
}

func DownloadVideoDtoFromPrimitives(primitives map[string]any) *DownloadVideoDto {
	format := application.GetTrimmedString(primitives, enums.VideosKeyEnum.FORMAT)
	if format == "" {
		format = enums.VideoFormatEnum.VIDEO
	}
	quality := application.GetTrimmedString(primitives, enums.VideosKeyEnum.QUALITY)
	if quality == "" {
		quality = enums.VideoQualityEnum.BEST
	}
	return &DownloadVideoDto{
		url:      application.GetTrimmedString(primitives, enums.VideosKeyEnum.URL),
		format:   format,
		quality:  quality,
		output:   application.GetTrimmedString(primitives, enums.VideosKeyEnum.OUTPUT),
		fileName: application.GetTrimmedString(primitives, enums.VideosKeyEnum.FILENAME),
	}
}

func (d *DownloadVideoDto) GetUrl() string      { return d.url }
func (d *DownloadVideoDto) GetFormat() string   { return d.format }
func (d *DownloadVideoDto) GetQuality() string  { return d.quality }
func (d *DownloadVideoDto) GetOutput() string   { return d.output }
func (d *DownloadVideoDto) GetFileName() string { return d.fileName }
