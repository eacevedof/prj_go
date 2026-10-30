package enums

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
// BEST = la mayor disponible hasta 1080p.
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

// GetQualityMaxHeight altura maxima en pixeles para el selector de yt-dlp.
func GetQualityMaxHeight(quality string) int {
	switch quality {
	case VideoQualityEnum.HD720:
		return 720
	case VideoQualityEnum.MEDIUM:
		return 360
	case VideoQualityEnum.SMALL:
		return 240
	default:
		return 1080
	}
}
