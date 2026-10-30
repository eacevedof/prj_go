package enums

// VideosKeyEnum claves de entrada (--opcion de consola) y de salida (JSON) del modulo videos.
var VideosKeyEnum = struct {
	URL      string
	FORMAT   string
	QUALITY  string
	OUTPUT   string
	FILENAME string

	VIDEO_ID       string
	TITLE          string
	AUTHOR         string
	DURATION_SECS  string
	VIEWS          string
	PUBLISHED_AT   string
	FORMATS        string
	ITAG           string
	MIME_TYPE      string
	QUALITY_LABEL  string
	BITRATE        string
	CONTENT_LENGTH string
	HAS_AUDIO      string
	HAS_VIDEO      string
	FILE_PATH      string
	BYTES          string

	PLAYLIST_ID string
	OUTPUT_PATH string
	TOTAL       string
	DOWNLOADED  string
	FAILED      string
	ERROR       string
	ENGINE      string
	MAX_HEIGHT  string
}{
	URL:      "url",
	FORMAT:   "format",
	QUALITY:  "quality",
	OUTPUT:   "output",
	FILENAME: "filename",

	VIDEO_ID:       "video_id",
	TITLE:          "title",
	AUTHOR:         "author",
	DURATION_SECS:  "duration_secs",
	VIEWS:          "views",
	PUBLISHED_AT:   "published_at",
	FORMATS:        "formats",
	ITAG:           "itag",
	MIME_TYPE:      "mime_type",
	QUALITY_LABEL:  "quality_label",
	BITRATE:        "bitrate",
	CONTENT_LENGTH: "content_length",
	HAS_AUDIO:      "has_audio",
	HAS_VIDEO:      "has_video",
	FILE_PATH:      "file_path",
	BYTES:          "bytes",

	PLAYLIST_ID: "playlist_id",
	OUTPUT_PATH: "output_path",
	TOTAL:       "total",
	DOWNLOADED:  "downloaded",
	FAILED:      "failed",
	ERROR:       "error",
	ENGINE:      "engine",
	MAX_HEIGHT:  "max_height",
}
