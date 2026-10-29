package getvideoinfo

import "yt-downloader-go/src/modules/videos/domain/entities"

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

func (d *GetVideoInfoResultDto) ToArray() entities.VideoPrimitives {
	return d.video
}
