package getvideoinfo

import (
	"yt-downloader-go/src/core/config"
	sharedenums "yt-downloader-go/src/modules/shared/domain/enums"
	"yt-downloader-go/src/modules/shared/infrastructure/exceptions"
	"yt-downloader-go/src/modules/shared/infrastructure/traits"
	"yt-downloader-go/src/modules/videos/domain/services"
	"yt-downloader-go/src/modules/videos/infrastructure/repositories"
)

type GetVideoInfoService struct {
	traits.ErrorStackTrait
	videoFieldsValidator          *services.VideoFieldsValidator
	videosReaderYoutubeRepository *repositories.VideosReaderYoutubeRepository

	getVideoInfoDto *GetVideoInfoDto
}

func GetVideoInfoServiceGetInstance() *GetVideoInfoService {
	return &GetVideoInfoService{
		videoFieldsValidator:          services.VideoFieldsValidatorGetInstance(),
		videosReaderYoutubeRepository: repositories.VideosReaderYoutubeRepositoryGetInstance(),
	}
}

// Invoke error: *FieldsException (400), *VideosException (404, 409) o uno inesperado.
func (s *GetVideoInfoService) Invoke(getVideoInfoDto *GetVideoInfoDto) (*GetVideoInfoResultDto, error) {
	s.getVideoInfoDto = getVideoInfoDto
	if err := s.failIfWrongInput(); err != nil {
		return nil, err
	}

	videoId := s.videoFieldsValidator.GetVideoIdFromUrl(getVideoInfoDto.GetUrl())
	videoEntity, err := s.videosReaderYoutubeRepository.GetVideoById(videoId)
	if err != nil {
		return nil, err
	}
	return GetVideoInfoResultDtoFromPrimitives(map[string]any{"video": videoEntity.ToPrimitives()}), nil
}

func (s *GetVideoInfoService) failIfWrongInput() error {
	errors := s.videoFieldsValidator.GetUrlErrors(s.getVideoInfoDto.GetUrl())
	if len(errors) == 0 {
		return nil
	}
	return exceptions.FieldsExceptionBadRequestCustom(
		s.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
	)
}
