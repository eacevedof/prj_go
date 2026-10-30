package ytdlp

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"yt-downloader-go/src/core/boot"
)

const (
	YTDLP_RELEASES_URL  = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/"
	YTDLP_LATEST_URL    = "https://github.com/yt-dlp/yt-dlp/releases/latest"
	YTDLP_CHECKSUMS     = "SHA2-256SUMS"
	ytdlpDownloadTimout = 5 * time.Minute
	ytdlpCheckInterval  = 24 * time.Hour
)

var (
	progressPattern = regexp.MustCompile(`\[download\]\s+([\d.]+)% of\s+~?\s*([\d.]+)(B|KiB|MiB|GiB)`)
	sizeUnits       = map[string]float64{"B": 1, "KiB": 1024, "MiB": 1024 * 1024, "GiB": 1024 * 1024 * 1024}
)

// ProgressFunc (bytes leidos, bytes totales).
type ProgressFunc func(read int64, total int64)

// YtDlp motor de descarga: envuelve el binario oficial yt-dlp (github.com/yt-dlp/yt-dlp). Es quien resuelve lo que
// la libreria Go no puede: pistas adaptativas (720p/1080p, audio solo) con el reto JS y el PO token de YouTube.
// El binario se busca en <root>/yt-tools/, junto al ejecutable y en el PATH; si no esta se descarga de la release
// oficial y se verifica su SHA-256 contra el fichero de sumas publicado en la misma release.
type YtDlp struct {
	exePath string
	mutex   sync.Mutex
}

var (
	instance     *YtDlp
	instanceOnce sync.Once
)

func GetInstance() *YtDlp {
	instanceOnce.Do(func() {
		instance = &YtDlp{}
	})
	return instance
}

// EnsureAvailable localiza o descarga el binario; onDownloadStart avisa (stderr) de que se va a bajar.
func (y *YtDlp) EnsureAvailable(onDownloadStart func(message string)) error {
	y.mutex.Lock()
	defer y.mutex.Unlock()
	if y.exePath != "" {
		return nil
	}
	if path := y.findExe(); path != "" {
		y.exePath = path
		y.updateIfOutdated(onDownloadStart)
		return nil
	}
	if onDownloadStart != nil {
		onDownloadStart("yt-dlp not found: downloading the official release from github.com/yt-dlp/yt-dlp (~25 MB)...")
	}
	path, err := y.downloadVerified()
	if err != nil {
		return err
	}
	y.exePath = path
	return nil
}

func (y *YtDlp) GetExePath() string {
	return y.exePath
}

// GetVersion "2026.10.26" o "" si no esta disponible.
func (y *YtDlp) GetVersion() string {
	if y.exePath == "" {
		return ""
	}
	output, err := exec.Command(y.exePath, "--version").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// Download baja la url en <outputBasePath>.<ext> y devuelve la ruta final.
//   - formatSelector: expresion -f de yt-dlp (la decide el servicio: calidad, con/sin ffmpeg).
//   - mergeFormat: "mp4" para unir video+audio; "" si no aplica (audio).
//   - extraArgs: p. ej. --extract-audio.
func (y *YtDlp) Download(url string, outputBasePath string, formatSelector string, mergeFormat string, extraArgs []string, onProgress ProgressFunc) (string, error) {
	if y.exePath == "" {
		return "", errors.New("yt-dlp is not available: call EnsureAvailable first")
	}
	args := []string{
		"--no-playlist", "--no-warnings", "--quiet", "--progress", "--newline", "--no-mtime",
		"-f", formatSelector,
		"-o", outputBasePath + ".%(ext)s",
	}
	if mergeFormat != "" {
		args = append(args, "--merge-output-format", mergeFormat)
	}
	args = append(args, extraArgs...)
	args = append(args, url)

	command := exec.Command(y.exePath, args...)
	stdout, err := command.StdoutPipe()
	if err != nil {
		return "", err
	}
	var stderrBuffer strings.Builder
	command.Stderr = &stderrBuffer
	if err := command.Start(); err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if onProgress != nil {
			if read, total, ok := parseProgress(scanner.Text()); ok {
				onProgress(read, total)
			}
		}
	}
	if err := command.Wait(); err != nil {
		return "", fmt.Errorf("yt-dlp failed: %s", getLastErrorLine(stderrBuffer.String(), err))
	}
	return findOutputFile(outputBasePath)
}

// updateIfOutdated una vez al dia (sello yt-tools/yt-dlp.last-check) compara la version instalada con la ultima
// release (redireccion de /releases/latest) y, si difiere, baja la nueva verificada. YouTube rompe yt-dlp cada pocas
// semanas: sin esto el motor caduca. Cualquier fallo de red se ignora: se sigue con la version instalada.
func (y *YtDlp) updateIfOutdated(onNotice func(message string)) {
	stampPath := filepath.Join(boot.BootPathsEnum.PATH_TOOLS, "yt-dlp.last-check")
	if info, err := os.Stat(stampPath); err == nil && time.Since(info.ModTime()) < ytdlpCheckInterval {
		return
	}
	_ = os.MkdirAll(boot.BootPathsEnum.PATH_TOOLS, 0o777)
	_ = os.WriteFile(stampPath, []byte(time.Now().Format(time.RFC3339)), 0o666)

	latest := getLatestVersion()
	current := y.GetVersion()
	if latest == "" || current == "" || latest == current {
		return
	}
	if onNotice != nil {
		onNotice("yt-dlp " + current + " is outdated: downloading " + latest + " (verified by sha256)...")
	}
	path, err := y.downloadVerified()
	if err != nil {
		if onNotice != nil {
			onNotice("yt-dlp update failed, keeping " + current + ": " + err.Error())
		}
		return
	}
	y.exePath = path
}

// getLatestVersion "2026.10.26" leyendo el Location de /releases/latest (sin API, sin limite de peticiones).
func getLatestVersion() string {
	client := &http.Client{
		Timeout:       15 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	response, err := client.Head(YTDLP_LATEST_URL)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	location := response.Header.Get("Location")
	if i := strings.LastIndex(location, "/tag/"); i >= 0 {
		return strings.TrimSpace(location[i+len("/tag/"):])
	}
	return ""
}

func (y *YtDlp) findExe() string {
	for _, candidate := range []string{
		filepath.Join(boot.BootPathsEnum.PATH_TOOLS, getExeName()),
		filepath.Join(boot.BootPathsEnum.PATH_ROOT, getExeName()),
	} {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate
		}
	}
	if path, err := exec.LookPath(getExeName()); err == nil {
		return path
	}
	return ""
}

// downloadVerified baja el asset de la release + SHA2-256SUMS, compara y solo entonces lo deja en yt-tools/.
func (y *YtDlp) downloadVerified() (string, error) {
	client := &http.Client{Timeout: ytdlpDownloadTimout}
	checksums, err := getBody(client, YTDLP_RELEASES_URL+YTDLP_CHECKSUMS)
	if err != nil {
		return "", fmt.Errorf("yt-dlp checksums download failed: %w", err)
	}
	expectedSha := getExpectedSha(string(checksums), getAssetName())
	if expectedSha == "" {
		return "", errors.New("yt-dlp checksums file has no entry for " + getAssetName())
	}
	binary, err := getBody(client, YTDLP_RELEASES_URL+getAssetName())
	if err != nil {
		return "", fmt.Errorf("yt-dlp download failed: %w", err)
	}
	actual := sha256.Sum256(binary)
	if actualSha := hex.EncodeToString(actual[:]); actualSha != expectedSha {
		return "", fmt.Errorf("yt-dlp download rejected: sha256 %s does not match the published %s", actualSha, expectedSha)
	}
	if err := os.MkdirAll(boot.BootPathsEnum.PATH_TOOLS, 0o777); err != nil {
		return "", err
	}
	exePath := filepath.Join(boot.BootPathsEnum.PATH_TOOLS, getExeName())
	if err := os.WriteFile(exePath, binary, 0o755); err != nil {
		return "", err
	}
	return exePath, nil
}

func getBody(client *http.Client, url string) ([]byte, error) {
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d for %s", response.StatusCode, url)
	}
	return io.ReadAll(response.Body)
}

// getExpectedSha linea "<sha>  <asset>" (o "<sha> *<asset>") del fichero de sumas.
func getExpectedSha(checksums string, assetName string) string {
	for _, line := range strings.Split(checksums, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && strings.TrimPrefix(fields[1], "*") == assetName {
			return strings.ToLower(fields[0])
		}
	}
	return ""
}

func getAssetName() string {
	switch runtime.GOOS {
	case "windows":
		return "yt-dlp.exe"
	case "darwin":
		return "yt-dlp_macos"
	default:
		return "yt-dlp_linux"
	}
}

func getExeName() string {
	if runtime.GOOS == "windows" {
		return "yt-dlp.exe"
	}
	return "yt-dlp"
}

// parseProgress "[download]  45.3% of   10.00MiB at ..." -> (bytes leidos, bytes totales).
func parseProgress(line string) (int64, int64, bool) {
	matches := progressPattern.FindStringSubmatch(line)
	if matches == nil {
		return 0, 0, false
	}
	percent, _ := strconv.ParseFloat(matches[1], 64)
	size, _ := strconv.ParseFloat(matches[2], 64)
	total := int64(size * sizeUnits[matches[3]])
	return int64(float64(total) * percent / 100), total, true
}

func getLastErrorLine(stderr string, err error) string {
	lines := strings.Split(strings.TrimSpace(stderr), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if line := strings.TrimSpace(lines[i]); line != "" {
			return strings.TrimPrefix(line, "ERROR: ")
		}
	}
	return err.Error()
}

// findOutputFile yt-dlp decide la extension: el fichero <base>.* que no sea temporal (.part, .fNNN.*).
func findOutputFile(outputBasePath string) (string, error) {
	matches, _ := filepath.Glob(outputBasePath + ".*")
	for _, match := range matches {
		if strings.HasSuffix(match, ".part") || strings.HasSuffix(match, ".ytdl") {
			continue
		}
		rest := strings.TrimPrefix(match, outputBasePath+".")
		if strings.Count(rest, ".") > 0 {
			continue
		}
		return match, nil
	}
	return "", errors.New("yt-dlp finished but no output file found for " + outputBasePath)
}
