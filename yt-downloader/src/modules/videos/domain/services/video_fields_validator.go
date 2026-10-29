package services

import (
	"regexp"
	"slices"
	"strings"

	"yt-downloader-go/src/core/config"
	"yt-downloader-go/src/modules/videos/domain/enums"
)

// youtubeVideoIdPattern un id de video son 11 caracteres [A-Za-z0-9_-]
var youtubeVideoIdPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

// youtubeUrlPatterns formas de url que se aceptan; el grupo 1 es el id
var youtubeUrlPatterns = []*regexp.Regexp{
	regexp.MustCompile(`^(?:https?://)?(?:www\.|m\.|music\.)?youtube\.com/watch\?(?:.*&)?v=([A-Za-z0-9_-]{11})`),
	regexp.MustCompile(`^(?:https?://)?(?:www\.)?youtube\.com/(?:shorts|embed|live|v)/([A-Za-z0-9_-]{11})`),
	regexp.MustCompile(`^(?:https?://)?youtu\.be/([A-Za-z0-9_-]{11})`),
}

// VideoFieldsValidator reglas de entrada de los comandos de videos. Devuelve un error "campo: mensaje" por cada
// regla incumplida; vacio = valido.
type VideoFieldsValidator struct{}

func VideoFieldsValidatorGetInstance() *VideoFieldsValidator {
	return &VideoFieldsValidator{}
}

func (v *VideoFieldsValidator) GetUrlErrors(url string) []string {
	if url == "" {
		return []string{enums.VideosKeyEnum.URL + ": " + config.Tr(enums.VideosMessageEnum.FIELD_IS_REQUIRED)}
	}
	if v.GetVideoIdFromUrl(url) == "" {
		return []string{enums.VideosKeyEnum.URL + ": " + config.Tr(enums.VideosMessageEnum.INVALID_YOUTUBE_URL)}
	}
	return nil
}

func (v *VideoFieldsValidator) GetDownloadErrors(url string, format string, quality string) []string {
	errors := v.GetUrlErrors(url)
	if !slices.Contains(enums.GetVideoFormats(), format) {
		errors = append(errors, enums.VideosKeyEnum.FORMAT+": "+config.Tr(
			enums.VideosMessageEnum.INVALID_FORMAT_ALLOWED_X, strings.Join(enums.GetVideoFormats(), ", "),
		))
	}
	if !slices.Contains(enums.GetVideoQualities(), quality) {
		errors = append(errors, enums.VideosKeyEnum.QUALITY+": "+config.Tr(
			enums.VideosMessageEnum.INVALID_QUALITY_ALLOWED_X, strings.Join(enums.GetVideoQualities(), ", "),
		))
	}
	return errors
}

// GetVideoIdFromUrl el id de video de una url de YouTube (watch, shorts, embed, youtu.be) o del propio id; "" si
// no se reconoce.
func (v *VideoFieldsValidator) GetVideoIdFromUrl(url string) string {
	url = strings.TrimSpace(url)
	if youtubeVideoIdPattern.MatchString(url) {
		return url
	}
	for _, pattern := range youtubeUrlPatterns {
		if matches := pattern.FindStringSubmatch(url); matches != nil {
			return matches[1]
		}
	}
	return ""
}
