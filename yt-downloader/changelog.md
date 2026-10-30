# changelog — eaf-yt-dl

La versión es `APP_VERSION` en `src/modules/shared/infrastructure/repositories/configuration/environment_reader_raw_repository.go`
(única fuente: la muestran `menu`, `help` y `version`, y da nombre a `build/dist/eaf-ytdl-<versión>.exe`).
Lo más nuevo arriba.

## 0.2.0 — 2026-10-30

1. Menú interactivo (`menu`, arranca sin argumentos o con doble clic): vídeo, info, playlist, audio, playlist en audio.
2. `.env` opcional: valores por defecto (`APP_ENV=develop`, `APP_DEBUG=1`, `APP_NAME=eaf-yt-dl`, `APP_DEFAULT_LANG=es`).
3. Modo exe suelto: fuera del repo las descargas van a `downloads\` y los logs a `logs\` junto al `.exe`.
4. `make dist` → `build/dist/eaf-ytdl-<versión>.exe` · comando `version`.

## 0.1.0 — 2026-10-29

1. `video-info`, `video-download` (mp4 progresivo o audio m4a vía ffmpeg), `playlist-download`.
2. Detección de descargas recortadas por YouTube (`was cut by YouTube`).
3. `install.bat` para Windows.
