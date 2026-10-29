package downloadplaylist

import (
	"path/filepath"

	"yt-downloader-go/src/core/config"
	sharedenums "yt-downloader-go/src/modules/shared/domain/enums"
	"yt-downloader-go/src/modules/shared/infrastructure/components/ffmpeg"
	"yt-downloader-go/src/modules/shared/infrastructure/exceptions"
	"yt-downloader-go/src/modules/shared/infrastructure/repositories/configuration"
	"yt-downloader-go/src/modules/shared/infrastructure/traits"
	"yt-downloader-go/src/modules/videos/application/downloadvideo"
	"yt-downloader-go/src/modules/videos/domain/entities"
	"yt-downloader-go/src/modules/videos/domain/enums"
	videosexceptions "yt-downloader-go/src/modules/videos/domain/exceptions"
	"yt-downloader-go/src/modules/videos/domain/services"
	"yt-downloader-go/src/modules/videos/infrastructure/repositories"
)

type ProgressFunc = repositories.ProgressFunc

// VideoStartFunc el command la pasa para anunciar cada video antes de bajarlo: (posicion 1..n, total, titulo).
type VideoStartFunc func(index int, total int, title string)

// VideoEndFunc tras cada video: err nil = descargado; si no, el motivo. La descarga de la playlist no se detiene.
type VideoEndFunc func(index int, total int, title string, err error)

// DownloadPlaylistService lista la playlist y delega cada video en DownloadVideoService, secuencialmente, en la
// subcarpeta <output>/<titulo de la playlist>. Un video que falla (privado, borrado, sin formato) se anota y se
// sigue con el siguiente; solo es error de dominio que no se pueda descargar ninguno.
type DownloadPlaylistService struct {
	traits.ErrorStackTrait
	videoFieldsValidator          *services.VideoFieldsValidator
	fileNameSanitizer             *services.FileNameSanitizer
	videosReaderYoutubeRepository *repositories.VideosReaderYoutubeRepository
	audioExtractor                *ffmpeg.AudioExtractor
	downloadVideoService          *downloadvideo.DownloadVideoService

	downloadPlaylistDto *DownloadPlaylistDto
	onVideoStart        VideoStartFunc
	onVideoEnd          VideoEndFunc
}

func DownloadPlaylistServiceGetInstance(onProgress ProgressFunc, onVideoStart VideoStartFunc, onVideoEnd VideoEndFunc) *DownloadPlaylistService {
	return &DownloadPlaylistService{
		videoFieldsValidator:          services.VideoFieldsValidatorGetInstance(),
		fileNameSanitizer:             services.FileNameSanitizerGetInstance(),
		videosReaderYoutubeRepository: repositories.VideosReaderYoutubeRepositoryGetInstance(),
		audioExtractor:                ffmpeg.GetInstance(),
		downloadVideoService:          downloadvideo.DownloadVideoServiceGetInstance(onProgress),
		onVideoStart:                  onVideoStart,
		onVideoEnd:                    onVideoEnd,
	}
}

// Invoke error: *FieldsException (400), *VideosException (404 playlist, 409 vacia/nada descargado) o inesperado.
func (s *DownloadPlaylistService) Invoke(downloadPlaylistDto *DownloadPlaylistDto) (*DownloadPlaylistResultDto, error) {
	s.downloadPlaylistDto = downloadPlaylistDto
	if err := s.failIfWrongInput(); err != nil {
		return nil, err
	}

	playlistId := s.videoFieldsValidator.GetPlaylistIdFromUrl(downloadPlaylistDto.GetUrl())
	playlistEntity, err := s.videosReaderYoutubeRepository.GetPlaylistById(playlistId)
	if err != nil {
		return nil, err
	}
	if playlistEntity.IsEmpty() {
		return nil, videosexceptions.ConflictCustom(config.Tr(enums.VideosMessageEnum.PLAYLIST_X_IS_EMPTY, playlistId))
	}

	outputPath := s.getOutputFolder(playlistEntity)
	downloaded, failed := s.downloadEntries(playlistEntity, outputPath)
	if len(downloaded) == 0 {
		return nil, videosexceptions.ConflictCustom(config.Tr(
			enums.VideosMessageEnum.PLAYLIST_X_NOTHING_DOWNLOADED, playlistId, playlistEntity.CountEntries(),
		))
	}

	return DownloadPlaylistResultDtoFromPrimitives(map[string]any{
		enums.VideosKeyEnum.PLAYLIST_ID: playlistEntity.GetPlaylistId(),
		enums.VideosKeyEnum.TITLE:       playlistEntity.GetTitle(),
		enums.VideosKeyEnum.OUTPUT_PATH: outputPath,
		enums.VideosKeyEnum.TOTAL:       playlistEntity.CountEntries(),
		enums.VideosKeyEnum.DOWNLOADED:  downloaded,
		enums.VideosKeyEnum.FAILED:      failed,
	}), nil
}

func (s *DownloadPlaylistService) failIfWrongInput() error {
	errors := s.videoFieldsValidator.GetPlaylistDownloadErrors(
		s.downloadPlaylistDto.GetUrl(), s.downloadPlaylistDto.GetFormat(), s.downloadPlaylistDto.GetQuality(),
	)
	if len(errors) > 0 {
		return exceptions.FieldsExceptionBadRequestCustom(
			s.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
		)
	}
	if s.downloadPlaylistDto.GetFormat() == enums.VideoFormatEnum.AUDIO && !s.audioExtractor.IsAvailable() {
		return videosexceptions.ConflictCustom(config.Tr(enums.VideosMessageEnum.FFMPEG_NOT_AVAILABLE))
	}
	return nil
}

// downloadEntries uno a uno; los fallos no cortan el bucle.
func (s *DownloadPlaylistService) downloadEntries(playlistEntity *entities.PlaylistEntity, outputPath string) ([]map[string]any, []map[string]any) {
	total := playlistEntity.CountEntries()
	downloaded := []map[string]any{}
	failed := []map[string]any{}

	for i, entry := range playlistEntity.GetEntries() {
		index := i + 1
		if s.onVideoStart != nil {
			s.onVideoStart(index, total, entry.GetTitle())
		}

		resultDto, err := s.downloadVideoService.Invoke(downloadvideo.DownloadVideoDtoFromPrimitives(map[string]any{
			enums.VideosKeyEnum.URL:     entry.GetVideoId(),
			enums.VideosKeyEnum.FORMAT:  s.downloadPlaylistDto.GetFormat(),
			enums.VideosKeyEnum.QUALITY: s.downloadPlaylistDto.GetQuality(),
			enums.VideosKeyEnum.OUTPUT:  outputPath,
		}))

		if s.onVideoEnd != nil {
			s.onVideoEnd(index, total, entry.GetTitle(), err)
		}
		if err != nil {
			failed = append(failed, map[string]any{
				enums.VideosKeyEnum.VIDEO_ID: entry.GetVideoId(),
				enums.VideosKeyEnum.TITLE:    entry.GetTitle(),
				enums.VideosKeyEnum.ERROR:    err.Error(),
			})
			continue
		}
		downloaded = append(downloaded, resultDto.ToArray())
	}
	return downloaded, failed
}

// getOutputFolder <--output o APP_DOWNLOADS_PATH>/<titulo saneado de la playlist>
func (s *DownloadPlaylistService) getOutputFolder(playlistEntity *entities.PlaylistEntity) string {
	base := s.downloadPlaylistDto.GetOutput()
	if base == "" {
		base = configuration.GetInstance().GetDownloadsPath()
	}
	folder := s.fileNameSanitizer.GetSafeFileName(playlistEntity.GetTitle(), playlistEntity.GetPlaylistId())
	return filepath.Join(base, folder)
}
