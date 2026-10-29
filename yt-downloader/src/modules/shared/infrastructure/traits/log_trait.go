package traits

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"yt-downloader-go/src/core/boot"
	"yt-downloader-go/src/modules/shared/infrastructure/components/logger"
	"yt-downloader-go/src/modules/shared/infrastructure/components/server"
)

// LogTrait logs de la app en storage/logs/YYYY-MM-DD-<nivel>.log. Cada entrada lleva "cli-<ip>" y su nivel.
type LogTrait struct{}

func NewLogTrait() LogTrait {
	return LogTrait{}
}

func (t LogTrait) LogDebug(mixed any, title string) {
	t.writeLevelLog("debug", "DEBUG", mixed, title)
}

func (t LogTrait) LogInfo(mixed any, title string) {
	t.writeLevelLog("info", "INFO", mixed, title)
}

func (t LogTrait) LogWarning(mixed any, title string) {
	t.writeLevelLog("warning", "WARNING", mixed, title)
}

func (t LogTrait) LogError(mixed any, title string) {
	t.writeLevelLog("error", "ERROR", mixed, title)
}

// LogException lo que va en todo error inesperado: tipo, mensaje y traza.
func (t LogTrait) LogException(err error, title string) {
	if title == "" {
		title = "ERROR"
	}
	t.writeLevelLog("error", "ERROR", strings.Join([]string{
		"ex type:\n\t" + fmt.Sprintf("%T", err),
		"ex message:\n\t" + err.Error(),
		"ex trace:\n\t" + string(debug.Stack()),
	}, "\n"), title)
}

// LogDebugCli solo en consola: al debug.log y ademas a pantalla con color.
func (t LogTrait) LogDebugCli(mixed any, title string) {
	if !boot.IsCli() {
		return
	}
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	print := getMixedAsString(mixed)
	logger.GetInstance().WriteLog("debug", "["+now+"]["+title+"]\n"+print, boot.BootPathsEnum.PATH_LOGS)
	fmt.Println("\033[97m[" + now + "]\033[0m\033[93m[" + title + "]\033[0m\n\033[92m" + print + "\033[0m")
}

func (t LogTrait) writeLevelLog(fileName string, level string, mixed any, title string) {
	content := []string{"cli-" + server.GetInstance().GetServerIp()}
	if title != "" {
		content = append(content, title)
	}
	content = append(content, getMixedAsString(mixed))
	logger.GetInstance().WriteLog(fileName, "["+level+"] "+strings.Join(content, "\n"), boot.BootPathsEnum.PATH_LOGS)
}

func getMixedAsString(mixed any) string {
	if text, isString := mixed.(string); isString {
		return text
	}
	return fmt.Sprintf("%+v", mixed)
}
