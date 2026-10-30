package enums

import (
	"regexp"
	"slices"
	"strconv"
)

// VideoFormatEnum que se descarga: VIDEO = mp4 progresivo (video+audio en un solo fichero, YouTube solo lo sirve
// hasta 720p) · AUDIO = solo la pista de audio (m4a).
var VideoFormatEnum = struct {
	VIDEO string
	AUDIO string
}{
	VIDEO: "video",
	AUDIO: "audio",
}

func GetVideoFormats() []string {
	return []string{VideoFormatEnum.VIDEO, VideoFormatEnum.AUDIO}
}

// VideoQualityEnum calidades que admite --quality (etiquetas "quality" de YouTube). Con el motor yt-dlp la calidad es
// una altura maxima (GetMaxHeight); con el progresivo nativo solo existen hd720 (itag 22) y medium (itag 18).
// Tambien se admite "<N>p" (1440p, 2160p...). BEST = la mayor que exista.
var VideoQualityEnum = struct {
	BEST   string
	HD1080 string
	HD720  string
	MEDIUM string
	SMALL  string
}{
	BEST:   "best",
	HD1080: "hd1080",
	HD720:  "hd720",
	MEDIUM: "medium",
	SMALL:  "small",
}

func GetVideoQualities() []string {
	return []string{VideoQualityEnum.BEST, VideoQualityEnum.HD1080, VideoQualityEnum.HD720, VideoQualityEnum.MEDIUM, VideoQualityEnum.SMALL}
}

var (
	qualityHeightPattern = regexp.MustCompile(`^(\d{3,4})p$`)
	qualityLabelPattern  = regexp.MustCompile(`^(\d{3,4})p`)
)

// IsValidQuality una de GetVideoQualities() o una altura "<N>p" (360p, 720p, 1080p, 1440p, 2160p...).
func IsValidQuality(quality string) bool {
	return slices.Contains(GetVideoQualities(), quality) || qualityHeightPattern.MatchString(quality)
}

// GetHeightFromQualityLabel "1080p60" -> 1080; "" o "medium" -> 0.
func GetHeightFromQualityLabel(label string) int {
	matches := qualityLabelPattern.FindStringSubmatch(label)
	if matches == nil {
		return 0
	}
	height, _ := strconv.Atoi(matches[1])
	return height
}

// GetQualityMaxHeight altura maxima en pixeles para el selector de yt-dlp. BEST = sin tope (la mayor que exista).
func GetQualityMaxHeight(quality string) int {
	if height := GetHeightFromQualityLabel(quality); height > 0 {
		return height
	}
	switch quality {
	case VideoQualityEnum.HD720:
		return 720
	case VideoQualityEnum.MEDIUM:
		return 360
	case VideoQualityEnum.SMALL:
		return 240
	case VideoQualityEnum.HD1080:
		return 1080
	default:
		return 4320
	}
}
