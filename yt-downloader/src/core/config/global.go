package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"yt-downloader-go/src/core/boot"
	"yt-downloader-go/src/modules/shared/domain/enums"
	"yt-downloader-go/src/modules/shared/infrastructure/components/logger"
)

var (
	quotedEnvValuePattern = regexp.MustCompile(`^(["'])(.*)`)
	envCommentPattern     = regexp.MustCompile(`(^|\s+)#.*$`)
)

// AppbootLoadenv carga <root>/.env si existe y rellena con valores por defecto lo que falte: el .env es OPCIONAL
// (quien solo tiene el .exe no configura nada). Prioridad: variable REAL del OS > .env > defecto.
func AppbootLoadenv() error {
	pathEnv := filepath.Join(boot.BootPathsEnum.PATH_ROOT, ".env")
	appbootLoadenvfile(pathEnv)
	appbootSetdefaultenvkeys()
	return nil
}

// Lgerr storage/logs/YYYY-MM-DD-<type>.log (mismo esquema que LogTrait/Logger).
func Lgerr(variable any, title string, logType string) {
	content, isString := variable.(string)
	if !isString {
		content = fmt.Sprintf("%+v", variable)
	}
	if title != "" {
		content = "(lgerr) <<  " + title + " >>\n" + content
	} else {
		content = "(lgerr)\n" + content
	}
	logger.GetInstance().WriteLog(logType+".log", content, boot.BootPathsEnum.PATH_LOGS)
}

// Tr sustituye {0}, {1}... por los args. Siempre en ingles, el msgid tal cual.
func Tr(msgId string, args ...any) string {
	msgChanged := strings.TrimSpace(msgId)
	if msgChanged == "" {
		return ""
	}
	for i, arg := range args {
		msgChanged = strings.ReplaceAll(msgChanged, "{"+strconv.Itoa(i)+"}", fmt.Sprint(arg))
	}
	return msgChanged
}

func appbootLoadenvfile(pathEnv string) {
	envFile, err := os.Open(pathEnv)
	if err != nil {
		return
	}
	defer envFile.Close()

	replacer := strings.NewReplacer(
		"%PATH_ROOT%", filepath.ToSlash(boot.BootPathsEnum.PATH_ROOT),
		"%PATH_SRC%", filepath.ToSlash(boot.BootPathsEnum.PATH_SRC),
	)

	scanner := bufio.NewScanner(envFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "#") || !strings.Contains(line, "=") {
			continue
		}
		name, value, _ := strings.Cut(line, "=")
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		// la variable de entorno REAL (OS) gana sobre el fichero .env
		if _, isSet := os.LookupEnv(name); isSet {
			continue
		}
		_ = os.Setenv(name, replacer.Replace(appbootGetenvfilevalue(value)))
	}
}

// appbootGetenvfilevalue entre comillas se toma el interior literal; sin comillas se corta el comentario en linea.
func appbootGetenvfilevalue(rawValue string) string {
	value := strings.TrimSpace(rawValue)
	if matches := quotedEnvValuePattern.FindStringSubmatch(value); matches != nil {
		if closingIndex := strings.LastIndex(matches[2], matches[1]); closingIndex >= 0 {
			return matches[2][:closingIndex]
		}
	}
	return strings.TrimSpace(envCommentPattern.ReplaceAllString(value, ""))
}

// appbootSetdefaultenvkeys toda clave ausente o vacia toma su valor por defecto.
func appbootSetdefaultenvkeys() {
	defaults := map[string]string{
		enums.AppEnvKeyEnum.APP_ENV:            enums.AppEnvDefaultEnum.APP_ENV,
		enums.AppEnvKeyEnum.APP_DEBUG:          enums.AppEnvDefaultEnum.APP_DEBUG,
		enums.AppEnvKeyEnum.APP_NAME:           enums.AppEnvDefaultEnum.APP_NAME,
		enums.AppEnvKeyEnum.APP_DEFAULT_LANG:   enums.AppEnvDefaultEnum.APP_DEFAULT_LANG,
		enums.AppEnvKeyEnum.APP_LOG_PATH:       boot.BootPathsEnum.PATH_LOGS,
		enums.AppEnvKeyEnum.APP_DOWNLOADS_PATH: boot.BootPathsEnum.PATH_DOWNLOADS,
	}
	for key, value := range defaults {
		if strings.TrimSpace(os.Getenv(key)) == "" {
			_ = os.Setenv(key, value)
		}
	}
}
