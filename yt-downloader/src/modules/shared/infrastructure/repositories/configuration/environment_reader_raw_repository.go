package configuration

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"yt-downloader-go/src/core/boot"
	"yt-downloader-go/src/modules/shared/domain/enums"
)

// Version de la app: la muestran el menu, `version` y `help`, y da nombre al instalable (make dist ->
// build/dist/eaf-ytdl-<APP_VERSION>.exe). Al subirla, anotar el cambio en changelog.md.
const (
	APP_VERSION        = "0.2.1"
	APP_VERSION_UPDATE = "2026-10-30"
)

// EnvironmentReaderRawRepository fuente UNICA de lectura de entorno: un getter semantico por cada clave del .env
// que la app consume (claves en AppEnvKeyEnum).
type EnvironmentReaderRawRepository struct{}

var (
	instance     *EnvironmentReaderRawRepository
	instanceOnce sync.Once
)

func GetInstance() *EnvironmentReaderRawRepository {
	instanceOnce.Do(func() {
		instance = &EnvironmentReaderRawRepository{}
	})
	return instance
}

func (e *EnvironmentReaderRawRepository) GetAppVersion() string {
	return APP_VERSION
}

func (e *EnvironmentReaderRawRepository) GetAppVersionUpdate() string {
	return APP_VERSION_UPDATE
}

func (e *EnvironmentReaderRawRepository) GetAppName() string {
	return e.get(enums.AppEnvKeyEnum.APP_NAME)
}

func (e *EnvironmentReaderRawRepository) GetDefaultLang() string {
	return e.get(enums.AppEnvKeyEnum.APP_DEFAULT_LANG)
}

func (e *EnvironmentReaderRawRepository) GetLogPath() string {
	return e.get(enums.AppEnvKeyEnum.APP_LOG_PATH)
}

func (e *EnvironmentReaderRawRepository) GetDownloadsPath() string {
	return e.get(enums.AppEnvKeyEnum.APP_DOWNLOADS_PATH)
}

func (e *EnvironmentReaderRawRepository) GetEnvironment() string {
	return e.get(enums.AppEnvKeyEnum.APP_ENV)
}

func (e *EnvironmentReaderRawRepository) IsProduction() bool {
	return e.GetEnvironment() == boot.ENV.PRODUCTION
}

func (e *EnvironmentReaderRawRepository) IsLocal() bool {
	return e.GetEnvironment() == boot.ENV.LOCAL
}

// IsDebug filter_var(FILTER_VALIDATE_BOOLEAN): "1", "true", "on", "yes".
func (e *EnvironmentReaderRawRepository) IsDebug() bool {
	debug := strings.ToLower(e.get(enums.AppEnvKeyEnum.APP_DEBUG))
	if debug == "on" || debug == "yes" {
		return true
	}
	isDebug, _ := strconv.ParseBool(debug)
	return isDebug
}

func (e *EnvironmentReaderRawRepository) get(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}
