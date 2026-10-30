package downloadvideo

import (
	"os"
	"path/filepath"
	"strings"

	"yt-downloader-go/src/core/config"
	sharedenums "yt-downloader-go/src/modules/shared/domain/enums"
	"yt-downloader-go/src/modules/shared/infrastructure/components/ffmpeg"
	"yt-downloader-go/src/modules/shared/infrastructure/components/slugger"
	"yt-downloader-go/src/modules/shared/infrastructure/exceptions"
	"yt-downloader-go/src/modules/shared/infrastructure/repositories/configuration"
	"yt-downloader-go/src/modules/shared/infrastructure/traits"
	"yt-downloader-go/src/modules/videos/domain/entities"
	"yt-downloader-go/src/modules/videos/domain/enums"
	videosexceptions "yt-downloader-go/src/modules/videos/domain/exceptions"
	"yt-downloader-go/src/modules/videos/domain/services"
	"yt-downloader-go/src/modules/videos/infrastructure/repositories"
)

// ProgressFunc el command la pasa para pintar el avance; (bytes leidos, bytes totales o 0 si no se saben).
type ProgressFunc = repositories.ProgressFunc

// DownloadVideoService descarga el mp4 PROGRESIVO (audio+video en un fichero) y, si se pide audio, extrae la pista
// con ffmpeg. No se usan las pistas adaptativas (audio solo, video 1080p sin audio): YouTube las corta con 403
// pasado ~1 MB en todos los clientes de la libreria (verificado 2026-10-29); el progresivo baja entero.
type DownloadVideoService struct {
	traits.ErrorStackTrait
	videoFieldsValidator          *services.VideoFieldsValidator
	slugger                       *slugger.Slugger
	videosReaderYoutubeRepository *repositories.VideosReaderYoutubeRepository
	videosWriterDiskRepository    *repositories.VideosWriterDiskRepository
	audioExtractor                *ffmpeg.AudioExtractor

	downloadVideoDto *DownloadVideoDto
	onProgress       ProgressFunc
}

func DownloadVideoServiceGetInstance(onProgress ProgressFunc) *DownloadVideoService {
	return &DownloadVideoService{
		videoFieldsValidator:          services.VideoFieldsValidatorGetInstance(),
		slugger:                       slugger.GetInstance(),
		videosReaderYoutubeRepository: repositories.VideosReaderYoutubeRepositoryGetInstance(),
		videosWriterDiskRepository:    repositories.VideosWriterDiskRepositoryGetInstance(),
		audioExtractor:                ffmpeg.GetInstance(),
		onProgress:                    onProgress,
	}
}

// Invoke error: *FieldsException (400), *VideosException (404, 409) o uno inesperado.
func (s *DownloadVideoService) Invoke(downloadVideoDto *DownloadVideoDto) (*DownloadVideoResultDto, error) {
	s.downloadVideoDto = downloadVideoDto
	if err := s.failIfWrongInput(); err != nil {
		return nil, err
	}

	videoId := s.videoFieldsValidator.GetVideoIdFromUrl(downloadVideoDto.GetUrl())
	videoEntity, err := s.videosReaderYoutubeRepository.GetVideoById(videoId)
	if err != nil {
		return nil, err
	}

	formatEntity, err := s.getFormatOrFail(videoEntity)
	if err != nil {
		return nil, err
	}

	baseName := s.getBaseName(videoEntity)
	filePath, bytes, err := s.videosWriterDiskRepository.SaveFormat(
		videoId, formatEntity, s.getOutputFolder(), baseName+"."+s.getVideoExtension(formatEntity), s.onProgress,
	)
	if err != nil {
		return nil, err
	}

	if s.isAudio() {
		if filePath, bytes, err = s.extractAudio(filePath, baseName); err != nil {
			return nil, err
		}
	}

	return DownloadVideoResultDtoFromPrimitives(map[string]any{
		enums.VideosKeyEnum.VIDEO_ID:  videoId,
		enums.VideosKeyEnum.TITLE:     videoEntity.GetTitle(),
		enums.VideosKeyEnum.FILE_PATH: filePath,
		enums.VideosKeyEnum.BYTES:     bytes,
		enums.VideosKeyEnum.FORMAT:    formatEntity.ToPrimitives(),
	}), nil
}

func (s *DownloadVideoService) failIfWrongInput() error {
	errors := s.videoFieldsValidator.GetDownloadErrors(
		s.downloadVideoDto.GetUrl(), s.downloadVideoDto.GetFormat(), s.downloadVideoDto.GetQuality(),
	)
	if len(errors) > 0 {
		return exceptions.FieldsExceptionBadRequestCustom(
			s.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
		)
	}
	if s.isAudio() && !s.audioExtractor.IsAvailable() {
		return videosexceptions.ConflictCustom(config.Tr(enums.VideosMessageEnum.FFMPEG_NOT_AVAILABLE))
	}
	return nil
}

func (s *DownloadVideoService) isAudio() bool {
	return s.downloadVideoDto.GetFormat() == enums.VideoFormatEnum.AUDIO
}

// getFormatOrFail siempre un progresivo; para audio la calidad da igual (la pista AAC es la misma), se coge la mejor.
func (s *DownloadVideoService) getFormatOrFail(videoEntity *entities.VideoEntity) (*entities.VideoFormatEntity, error) {
	quality := s.downloadVideoDto.GetQuality()
	if s.isAudio() {
		quality = enums.VideoQualityEnum.BEST
	}
	formatEntity := videoEntity.GetBestProgressiveFormat(quality)
	if formatEntity == nil {
		return nil, videosexceptions.ConflictCustom(config.Tr(
			enums.VideosMessageEnum.NO_DOWNLOADABLE_FORMAT_FOR,
			s.downloadVideoDto.GetFormat(), quality, videoEntity.GetVideoId(),
		))
	}
	return formatEntity, nil
}

// extractAudio <base>.mp4 -> <base>.m4a y se borra el mp4 intermedio.
func (s *DownloadVideoService) extractAudio(videoPath string, baseName string) (string, int64, error) {
	audioPath := filepath.Join(filepath.Dir(videoPath), baseName+".m4a")
	defer os.Remove(videoPath)
	if err := s.audioExtractor.ExtractAudio(videoPath, audioPath); err != nil {
		_ = os.Remove(audioPath)
		return "", 0, err
	}
	info, err := os.Stat(audioPath)
	if err != nil {
		return "", 0, err
	}
	return audioPath, info.Size(), nil
}

func (s *DownloadVideoService) getOutputFolder() string {
	if output := s.downloadVideoDto.GetOutput(); output != "" {
		return output
	}
	return configuration.GetInstance().GetDownloadsPath()
}

// getBaseName slug del nombre pedido o del titulo + fecha-hora, sin extension: "la-nina-de-mis-ojos-20261030-153045".
func (s *DownloadVideoService) getBaseName(videoEntity *entities.VideoEntity) string {
	name := s.downloadVideoDto.GetFileName()
	if name == "" {
		name = videoEntity.GetTitle()
	}
	return s.slugger.GetSlugWithNow(name, videoEntity.GetVideoId())
}

func (s *DownloadVideoService) getVideoExtension(formatEntity *entities.VideoFormatEntity) string {
	if strings.HasPrefix(formatEntity.GetMimeType(), "video/webm") {
		return "webm"
	}
	return "mp4"
}
