package ffmpeg

import (
	"fmt"
	"os/exec"
	"strings"
)

// AudioExtractor saca la pista de audio de un contenedor de video con ffmpeg, sin recodificar (-acodec copy): un
// mp4 progresivo de YouTube lleva AAC, que cabe tal cual en un .m4a.
type AudioExtractor struct {
	binary string
}

func GetInstance() *AudioExtractor {
	return &AudioExtractor{binary: "ffmpeg"}
}

func (e *AudioExtractor) IsAvailable() bool {
	_, err := exec.LookPath(e.binary)
	return err == nil
}

func (e *AudioExtractor) ExtractAudio(inputPath string, outputPath string) error {
	command := exec.Command(e.binary, "-y", "-loglevel", "error", "-i", inputPath, "-vn", "-acodec", "copy", outputPath)
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}
