package services_test

import (
	"testing"

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

func TestFileNameSanitizerGetSafeFileName(t *testing.T) {
	sanitizer := services.FileNameSanitizerGetInstance()
	tests.AssertEquals(t, "Mi video parte 1 2", sanitizer.GetSafeFileName(`Mi video: parte 1/2?`, "x"))
	tests.AssertEquals(t, "fallback", sanitizer.GetSafeFileName("???", "fallback"))
}
