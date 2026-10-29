package repositories

import (
	"errors"
	"strings"
	"sync"

	"github.com/kkdai/youtube/v2"

	"yt-downloader-go/src/core/config"
	"yt-downloader-go/src/modules/videos/domain/entities"
	"yt-downloader-go/src/modules/videos/domain/enums"
	"yt-downloader-go/src/modules/videos/domain/exceptions"
)

// VideosReaderYoutubeRepository lectura de metadatos y streams de YouTube (kkdai/youtube). Es el unico sitio que
// conoce la libreria: el dominio solo ve VideoEntity.
//
// Dos clientes porque YouTube no sirve todo a uno solo (verificado 2026-10-29):
//   - ANDROID_VR (el por defecto de la libreria): metadatos y pistas ADAPTATIVAS (audio m4a, video sin audio).
//     El mp4 progresivo (itag 18) responde 403.
//   - ANDROID: el mp4 PROGRESIVO (audio+video). Las adaptativas responden 403.
type VideosReaderYoutubeRepository struct {
	adaptiveClient    *youtube.Client
	progressiveClient *youtube.Client
	// cache por proceso por cliente: el writer necesita el *youtube.Video y su Format para abrir el stream
	adaptiveVideos    map[string]*youtube.Video
	progressiveVideos map[string]*youtube.Video
	mutex             sync.Mutex
}

var (
	readerInstance     *VideosReaderYoutubeRepository
	readerInstanceOnce sync.Once
)

func VideosReaderYoutubeRepositoryGetInstance() *VideosReaderYoutubeRepository {
	readerInstanceOnce.Do(func() {
		readerInstance = &VideosReaderYoutubeRepository{
			adaptiveClient:    newClient(youtube.AndroidVRClient),
			progressiveClient: newClient(youtube.AndroidClient),
			adaptiveVideos:    map[string]*youtube.Video{},
			progressiveVideos: map[string]*youtube.Video{},
		}
	})
	return readerInstance
}

// newClient la libreria fija el tipo de cliente desde youtube.DefaultClient en la primera llamada; se fuerza aqui
// con una llamada que no sale a la red (GetStreamURL con un formato vacio: devuelve ErrCipherNotFound tras fijarlo).
func newClient(clientInfo youtube.ClientInfo) *youtube.Client {
	youtube.DefaultClient = clientInfo
	client := &youtube.Client{}
	_, _ = client.GetStreamURL(&youtube.Video{}, &youtube.Format{})
	return client
}

// GetVideoById error: *VideosException (404 no existe · 409 privado/edad/no disponible) o uno inesperado.
func (r *VideosReaderYoutubeRepository) GetVideoById(videoId string) (*entities.VideoEntity, error) {
	video, err := r.getYoutubeVideo(videoId, r.adaptiveClient, r.adaptiveVideos)
	if err != nil {
		return nil, err
	}

	formats := make([]*entities.VideoFormatEntity, 0, len(video.Formats))
	for _, format := range video.Formats {
		formats = append(formats, r.getFormatEntity(format))
	}

	return entities.VideoEntityFromPrimitives(map[string]any{
		enums.VideosKeyEnum.VIDEO_ID:      video.ID,
		enums.VideosKeyEnum.TITLE:         video.Title,
		enums.VideosKeyEnum.AUTHOR:        video.Author,
		enums.VideosKeyEnum.DURATION_SECS: int(video.Duration.Seconds()),
		enums.VideosKeyEnum.VIEWS:         video.Views,
		enums.VideosKeyEnum.PUBLISHED_AT:  video.PublishDate.Format("2006-01-02"),
	}, formats), nil
}

// GetStreamReader el cliente, video y formato con los que abrir el stream del itag pedido (firma ya descifrada
// por la libreria). Progresivo (audio+video) -> cliente ANDROID; adaptativo -> ANDROID_VR.
func (r *VideosReaderYoutubeRepository) GetStreamReader(videoId string, itag int) (*youtube.Client, *youtube.Video, *youtube.Format, error) {
	adaptiveVideo, err := r.getYoutubeVideo(videoId, r.adaptiveClient, r.adaptiveVideos)
	if err != nil {
		return nil, nil, nil, err
	}
	format := findFormat(adaptiveVideo, itag)
	if format == nil {
		return nil, nil, nil, exceptions.NotFoundCustom(config.Tr(enums.VideosMessageEnum.VIDEO_X_NOT_FOUND, videoId))
	}
	if format.InitRange != nil {
		return r.adaptiveClient, adaptiveVideo, format, nil
	}

	progressiveVideo, err := r.getYoutubeVideo(videoId, r.progressiveClient, r.progressiveVideos)
	if err != nil {
		return nil, nil, nil, err
	}
	if format = findFormat(progressiveVideo, itag); format == nil {
		return nil, nil, nil, exceptions.NotFoundCustom(config.Tr(enums.VideosMessageEnum.VIDEO_X_NOT_FOUND, videoId))
	}
	return r.progressiveClient, progressiveVideo, format, nil
}

func findFormat(video *youtube.Video, itag int) *youtube.Format {
	for i := range video.Formats {
		if video.Formats[i].ItagNo == itag {
			return &video.Formats[i]
		}
	}
	return nil
}

func (r *VideosReaderYoutubeRepository) getYoutubeVideo(videoId string, client *youtube.Client, cache map[string]*youtube.Video) (*youtube.Video, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if video, isCached := cache[videoId]; isCached {
		return video, nil
	}

	video, err := client.GetVideo(videoId)
	if err != nil {
		return nil, r.getDomainError(videoId, err)
	}
	cache[videoId] = video
	return video, nil
}

// getDomainError los errores que la libreria tipifica se traducen a dominio; el resto sube como inesperado.
func (r *VideosReaderYoutubeRepository) getDomainError(videoId string, err error) error {
	var playabilityError youtube.ErrPlayabiltyStatus
	switch {
	case errors.Is(err, youtube.ErrVideoIDMinLength), errors.Is(err, youtube.ErrInvalidCharactersInVideoID):
		return exceptions.NotFoundCustom(config.Tr(enums.VideosMessageEnum.VIDEO_X_NOT_FOUND, videoId))
	case errors.Is(err, youtube.ErrVideoPrivate), errors.Is(err, youtube.ErrLoginRequired), errors.Is(err, youtube.ErrNotPlayableInEmbed):
		return exceptions.ConflictCustom(config.Tr(enums.VideosMessageEnum.VIDEO_X_NOT_AVAILABLE, videoId, err.Error()))
	case errors.As(err, &playabilityError):
		if strings.Contains(strings.ToLower(playabilityError.Reason), "unavailable") {
			return exceptions.NotFoundCustom(config.Tr(enums.VideosMessageEnum.VIDEO_X_NOT_FOUND, videoId))
		}
		return exceptions.ConflictCustom(config.Tr(enums.VideosMessageEnum.VIDEO_X_NOT_AVAILABLE, videoId, playabilityError.Reason))
	}
	return err
}

func (r *VideosReaderYoutubeRepository) getFormatEntity(format youtube.Format) *entities.VideoFormatEntity {
	mimeType := format.MimeType
	// adaptativo = solo una pista (InitRange presente); progresivo = audio+video
	isAdaptive := format.InitRange != nil
	isVideoMime := strings.HasPrefix(mimeType, "video/")
	hasVideo := isVideoMime
	hasAudio := !isVideoMime || !isAdaptive
	if isVideoMime && isAdaptive {
		hasAudio = false
	}
	return entities.VideoFormatEntityFromPrimitives(map[string]any{
		enums.VideosKeyEnum.ITAG:           format.ItagNo,
		enums.VideosKeyEnum.MIME_TYPE:      mimeType,
		enums.VideosKeyEnum.QUALITY:        format.Quality,
		enums.VideosKeyEnum.QUALITY_LABEL:  format.QualityLabel,
		enums.VideosKeyEnum.BITRATE:        format.Bitrate,
		enums.VideosKeyEnum.CONTENT_LENGTH: format.ContentLength,
		enums.VideosKeyEnum.HAS_AUDIO:      hasAudio,
		enums.VideosKeyEnum.HAS_VIDEO:      hasVideo,
	})
}
