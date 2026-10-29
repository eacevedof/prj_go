package entities

import "yt-downloader-go/src/modules/videos/domain/enums"

// PlaylistEntryEntity un video de la playlist tal como lo lista YouTube: solo id, titulo y duracion (los formatos
// se piden despues, video a video).
type PlaylistEntryEntity struct {
	videoId      string
	title        string
	author       string
	durationSecs int
}

func PlaylistEntryEntityFromPrimitives(primitives map[string]any) *PlaylistEntryEntity {
	videoId, _ := primitives[enums.VideosKeyEnum.VIDEO_ID].(string)
	title, _ := primitives[enums.VideosKeyEnum.TITLE].(string)
	author, _ := primitives[enums.VideosKeyEnum.AUTHOR].(string)
	durationSecs, _ := primitives[enums.VideosKeyEnum.DURATION_SECS].(int)
	return &PlaylistEntryEntity{videoId: videoId, title: title, author: author, durationSecs: durationSecs}
}

func (e *PlaylistEntryEntity) GetVideoId() string   { return e.videoId }
func (e *PlaylistEntryEntity) GetTitle() string     { return e.title }
func (e *PlaylistEntryEntity) GetAuthor() string    { return e.author }
func (e *PlaylistEntryEntity) GetDurationSecs() int { return e.durationSecs }

func (e *PlaylistEntryEntity) ToPrimitives() map[string]any {
	return map[string]any{
		enums.VideosKeyEnum.VIDEO_ID:      e.videoId,
		enums.VideosKeyEnum.TITLE:         e.title,
		enums.VideosKeyEnum.AUTHOR:        e.author,
		enums.VideosKeyEnum.DURATION_SECS: e.durationSecs,
	}
}

type PlaylistEntity struct {
	playlistId string
	title      string
	author     string
	entries    []*PlaylistEntryEntity
}

func PlaylistEntityFromPrimitives(primitives map[string]any, entries []*PlaylistEntryEntity) *PlaylistEntity {
	playlistId, _ := primitives[enums.VideosKeyEnum.PLAYLIST_ID].(string)
	title, _ := primitives[enums.VideosKeyEnum.TITLE].(string)
	author, _ := primitives[enums.VideosKeyEnum.AUTHOR].(string)
	return &PlaylistEntity{playlistId: playlistId, title: title, author: author, entries: entries}
}

func (p *PlaylistEntity) GetPlaylistId() string              { return p.playlistId }
func (p *PlaylistEntity) GetTitle() string                   { return p.title }
func (p *PlaylistEntity) GetAuthor() string                  { return p.author }
func (p *PlaylistEntity) GetEntries() []*PlaylistEntryEntity { return p.entries }
func (p *PlaylistEntity) CountEntries() int                  { return len(p.entries) }
func (p *PlaylistEntity) IsEmpty() bool                      { return len(p.entries) == 0 }
