package repositories

import (
	"io"
	"os"
	"path/filepath"

	"yt-downloader-go/src/core/config"
	"yt-downloader-go/src/modules/videos/domain/entities"
	"yt-downloader-go/src/modules/videos/domain/enums"
	"yt-downloader-go/src/modules/videos/domain/exceptions"
)

// ProgressFunc (bytes leidos, bytes totales o 0 si no se saben).
type ProgressFunc func(read int64, total int64)

// VideosWriterDiskRepository escribe el stream en disco. Descarga a <fichero>.part y renombra al acabar: un corte
// no deja un mp4 a medias con nombre definitivo.
type VideosWriterDiskRepository struct {
	videosReaderYoutubeRepository *VideosReaderYoutubeRepository
}

func VideosWriterDiskRepositoryGetInstance() *VideosWriterDiskRepository {
	return &VideosWriterDiskRepository{videosReaderYoutubeRepository: VideosReaderYoutubeRepositoryGetInstance()}
}

// SaveFormat devuelve la ruta absoluta y los bytes escritos. error: *VideosException (409 carpeta no escribible)
// o uno inesperado.
func (w *VideosWriterDiskRepository) SaveFormat(
	videoId string, formatEntity *entities.VideoFormatEntity, outputFolder string, fileName string, onProgress ProgressFunc,
) (string, int64, error) {
	if err := os.MkdirAll(outputFolder, 0o777); err != nil {
		return "", 0, exceptions.ConflictCustom(config.Tr(enums.VideosMessageEnum.OUTPUT_FOLDER_X_NOT_WRITABLE, outputFolder))
	}
	filePath, err := filepath.Abs(filepath.Join(outputFolder, fileName))
	if err != nil {
		return "", 0, err
	}

	client, video, format, err := w.videosReaderYoutubeRepository.GetStreamReader(videoId, formatEntity.GetItag())
	if err != nil {
		return "", 0, err
	}
	stream, total, err := client.GetStream(video, format)
	if err != nil {
		return "", 0, err
	}
	defer stream.Close()

	partPath := filePath + ".part"
	partFile, err := os.Create(partPath)
	if err != nil {
		return "", 0, exceptions.ConflictCustom(config.Tr(enums.VideosMessageEnum.OUTPUT_FOLDER_X_NOT_WRITABLE, outputFolder))
	}

	written, err := io.Copy(partFile, &progressReader{reader: stream, total: total, onProgress: onProgress})
	closeErr := partFile.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(partPath)
		return "", 0, err
	}
	if err := os.Rename(partPath, filePath); err != nil {
		_ = os.Remove(partPath)
		return "", 0, err
	}
	return filePath, written, nil
}

// progressReader cuenta los bytes que pasan y avisa al callback.
type progressReader struct {
	reader     io.Reader
	total      int64
	read       int64
	onProgress ProgressFunc
}

func (p *progressReader) Read(buffer []byte) (int, error) {
	n, err := p.reader.Read(buffer)
	p.read += int64(n)
	if p.onProgress != nil && n > 0 {
		p.onProgress(p.read, p.total)
	}
	return n, err
}
