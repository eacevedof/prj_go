package enums

var VideosMessageEnum = struct {
	FIELD_IS_REQUIRED            string
	INVALID_YOUTUBE_URL          string
	INVALID_FORMAT_ALLOWED_X     string
	INVALID_QUALITY_ALLOWED_X    string
	VIDEO_X_NOT_FOUND            string
	VIDEO_X_NOT_AVAILABLE        string
	NO_DOWNLOADABLE_FORMAT_FOR   string
	OUTPUT_FOLDER_X_NOT_WRITABLE string
	FFMPEG_NOT_AVAILABLE         string
}{
	FFMPEG_NOT_AVAILABLE:         "ffmpeg is required to extract the audio track and it is not in PATH",
	FIELD_IS_REQUIRED:            "field is required",
	INVALID_YOUTUBE_URL:          "invalid YouTube url or video id",
	INVALID_FORMAT_ALLOWED_X:     "invalid format, allowed: {0}",
	INVALID_QUALITY_ALLOWED_X:    "invalid quality, allowed: {0}",
	VIDEO_X_NOT_FOUND:            "Video {0} not found",
	VIDEO_X_NOT_AVAILABLE:        "Video {0} is not available: {1}",
	NO_DOWNLOADABLE_FORMAT_FOR:   "No downloadable {0} format with quality {1} for video {2}",
	OUTPUT_FOLDER_X_NOT_WRITABLE: "Output folder {0} is not writable",
}
