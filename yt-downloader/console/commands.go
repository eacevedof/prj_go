// Package console registro de comandos de consola.
//
//	./console <command-alias> --clave=valor .. --clave=valor
package console

import (
	"yt-downloader-go/src/console/application"
	videoscommands "yt-downloader-go/src/modules/videos/infrastructure/commands"
)

// GetCommands alias => command, que hace (una linea) y los argumentos de un ejemplo. El orden es el del help.
func GetCommands() []application.CommandConfig {
	return []application.CommandConfig{
		{
			Alias: "help",
			Service: func(input []string) application.InterfaceConsole {
				return application.HelpCommandGetInstance(input, GetCommands())
			},
			Info:    "lista los comandos con un ejemplo de cada uno",
			Example: "",
		},
		{
			Alias: "video-info",
			Service: func(input []string) application.InterfaceConsole {
				return videoscommands.GetVideoInfoCommandGetInstance(input)
			},
			Info:    "metadatos del video y formatos descargables (JSON)",
			Example: "--url=https://www.youtube.com/watch?v=<id>",
		},
		{
			Alias: "video-download",
			Service: func(input []string) application.InterfaceConsole {
				return videoscommands.DownloadVideoCommandGetInstance(input)
			},
			Info:    "descarga el mp4 (video+audio) o solo el audio (m4a, via ffmpeg) en APP_DOWNLOADS_PATH o --output",
			Example: "--url=https://youtu.be/<id> --format=video|audio --quality=best|hd720|medium --output=<carpeta> --filename=<nombre>",
		},
		{
			Alias: "playlist-download",
			Service: func(input []string) application.InterfaceConsole {
				return videoscommands.DownloadPlaylistCommandGetInstance(input)
			},
			Info:    "descarga todos los videos de una playlist en una subcarpeta con su titulo; sigue si uno falla",
			Example: "--url=https://www.youtube.com/playlist?list=<id> --format=video|audio --quality=best|hd720|medium --output=<carpeta>",
		},
	}
}
