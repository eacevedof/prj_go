package enums

var VideosMessageEnum = struct {
	FIELD_IS_REQUIRED             string
	INVALID_YOUTUBE_URL           string
	INVALID_FORMAT_ALLOWED_X      string
	INVALID_QUALITY_ALLOWED_X     string
	VIDEO_X_NOT_FOUND             string
	VIDEO_X_NOT_AVAILABLE         string
	NO_DOWNLOADABLE_FORMAT_FOR    string
	OUTPUT_FOLDER_X_NOT_WRITABLE  string
	FFMPEG_NOT_AVAILABLE          string
	INVALID_YOUTUBE_PLAYLIST_URL  string
	PLAYLIST_X_NOT_FOUND          string
	PLAYLIST_X_IS_EMPTY           string
	PLAYLIST_X_NOTHING_DOWNLOADED string
	DOWNLOAD_X_TRUNCATED          string
}{
	FFMPEG_NOT_AVAILABLE:          "ffmpeg is required to extract the audio track and it is not in PATH",
	FIELD_IS_REQUIRED:             "field is required",
	INVALID_YOUTUBE_URL:           "invalid YouTube url or video id",
	INVALID_FORMAT_ALLOWED_X:      "invalid format, allowed: {0}",
	INVALID_QUALITY_ALLOWED_X:     "invalid quality, allowed: {0}",
	VIDEO_X_NOT_FOUND:             "Video {0} not found",
	VIDEO_X_NOT_AVAILABLE:         "Video {0} is not available: {1}",
	NO_DOWNLOADABLE_FORMAT_FOR:    "No downloadable {0} format with quality {1} for video {2}",
	OUTPUT_FOLDER_X_NOT_WRITABLE:  "Output folder {0} is not writable",
	INVALID_YOUTUBE_PLAYLIST_URL:  "invalid YouTube playlist url or playlist id",
	PLAYLIST_X_NOT_FOUND:          "Playlist {0} not found",
	PLAYLIST_X_IS_EMPTY:           "Playlist {0} has no videos",
	PLAYLIST_X_NOTHING_DOWNLOADED: "Playlist {0}: none of the {1} videos could be downloaded",
	DOWNLOAD_X_TRUNCATED:          "Download of video {0} was cut by YouTube: {1} of {2} bytes",
}
