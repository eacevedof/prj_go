package repositories

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/kkdai/youtube/v2"

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
	isLengthFromMetadata := format.ContentLength > 0
	w.fillContentLengthIfUnknown(client, video, format)
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
	if expected := getExpectedBytes(total, format, isLengthFromMetadata); written < expected {
		_ = os.Remove(partPath)
		return "", 0, exceptions.ConflictCustom(config.Tr(
			enums.VideosMessageEnum.DOWNLOAD_X_TRUNCATED, videoId, strconv.FormatInt(written, 10), strconv.FormatInt(expected, 10),
		))
	}
	if err := os.Rename(partPath, filePath); err != nil {
		_ = os.Remove(partPath)
		return "", 0, err
	}
	return filePath, written, nil
}

// getExpectedBytes el minimo de bytes para dar la descarga por buena. Con content_length en los METADATOS, ese.
// Sin el, la mitad de bitrate x duracion: hay videos (content_length 0) para los que YouTube anuncia (HEAD) y sirve
// un fichero de ~240 KB con la cabecera mp4 de la duracion entera (verificado 2026-10-29, itag 18 cliente
// ANDROID); el tamaño servido es coherente consigo mismo, asi que solo se detecta comparando con lo que deberia pesar.
func getExpectedBytes(total int64, format *youtube.Format, isLengthFromMetadata bool) int64 {
	if isLengthFromMetadata {
		return total
	}
	durationMs, _ := strconv.ParseInt(format.ApproxDurationMs, 10, 64)
	if format.Bitrate <= 0 || durationMs <= 0 {
		return total
	}
	expected := int64(format.Bitrate) * durationMs / 8000
	return expected / 2
}

// fillContentLengthIfUnknown algunos formatos llegan sin content_length; con 0 la libreria baja todo en UNA
// peticion y YouTube la corta a los pocos cientos de KB (verificado 2026-10-29). Con el tamaño (HEAD a la url del
// stream) la libreria descarga por rangos de 10 MB y valida cada trozo. Si el HEAD falla se sigue sin tamaño y
// lo pilla el chequeo de truncado.
func (w *VideosWriterDiskRepository) fillContentLengthIfUnknown(client *youtube.Client, video *youtube.Video, format *youtube.Format) {
	if format.ContentLength > 0 {
		return
	}
	streamUrl, err := client.GetStreamURL(video, format)
	if err != nil {
		return
	}
	httpClient := &http.Client{Timeout: 20 * time.Second}
	response, err := httpClient.Head(streamUrl)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK && response.ContentLength > 0 {
		format.ContentLength = response.ContentLength
	}
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
