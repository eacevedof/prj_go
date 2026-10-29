package getvideoinfo

import "yt-downloader-go/src/modules/shared/application"

type GetVideoInfoDto struct {
	url string
}

func GetVideoInfoDtoFromPrimitives(primitives map[string]any) *GetVideoInfoDto {
	return &GetVideoInfoDto{url: application.GetTrimmedString(primitives, "url")}
}

func (d *GetVideoInfoDto) GetUrl() string {
	return d.url
}
