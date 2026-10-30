package getvideoinfo

import (
	"yt-downloader-go/src/modules/videos/domain/entities"
	"yt-downloader-go/src/modules/videos/domain/enums"
)

type GetVideoInfoResultDto struct {
	video entities.VideoPrimitives
}

func GetVideoInfoResultDtoFromPrimitives(primitives map[string]any) *GetVideoInfoResultDto {
	video, _ := primitives["video"].(entities.VideoPrimitives)
	return &GetVideoInfoResultDto{video: video}
}

func (d *GetVideoInfoResultDto) GetVideo() entities.VideoPrimitives {
	return d.video
}

func (d *GetVideoInfoResultDto) GetTitle() string {
	title, _ := d.video[enums.VideosKeyEnum.TITLE].(string)
	return title
}

// GetAvailableHeights 1080, 720, 360... de mayor a menor.
func (d *GetVideoInfoResultDto) GetAvailableHeights() []int {
	heights, _ := d.video[enums.VideosKeyEnum.HEIGHTS].([]int)
	return heights
}

func (d *GetVideoInfoResultDto) ToArray() entities.VideoPrimitives {
	return d.video
}
