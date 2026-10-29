package exceptions

import "yt-downloader-go/src/modules/shared/domain/enums"

// VideosException error de dominio de videos con su codigo HTTP equivalente (400, 404, 409).
type VideosException struct {
	message string
	code    int
}

func (e *VideosException) Error() string {
	return e.message
}

func (e *VideosException) GetCode() int {
	return e.code
}

func BadRequestCustom(message string) *VideosException {
	return &VideosException{message: message, code: enums.ExceptionEnum.CODE_BAD_REQUEST}
}

func NotFoundCustom(message string) *VideosException {
	return &VideosException{message: message, code: enums.ExceptionEnum.CODE_NOT_FOUND}
}

func ConflictCustom(message string) *VideosException {
	return &VideosException{message: message, code: enums.ExceptionEnum.CODE_CONFLICT}
}
