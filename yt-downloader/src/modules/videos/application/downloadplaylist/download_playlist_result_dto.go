package downloadplaylist

import (
	"yt-downloader-go/src/modules/shared/application"
	"yt-downloader-go/src/modules/videos/domain/enums"
)

// DownloadPlaylistResultDto resumen de la playlist: lo descargado (resultado de cada video) y lo que fallo
// (video + motivo), para que el llamador pueda reintentar solo los fallidos.
type DownloadPlaylistResultDto struct {
	playlistId string
	title      string
	outputPath string
	total      int
	downloaded []map[string]any
	failed     []map[string]any
}

func DownloadPlaylistResultDtoFromPrimitives(primitives map[string]any) *DownloadPlaylistResultDto {
	total, _ := primitives[enums.VideosKeyEnum.TOTAL].(int)
	downloaded, _ := primitives[enums.VideosKeyEnum.DOWNLOADED].([]map[string]any)
	failed, _ := primitives[enums.VideosKeyEnum.FAILED].([]map[string]any)
	if downloaded == nil {
		downloaded = []map[string]any{}
	}
	if failed == nil {
		failed = []map[string]any{}
	}
	return &DownloadPlaylistResultDto{
		playlistId: application.GetTrimmedString(primitives, enums.VideosKeyEnum.PLAYLIST_ID),
		title:      application.GetTrimmedString(primitives, enums.VideosKeyEnum.TITLE),
		outputPath: application.GetTrimmedString(primitives, enums.VideosKeyEnum.OUTPUT_PATH),
		total:      total,
		downloaded: downloaded,
		failed:     failed,
	}
}

func (d *DownloadPlaylistResultDto) CountDownloaded() int { return len(d.downloaded) }
func (d *DownloadPlaylistResultDto) CountFailed() int     { return len(d.failed) }

func (d *DownloadPlaylistResultDto) ToArray() map[string]any {
	return map[string]any{
		enums.VideosKeyEnum.PLAYLIST_ID: d.playlistId,
		enums.VideosKeyEnum.TITLE:       d.title,
		enums.VideosKeyEnum.OUTPUT_PATH: d.outputPath,
		enums.VideosKeyEnum.TOTAL:       d.total,
		enums.VideosKeyEnum.DOWNLOADED:  d.downloaded,
		enums.VideosKeyEnum.FAILED:      d.failed,
	}
}
