package services_test

import (
	"strings"
	"testing"

	"yt-downloader-go/src/modules/shared/infrastructure/components/slugger"
	"yt-downloader-go/src/modules/videos/domain/services"
	"yt-downloader-go/tests"
)

func TestVideoFieldsValidatorGetVideoIdFromUrl(t *testing.T) {
	validator := services.VideoFieldsValidatorGetInstance()
	cases := map[string]string{
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ":          "dQw4w9WgXcQ",
		"https://www.youtube.com/watch?list=abc&v=dQw4w9WgXcQ": "dQw4w9WgXcQ",
		"https://youtu.be/dQw4w9WgXcQ?t=10":                    "dQw4w9WgXcQ",
		"https://youtube.com/shorts/dQw4w9WgXcQ":               "dQw4w9WgXcQ",
		"https://m.youtube.com/watch?v=dQw4w9WgXcQ":            "dQw4w9WgXcQ",
		"dQw4w9WgXcQ":           "dQw4w9WgXcQ",
		"https://vimeo.com/123": "",
		"not a url":             "",
		"":                      "",
	}
	for url, expected := range cases {
		tests.AssertEquals(t, expected, validator.GetVideoIdFromUrl(url), url)
	}
}

func TestVideoFieldsValidatorGetDownloadErrors(t *testing.T) {
	validator := services.VideoFieldsValidatorGetInstance()

	tests.AssertEmpty(t, validator.GetDownloadErrors("https://youtu.be/dQw4w9WgXcQ", "video", "best"))
	tests.AssertEmpty(t, validator.GetDownloadErrors("dQw4w9WgXcQ", "audio", "hd720"))

	errors := validator.GetDownloadErrors("", "mp3", "4k")
	tests.AssertEquals(t, 3, len(errors))
	tests.AssertStringContainsString(t, "url: field is required", errors[0])
	tests.AssertStringContainsString(t, "format: invalid format", errors[1])
	tests.AssertStringContainsString(t, "quality: invalid quality", errors[2])
}

func TestVideoFieldsValidatorGetPlaylistIdFromUrl(t *testing.T) {
	validator := services.VideoFieldsValidatorGetInstance()
	cases := map[string]string{
		"https://www.youtube.com/playlist?list=PLrAXtmErZgOeiKm4sgNOknGvNjby9efdf":            "PLrAXtmErZgOeiKm4sgNOknGvNjby9efdf",
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ&list=PLrAXtmErZgOeiKm4sgNOknGvNjby9efdf": "PLrAXtmErZgOeiKm4sgNOknGvNjby9efdf",
		"https://music.youtube.com/playlist?list=OLAK5uy_kX1234567890abcdef":                  "OLAK5uy_kX1234567890abcdef",
		"PLrAXtmErZgOeiKm4sgNOknGvNjby9efdf":                                                  "PLrAXtmErZgOeiKm4sgNOknGvNjby9efdf",
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ":                                         "",
		"dQw4w9WgXcQ":           "",
		"https://vimeo.com/123": "",
		"":                      "",
	}
	for url, expected := range cases {
		tests.AssertEquals(t, expected, validator.GetPlaylistIdFromUrl(url), url)
	}
}

func TestVideoFieldsValidatorGetPlaylistDownloadErrors(t *testing.T) {
	validator := services.VideoFieldsValidatorGetInstance()

	tests.AssertEmpty(t, validator.GetPlaylistDownloadErrors("PLrAXtmErZgOeiKm4sgNOknGvNjby9efdf", "audio", "best"))

	errors := validator.GetPlaylistDownloadErrors("dQw4w9WgXcQ", "mp3", "4k")
	tests.AssertEquals(t, 3, len(errors))
	tests.AssertStringContainsString(t, "url: invalid YouTube playlist url", errors[0])
	tests.AssertStringContainsString(t, "format: invalid format", errors[1])
	tests.AssertStringContainsString(t, "quality: invalid quality", errors[2])
}

func TestSluggerGetSlug(t *testing.T) {
	slug := slugger.GetInstance()
	tests.AssertEquals(t, "la-nina-de-mis-ojos", slug.GetSlug("La Niña de mis Ojos", "x"))
	tests.AssertEquals(t, "mi-video-parte-1-2", slug.GetSlug(`Mi vídeo: parte 1/2?`, "x"))
	tests.AssertEquals(t, "rock-and-roll", slug.GetSlug("Rock & Roll!!", "x"))
	tests.AssertEquals(t, "fallback", slug.GetSlug("???", "fallback"))
	tests.AssertEquals(t, "file", slug.GetSlug("", ""))
	tests.AssertEquals(t, true, strings.HasPrefix(slug.GetSlugWithNow("Eva Ayllón - Mal Paso (En Vivo)", "x"), "eva-ayllon-mal-paso-en-vivo-20"))
}
