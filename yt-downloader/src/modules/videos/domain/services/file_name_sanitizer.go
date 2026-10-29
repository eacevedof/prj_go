package services

import (
	"regexp"
	"strings"
)

var (
	forbiddenFileNameChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]+`)
	repeatedSpaces         = regexp.MustCompile(`\s+`)
)

const FILE_NAME_MAX_LENGTH = 120

// FileNameSanitizer convierte el titulo del video en un nombre de fichero valido en Windows y Linux.
type FileNameSanitizer struct{}

func FileNameSanitizerGetInstance() *FileNameSanitizer {
	return &FileNameSanitizer{}
}

// GetSafeFileName "Mi video: parte 1/2?" -> "Mi video parte 1 2"; vacio -> fallback.
func (s *FileNameSanitizer) GetSafeFileName(title string, fallback string) string {
	name := forbiddenFileNameChars.ReplaceAllString(title, " ")
	name = repeatedSpaces.ReplaceAllString(name, " ")
	name = strings.Trim(name, " .")
	if len(name) > FILE_NAME_MAX_LENGTH {
		name = strings.TrimSpace(name[:FILE_NAME_MAX_LENGTH])
	}
	if name == "" {
		return fallback
	}
	return name
}
