package downloadvideo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"yt-downloader-go/src/core/config"
	sharedenums "yt-downloader-go/src/modules/shared/domain/enums"
	"yt-downloader-go/src/modules/shared/infrastructure/components/ffmpeg"
	"yt-downloader-go/src/modules/shared/infrastructure/components/slugger"
	"yt-downloader-go/src/modules/shared/infrastructure/components/ytdlp"
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

// DownloadVideoService dos motores:
//   - yt-dlp (APP_ENGINE=ytdlp, defecto): el binario oficial baja video+audio adaptativos (hasta 1080p) y los une con
//     ffmpeg; el audio solo (m4a) sale directo sin ffmpeg. Sin ffmpeg se limita al progresivo.
//   - nativo (APP_ENGINE=native o si yt-dlp no se puede obtener): el mp4 PROGRESIVO (audio+video en un fichero,
//     360p/720p) con la libreria Go, y para audio extrae la pista con ffmpeg. Las pistas adaptativas no se usan
//     aqui: YouTube las corta con 403 pasado ~1 MB en todos los clientes de la libreria (verificado 2026-10-29).
//
// Los metadatos (titulo, formatos) salen siempre de la libreria Go.
type DownloadVideoService struct {
	traits.ErrorStackTrait
	videoFieldsValidator          *services.VideoFieldsValidator
	slugger                       *slugger.Slugger
	videosReaderYoutubeRepository *repositories.VideosReaderYoutubeRepository
	videosWriterDiskRepository    *repositories.VideosWriterDiskRepository
	audioExtractor                *ffmpeg.AudioExtractor
	ytDlp                         *ytdlp.YtDlp

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
		ytDlp:                         ytdlp.GetInstance(),
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

	baseName := s.getBaseName(videoEntity)
	if s.isYtDlpEngine() {
		return s.invokeWithYtDlp(videoId, videoEntity, baseName)
	}

	formatEntity, err := s.getFormatOrFail(videoEntity)
	if err != nil {
		return nil, err
	}

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
	if s.isAudio() && !s.audioExtractor.IsAvailable() && !s.isYtDlpEngine() {
		return videosexceptions.ConflictCustom(config.Tr(enums.VideosMessageEnum.FFMPEG_NOT_AVAILABLE))
	}
	return nil
}

// isYtDlpEngine APP_ENGINE=ytdlp y el binario disponible (se descarga la primera vez); si no, aviso y motor nativo.
func (s *DownloadVideoService) isYtDlpEngine() bool {
	if !configuration.GetInstance().IsYtDlpEngine() {
		return false
	}
	if err := s.ytDlp.EnsureAvailable(s.printNotice); err != nil {
		s.printNotice(config.Tr(enums.VideosMessageEnum.YTDLP_NOT_AVAILABLE_X, err.Error()))
		return false
	}
	return true
}

func (s *DownloadVideoService) printNotice(message string) {
	fmt.Fprintln(os.Stderr, "  "+message)
}

// invokeWithYtDlp selector -f segun formato/calidad/ffmpeg; yt-dlp decide la extension final.
func (s *DownloadVideoService) invokeWithYtDlp(videoId string, videoEntity *entities.VideoEntity, baseName string) (*DownloadVideoResultDto, error) {
	outputFolder := s.getOutputFolder()
	if err := os.MkdirAll(outputFolder, 0o777); err != nil {
		return nil, videosexceptions.ConflictCustom(config.Tr(enums.VideosMessageEnum.OUTPUT_FOLDER_X_NOT_WRITABLE, outputFolder))
	}
	maxHeight := enums.GetQualityMaxHeight(s.downloadVideoDto.GetQuality())
	selector, mergeFormat := s.getYtDlpSelector(maxHeight)
	filePath, err := s.ytDlp.Download(
		enums.GetVideoUrl(videoId), filepath.Join(outputFolder, baseName), selector, mergeFormat, nil, ytdlp.ProgressFunc(s.onProgress),
	)
	if err != nil {
		return nil, videosexceptions.ConflictCustom(err.Error())
	}
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	return DownloadVideoResultDtoFromPrimitives(map[string]any{
		enums.VideosKeyEnum.VIDEO_ID:  videoId,
		enums.VideosKeyEnum.TITLE:     videoEntity.GetTitle(),
		enums.VideosKeyEnum.FILE_PATH: filePath,
		enums.VideosKeyEnum.BYTES:     info.Size(),
		enums.VideosKeyEnum.FORMAT: entities.VideoFormatPrimitives{
			enums.VideosKeyEnum.ENGINE:     "yt-dlp " + s.ytDlp.GetVersion(),
			enums.VideosKeyEnum.QUALITY:    s.downloadVideoDto.GetQuality(),
			enums.VideosKeyEnum.MAX_HEIGHT: maxHeight,
		},
	}), nil
}

// getYtDlpSelector audio: la mejor pista m4a (sin ffmpeg). Video con ffmpeg: mejor video H.264 hasta maxHeight + mejor
// audio m4a, unidos en mp4 (fallback a cualquier codec). Sin ffmpeg: solo progresivos (b), que YouTube da hasta 720p.
func (s *DownloadVideoService) getYtDlpSelector(maxHeight int) (string, string) {
	if s.isAudio() {
		return "ba[ext=m4a]/ba", ""
	}
	height := fmt.Sprint(maxHeight)
	if !s.audioExtractor.IsAvailable() {
		s.printNotice(config.Tr(enums.VideosMessageEnum.FFMPEG_NEEDED_FOR_QUALITY_X, "the best progressive mp4 up to "+height+"p"))
		return "b[height<=" + height + "][ext=mp4]/b[height<=" + height + "]/b", ""
	}
	return "bv*[height<=" + height + "][vcodec^=avc1]+ba[ext=m4a]/bv*[height<=" + height + "]+ba/b[height<=" + height + "]/b", "mp4"
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
