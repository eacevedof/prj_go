package enums

// AppEnvKeyEnum claves del .env que la app consume. Solo las usan EnvironmentReaderRawRepository y AppbootLoadenv():
// el resto del codigo nunca escribe un nombre de variable de entorno.
var AppEnvKeyEnum = struct {
	APP_ENV            string
	APP_DEBUG          string
	APP_NAME           string
	APP_LOG_PATH       string
	APP_DOWNLOADS_PATH string
}{
	APP_ENV:            "APP_ENV",
	APP_DEBUG:          "APP_DEBUG",
	APP_NAME:           "APP_NAME",
	APP_LOG_PATH:       "APP_LOG_PATH",
	APP_DOWNLOADS_PATH: "APP_DOWNLOADS_PATH",
}
