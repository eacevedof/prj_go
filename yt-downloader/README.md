# yt-downloader

## 1 - instalación

### 1.1 Solo usarlo en Windows (sin instalar nada)

1. Descargar `eaf-ytdl-x.y.z.exe` (la `x.y.z` es la versión, p. ej. `eaf-ytdl-0.2.0.exe`) de la página de
   **Releases** del repo: `https://github.com/eacevedof/prj_go/releases` (el asset de la release `yt-downloader-vx.y.z`),
   y guardarlo en cualquier carpeta, p. ej. `C:\eaf-ytdl\`. Si no hay release, pedírselo a quien compila (§1.3).
2. Doble clic (o ejecutarlo sin argumentos): se abre un **menú** que pregunta qué hacer (video, info, playlist, audio,
   playlist en audio), la URL y la calidad. No hace falta `.env`: los vídeos caen en `downloads\` junto al `.exe` y los
   logs en `logs\`. La cabecera del menú muestra la versión y su fecha.
3. Solo para audio (`m4a`) hace falta `ffmpeg` en el PATH: `winget install Gyan.FFmpeg`.

### 1.2 Compilar en Windows (sin `make`)

4. Clonar el repo y abrir una consola (`cmd` o PowerShell) en `yt-downloader\`.
5. Ejecutar `.\install.bat` (en `cmd` y en PowerShell; el `.\` es obligatorio en PowerShell). Hace, en orden: instala **Go** con `winget` si falta · instala **ffmpeg** con `winget` si
   falta (solo hace falta para `--format=audio`) · `go mod download` · crea `.env` desde `.env.example` si no existe ·
   compila `build\console.exe`.
6. Si el `.bat` ha tenido que instalar Go, **cerrar la consola, abrir una nueva** (para que coja el PATH) y volver a
   ejecutar `.\install.bat`.
7. Comprobar: `.\build\console.exe help` debe listar los comandos.
8. (Opcional) Revisar `.env` — ver §2.3; sin él ya funciona.

### 1.3 Linux / WSL / macOS

9. Instalar Go ≥ 1.26 (`https://go.dev/dl/`) y `ffmpeg` (`apt install ffmpeg` · `brew install ffmpeg`).
10. Clonar el repo y entrar en `yt-downloader/`.
11. Con `make`: `make build` (crea `.env` y compila `build/console`). Sin `make`:
    `cp .env.example .env && go mod download && go build -o build/console ./cmd/console`.
12. Comprobar: `build/console help` (o `make console c="help"`).
13. **Publicar una versión**: subir `APP_VERSION` y `APP_VERSION_UPDATE` en
    `src/modules/shared/infrastructure/repositories/configuration/environment_reader_raw_repository.go`, anotar el
    cambio en `changelog.md`, `make dist` (genera `build/dist/eaf-ytdl-<versión>.exe`, Windows amd64, ~14 MB) y
    subir ese `.exe` como asset de una release de GitHub con tag `yt-downloader-v<versión>`. `make version` imprime la
    versión actual.

## 2 - uso

### 2.1 Ejemplos

Sin argumentos (o `menu`) se abre el menú interactivo de §1.1; los comandos de abajo son para scripts o consola.

Windows (`cmd` / PowerShell, con `.\` delante; en **Git Bash** usar `./build/console.exe`):

```
.\build\console.exe help
.\build\console.exe video-info --url=https://www.youtube.com/watch?v=dQw4w9WgXcQ
.\build\console.exe video-download --url=dQw4w9WgXcQ --format=video --quality=best
.\build\console.exe video-download --url=dQw4w9WgXcQ --format=video --output=..\mis-videos
.\build\console.exe video-download --url=https://youtu.be/dQw4w9WgXcQ --format=audio --filename=cancion
.\build\console.exe playlist-download --url="https://www.youtube.com/playlist?list=PL4cUxeGkcC9gcy9lrvMJ75z9maRw4byYp" --format=audio
go test -count=1 ./...
```

Linux / WSL / macOS (con `make`; sin él, `build/console <comando>` directamente):

```
make console c="help"
make console c="video-info --url=https://www.youtube.com/watch?v=dQw4w9WgXcQ"
make console c="video-download --url=dQw4w9WgXcQ --format=video --quality=best"
make console c="video-download --url=dQw4w9WgXcQ --format=video --output=../mis-videos"
make console c="video-download --url=https://youtu.be/dQw4w9WgXcQ --format=audio --filename=cancion"
make console c="playlist-download --url=PL4cUxeGkcC9gcy9lrvMJ75z9maRw4byYp --format=video --output=../mis-videos"
make test
```

Sin `--output` el fichero cae en `storage/downloads` del repo, o en `downloads\` junto al `.exe` si se ejecuta fuera
del repo (`APP_DOWNLOADS_PATH`, §2.3). Con `--output` la
carpeta se crea si no existe. Sin ffmpeg el `--format=video` funciona; el `--format=audio` no.

### 2.2 Comandos y opciones

| # | Comando | Opciones |
|---|---|---|
| 1 | `menu` | — · menú interactivo; es lo que arranca sin argumentos |
| 2 | `help` | — · lista los comandos; la cabecera lleva la versión |
| 3 | `version` | — · versión y fecha (`eaf-yt-dl v0.2.0 (2026-10-30)`) |
| 4 | `video-info --url=<url o id>` | — · devuelve JSON con título, duración y formatos |
| 5 | `video-download --url=<url o id>` | `--format=video\|audio` (defecto `video`) · `--quality=best\|hd720\|medium` (defecto `best`) · `--output=<carpeta>` (defecto `APP_DOWNLOADS_PATH`, §2.3; se crea si no existe) · `--filename=<nombre sin extensión>` |
| 6 | `playlist-download --url=<url o id de playlist>` | `--format`, `--quality` y `--output` como en `video-download`. Descarga todos los vídeos, en orden, en la subcarpeta `<output>/<título de la playlist>`; si uno falla (privado, borrado, recortado) lo anota en `failed` y sigue. Exit 1 solo si no baja ninguno |

`--url` acepta la URL completa (`watch?v=`, `youtu.be/`) o el id a secas. Salida: JSON por stdout (apto para `jq`),
progreso y errores por stderr. Exit 0 ok · 1 error de dominio · 2 inesperado (se loguea en `storage/logs`).

### 2.3 Variables de entorno (`.env`, opcional)

El `.env` es **opcional**: toda variable tiene un valor por defecto. Se lee `<raíz>/.env` si existe (lo crea
`install.bat` / `make env` desde `.env.example`); prioridad: variable **real del sistema** > `.env` > defecto.
`%PATH_ROOT%` se sustituye por la raíz. La raíz es la carpeta del `go.mod` si se ejecuta dentro del repo, o la del
`.exe` si no.

| # | Variable | Defecto | Para qué |
|---|---|---|---|
| 1 | `APP_ENV` | `develop` | entorno (`local`, `develop`, `test`, `production`) |
| 2 | `APP_DEBUG` | `1` | `1` muestra el detalle de los errores inesperados |
| 3 | `APP_NAME` | `eaf-yt-dl` | nombre en el menú y en logs |
| 4 | `APP_DEFAULT_LANG` | `es` | idioma por defecto (reservado) |
| 5 | `APP_DOWNLOADS_PATH` | repo: `storage/downloads` · exe: `downloads/` junto al `.exe` | carpeta de descargas; `--output` la sobrescribe |
| 6 | `APP_LOG_PATH` | repo: `storage/logs` · exe: `logs/` junto al `.exe` | logs `YYYY-MM-DD-<tipo>.log` |

## 3 - cómo descarga (y por qué así)

1. Metadatos y formatos con `kkdai/youtube` (cliente ANDROID_VR).
2. `--format=video`: el mp4 **progresivo** (video+audio en un fichero, itag 18/22) con el cliente ANDROID. Con
   ANDROID_VR responde 403.
3. `--format=audio`: el mismo mp4 progresivo y se extrae la pista AAC con **ffmpeg** (`-acodec copy`, sin recodificar)
   a `.m4a`; el mp4 intermedio se borra. Requiere `ffmpeg` en el PATH.
4. No se usan las pistas **adaptativas** (audio solo, 1080p sin audio): YouTube las corta con 403 pasado ~1 MB en todos
   los clientes de la librería (verificado 2026-10-29). Por eso la calidad máxima real es la del progresivo (360p/720p).
5. Hay vídeos cuyo progresivo llega **sin `content_length`** y YouTube sirve solo ~200 KB con la cabecera mp4 de la
   duración entera (ningún cliente de la librería lo da completo, verificado 2026-10-29). Se detecta comparando con
   `bitrate × duración` y se devuelve el error `was cut by YouTube` en vez de un fichero inservible; en una playlist
   ese vídeo queda en `failed`.
6. `playlist-download` lee la lista (id + título de cada vídeo) y delega cada uno en la misma descarga de
   `video-download`, secuencialmente.

## 4 - dependencias

- Go 1.26 · `github.com/kkdai/youtube/v2`
- `ffmpeg` solo para `--format=audio` (Windows: `winget install Gyan.FFmpeg`).
