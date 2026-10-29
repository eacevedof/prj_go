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

// VideoQualityEnum calidades que admite --quality (las etiquetas "quality" de YouTube para los formatos
// progresivos). BEST elige la mayor disponible.
var VideoQualityEnum = struct {
	BEST   string
	HD720  string
	MEDIUM string
	SMALL  string
}{
	BEST:   "best",
	HD720:  "hd720",
	MEDIUM: "medium",
	SMALL:  "small",
}

func GetVideoQualities() []string {
	return []string{VideoQualityEnum.BEST, VideoQualityEnum.HD720, VideoQualityEnum.MEDIUM, VideoQualityEnum.SMALL}
}
