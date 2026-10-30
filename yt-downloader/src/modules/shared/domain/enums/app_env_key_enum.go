package enums

// AppEnvKeyEnum claves del .env que la app consume. Solo las usan EnvironmentReaderRawRepository y AppbootLoadenv():
// el resto del codigo nunca escribe un nombre de variable de entorno. Todas tienen valor por defecto (el .env es
// opcional): ver AppEnvDefaultEnum.
var AppEnvKeyEnum = struct {
	APP_ENV            string
	APP_DEBUG          string
	APP_NAME           string
	APP_DEFAULT_LANG   string
	APP_LOG_PATH       string
	APP_DOWNLOADS_PATH string
}{
	APP_ENV:            "APP_ENV",
	APP_DEBUG:          "APP_DEBUG",
	APP_NAME:           "APP_NAME",
	APP_DEFAULT_LANG:   "APP_DEFAULT_LANG",
	APP_LOG_PATH:       "APP_LOG_PATH",
	APP_DOWNLOADS_PATH: "APP_DOWNLOADS_PATH",
}

// AppEnvDefaultEnum valores que se asumen cuando ni el OS ni el .env definen la clave. Las rutas no estan aqui: salen
// de BootPathsEnum (dependen de si se ejecuta desde el repo o como .exe suelto).
var AppEnvDefaultEnum = struct {
	APP_ENV          string
	APP_DEBUG        string
	APP_NAME         string
	APP_DEFAULT_LANG string
}{
	APP_ENV:          "develop",
	APP_DEBUG:        "1",
	APP_NAME:         "eaf-yt-dl",
	APP_DEFAULT_LANG: "es",
}
