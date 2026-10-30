package entities

import (
	"slices"
	"sort"

	"yt-downloader-go/src/modules/videos/domain/enums"
)

// VideoFormatPrimitives un formato descargable tal cual sale por JSON.
type VideoFormatPrimitives map[string]any

// VideoPrimitives el video tal cual sale por JSON.
type VideoPrimitives map[string]any

// VideoFormatEntity un stream que ofrece YouTube: progresivo (audio+video) o adaptativo (solo una pista).
type VideoFormatEntity struct {
	itag          int
	mimeType      string
	quality       string
	qualityLabel  string
	bitrate       int
	contentLength int64
	hasAudio      bool
	hasVideo      bool
}

func VideoFormatEntityFromPrimitives(primitives map[string]any) *VideoFormatEntity {
	itag, _ := primitives[enums.VideosKeyEnum.ITAG].(int)
	mimeType, _ := primitives[enums.VideosKeyEnum.MIME_TYPE].(string)
	quality, _ := primitives[enums.VideosKeyEnum.QUALITY].(string)
	qualityLabel, _ := primitives[enums.VideosKeyEnum.QUALITY_LABEL].(string)
	bitrate, _ := primitives[enums.VideosKeyEnum.BITRATE].(int)
	contentLength, _ := primitives[enums.VideosKeyEnum.CONTENT_LENGTH].(int64)
	hasAudio, _ := primitives[enums.VideosKeyEnum.HAS_AUDIO].(bool)
	hasVideo, _ := primitives[enums.VideosKeyEnum.HAS_VIDEO].(bool)
	return &VideoFormatEntity{
		itag:          itag,
		mimeType:      mimeType,
		quality:       quality,
		qualityLabel:  qualityLabel,
		bitrate:       bitrate,
		contentLength: contentLength,
		hasAudio:      hasAudio,
		hasVideo:      hasVideo,
	}
}

func (f *VideoFormatEntity) GetItag() int            { return f.itag }
func (f *VideoFormatEntity) GetMimeType() string     { return f.mimeType }
func (f *VideoFormatEntity) GetQuality() string      { return f.quality }
func (f *VideoFormatEntity) GetQualityLabel() string { return f.qualityLabel }
func (f *VideoFormatEntity) GetBitrate() int         { return f.bitrate }
func (f *VideoFormatEntity) GetContentLength() int64 { return f.contentLength }
func (f *VideoFormatEntity) HasAudio() bool          { return f.hasAudio }
func (f *VideoFormatEntity) HasVideo() bool          { return f.hasVideo }

// IsProgressive audio y video en el mismo fichero: lo que se puede guardar sin mezclar (sin ffmpeg).
func (f *VideoFormatEntity) IsProgressive() bool {
	return f.hasAudio && f.hasVideo
}

func (f *VideoFormatEntity) IsAudioOnly() bool {
	return f.hasAudio && !f.hasVideo
}

func (f *VideoFormatEntity) ToPrimitives() VideoFormatPrimitives {
	return VideoFormatPrimitives{
		enums.VideosKeyEnum.ITAG:           f.itag,
		enums.VideosKeyEnum.MIME_TYPE:      f.mimeType,
		enums.VideosKeyEnum.QUALITY:        f.quality,
		enums.VideosKeyEnum.QUALITY_LABEL:  f.qualityLabel,
		enums.VideosKeyEnum.BITRATE:        f.bitrate,
		enums.VideosKeyEnum.CONTENT_LENGTH: f.contentLength,
		enums.VideosKeyEnum.HAS_AUDIO:      f.hasAudio,
		enums.VideosKeyEnum.HAS_VIDEO:      f.hasVideo,
	}
}

// VideoEntity metadatos de un video de YouTube y sus formatos.
type VideoEntity struct {
	videoId      string
	title        string
	author       string
	durationSecs int
	views        int
	publishedAt  string
	formats      []*VideoFormatEntity
}

func VideoEntityFromPrimitives(primitives map[string]any, formats []*VideoFormatEntity) *VideoEntity {
	videoId, _ := primitives[enums.VideosKeyEnum.VIDEO_ID].(string)
	title, _ := primitives[enums.VideosKeyEnum.TITLE].(string)
	author, _ := primitives[enums.VideosKeyEnum.AUTHOR].(string)
	durationSecs, _ := primitives[enums.VideosKeyEnum.DURATION_SECS].(int)
	views, _ := primitives[enums.VideosKeyEnum.VIEWS].(int)
	publishedAt, _ := primitives[enums.VideosKeyEnum.PUBLISHED_AT].(string)
	return &VideoEntity{
		videoId:      videoId,
		title:        title,
		author:       author,
		durationSecs: durationSecs,
		views:        views,
		publishedAt:  publishedAt,
		formats:      formats,
	}
}

func (v *VideoEntity) GetVideoId() string               { return v.videoId }
func (v *VideoEntity) GetTitle() string                 { return v.title }
func (v *VideoEntity) GetAuthor() string                { return v.author }
func (v *VideoEntity) GetDurationSecs() int             { return v.durationSecs }
func (v *VideoEntity) GetFormats() []*VideoFormatEntity { return v.formats }

// GetBestProgressiveFormat el mp4 con audio de mayor bitrate; con quality != best, solo el de esa calidad.
func (v *VideoEntity) GetBestProgressiveFormat(quality string) *VideoFormatEntity {
	var best *VideoFormatEntity
	for _, format := range v.formats {
		if !format.IsProgressive() {
			continue
		}
		if quality != enums.VideoQualityEnum.BEST && format.GetQuality() != quality {
			continue
		}
		if best == nil || format.GetBitrate() > best.GetBitrate() {
			best = format
		}
	}
	return best
}

// GetBestAudioFormat la pista de audio mp4 (m4a) de mayor bitrate.
func (v *VideoEntity) GetBestAudioFormat() *VideoFormatEntity {
	var best *VideoFormatEntity
	for _, format := range v.formats {
		if !format.IsAudioOnly() {
			continue
		}
		if best == nil || format.GetBitrate() > best.GetBitrate() {
			best = format
		}
	}
	return best
}

// GetAvailableHeights alturas de video distintas (1080, 720, 360...) de mayor a menor, segun las etiquetas
// "1080p"/"720p60" de todos los formatos con video (progresivos y adaptativos).
func (v *VideoEntity) GetAvailableHeights() []int {
	var heights []int
	for _, format := range v.formats {
		if !format.HasVideo() {
			continue
		}
		height := enums.GetHeightFromQualityLabel(format.GetQualityLabel())
		if height > 0 && !slices.Contains(heights, height) {
			heights = append(heights, height)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(heights)))
	return heights
}

func (v *VideoEntity) ToPrimitives() VideoPrimitives {
	formats := make([]VideoFormatPrimitives, 0, len(v.formats))
	for _, format := range v.formats {
		formats = append(formats, format.ToPrimitives())
	}
	return VideoPrimitives{
		enums.VideosKeyEnum.VIDEO_ID:      v.videoId,
		enums.VideosKeyEnum.TITLE:         v.title,
		enums.VideosKeyEnum.AUTHOR:        v.author,
		enums.VideosKeyEnum.DURATION_SECS: v.durationSecs,
		enums.VideosKeyEnum.VIEWS:         v.views,
		enums.VideosKeyEnum.PUBLISHED_AT:  v.publishedAt,
		enums.VideosKeyEnum.FORMATS:       formats,
		enums.VideosKeyEnum.HEIGHTS:       v.GetAvailableHeights(),
	}
}
