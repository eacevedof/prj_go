# changelog — eaf-yt-dl

La versión es `APP_VERSION` en `src/modules/shared/infrastructure/repositories/configuration/environment_reader_raw_repository.go`
(única fuente: la muestran `menu`, `help` y `version`, y da nombre a `build/dist/yt-{win,linux,macos}-<versión>`).
Lo más nuevo arriba.

## 0.6.0 — 2026-10-30

1. **Linux y macOS**: `make dist` genera los tres binarios `yt-win-<versión>.exe`, `yt-linux-<versión>` (amd64) y
   `yt-macos-<versión>` (Apple Silicon). Mismo menú, carpetas `yt-downloads/` y `yt-tools/` junto al binario; yt-dlp se
   baja como `yt-dlp_linux` / `yt-dlp_macos` verificado por SHA-256. Linux probado en WSL (Ubuntu); macOS sin probar.

## 0.5.3 — 2026-10-30

1. `video_url` (enlace de YouTube) en el resultado de `video-download`, en `video-info` y en cada vídeo de una playlist; en la tarjeta del menú es clicable.

## 0.5.2 — 2026-10-30

1. Carpetas junto al `.exe` renombradas: `yt-downloads\` (vídeos) y `yt-tools\` (yt-dlp). El compilado pasa a llamarse `yt-<versión>.exe`.
2. Fix: el menú entraba en bucle infinito si se agotaba la entrada (stdin cerrado); ahora termina.

## 0.5.1 — 2026-10-30

1. Menú: `select an option [1-6]`, salir es la 6, preguntas en minúsculas.
2. En el menú el resultado se muestra como tarjeta legible (no JSON) y la ruta del fichero es un enlace clicable (OSC 8, Windows Terminal).

## 0.5.0 — 2026-10-30

1. Menú: al descargar un vídeo se lee su info y se listan las **resoluciones reales** (1080p, 720p, 480p…) para elegir una.
2. `--quality` admite alturas `<N>p` (360p… 2160p); `best` ya no tiene tope (la mayor que exista, 4K incluido). `video-info` devuelve `heights`.
3. yt-dlp se **actualiza solo**: una vez al día compara la versión instalada con la última release y, si hay nueva, la baja verificada por SHA-256.

## 0.4.0 — 2026-10-30

1. Motor **yt-dlp** (`APP_ENGINE=ytdlp`, defecto): vídeo hasta 1080p (H.264 + AAC unidos con ffmpeg) y audio m4a sin ffmpeg.
   El binario oficial se descarga solo a `tools\` y se verifica contra el `SHA2-256SUMS` de la release.
2. Nueva calidad `hd1080`; `best` = hasta 1080p. Sin ffmpeg el vídeo se limita al progresivo (360p/720p) con aviso.
3. `APP_ENGINE=native` conserva el motor anterior (solo progresivo, sin binarios externos).

## 0.3.0 — 2026-10-30

1. Nombres de fichero en slug ASCII + fecha-hora (componente `slugger`): "La Niña de mis Ojos" → `la-nina-de-mis-ojos-20261030-153045.mp4`; la carpeta de una playlist es el slug de su título.
2. Todo error (de validación, de YouTube o inesperado) se anota además en `<downloads>/error-YYYYMMDD.log` con el comando y el mensaje.

## 0.2.2 — 2026-10-30

1. Menú entero en inglés; las preguntas de la app en naranja; las calidades en vertical, explicando qué da cada una.
2. Fix: restos en la línea de progreso (`8.8 MBMB`) al pasar de total desconocido a conocido.

## 0.2.1 — 2026-10-30

1. Fix: colores ANSI activados en el conhost clásico (doble clic en el `.exe`); sin color si la terminal no los soporta (antes salían `←[94m`).
2. Menú: opciones de audio avisan de que hay que tener ffmpeg instalado, con el comando `winget` y la URL de descarga.

## 0.2.0 — 2026-10-30

1. Menú interactivo (`menu`, arranca sin argumentos o con doble clic): vídeo, info, playlist, audio, playlist en audio.
2. `.env` opcional: valores por defecto (`APP_ENV=develop`, `APP_DEBUG=1`, `APP_NAME=eaf-yt-dl`, `APP_DEFAULT_LANG=es`).
3. Modo exe suelto: fuera del repo las descargas van a `downloads\` y los logs a `logs\` junto al `.exe`.
4. `make dist` → `build/dist/yt-<versión>.exe` · comando `version`.

## 0.1.0 — 2026-10-29

1. `video-info`, `video-download` (mp4 progresivo o audio m4a vía ffmpeg), `playlist-download`.
2. Detección de descargas recortadas por YouTube (`was cut by YouTube`).
3. `install.bat` para Windows.
