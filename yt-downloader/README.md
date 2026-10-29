# yt-downloader

Descargador de YouTube por consola en Go, con la arquitectura DDD de `task-go` (prj-task-cruds): `cmd/console` →
`console/commands.go` → `src/modules/videos/{application,domain,infrastructure}`. Sin API ni BD.
Sustituye al PoC `youtube-downloader/` (legacy).

## 1 - arrancar

### 1.1 Windows (sin `make`)

```
install.bat                                  # instala Go y ffmpeg (winget) si faltan, go mod download, .env, build\console.exe
build\console.exe help
build\console.exe video-info --url=https://www.youtube.com/watch?v=dQw4w9WgXcQ
build\console.exe video-download --url=dQw4w9WgXcQ --format=video --quality=best --output=D:\videos
build\console.exe video-download --url=https://youtu.be/dQw4w9WgXcQ --format=audio --filename=cancion
go test -count=1 ./...
```

Si `install.bat` instala Go, hay que **abrir una consola nueva** (PATH) y lanzarlo otra vez. Sin ffmpeg el
`--format=video` funciona; el `--format=audio` no.

A mano, lo mismo que hace el `.bat`:

```
winget install --id GoLang.Go -e
winget install --id Gyan.FFmpeg -e
go mod download
copy .env.example .env
go build -o build\console.exe .\cmd\console
```

### 1.2 Linux / WSL / macOS (con `make`)

```
make build                                   # .env desde .env.example + build/console
make console c="help"
make console c="video-info --url=https://www.youtube.com/watch?v=dQw4w9WgXcQ"
make console c="video-download --url=dQw4w9WgXcQ --format=video --quality=best --output=/tmp/videos"
make console c="video-download --url=https://youtu.be/dQw4w9WgXcQ --format=audio --filename=cancion"
make test
```

Sin `make`: `cp .env.example .env && go build -o build/console ./cmd/console` y luego `build/console <comando>`.

### 1.3 Comandos y opciones

| # | Comando | Opciones |
|---|---|---|
| 1 | `help` | — |
| 2 | `video-info --url=<url o id>` | — · devuelve JSON con título, duración y formatos |
| 3 | `video-download --url=<url o id>` | `--format=video\|audio` (defecto `video`) · `--quality=best\|hd720\|medium` (defecto `best`) · `--output=<carpeta>` (defecto `APP_DOWNLOADS_PATH` del `.env` = `storage/downloads`) · `--filename=<nombre sin extensión>` |

`--url` acepta la URL completa (`watch?v=`, `youtu.be/`) o el id a secas. Salida: JSON por stdout (apto para `jq`),
progreso y errores por stderr. Exit 0 ok · 1 error de dominio · 2 inesperado (se loguea en `storage/logs`).

## 2 - cómo descarga (y por qué así)

1. Metadatos y formatos con `kkdai/youtube` (cliente ANDROID_VR).
2. `--format=video`: el mp4 **progresivo** (video+audio en un fichero, itag 18/22) con el cliente ANDROID. Con
   ANDROID_VR responde 403.
3. `--format=audio`: el mismo mp4 progresivo y se extrae la pista AAC con **ffmpeg** (`-acodec copy`, sin recodificar)
   a `.m4a`; el mp4 intermedio se borra. Requiere `ffmpeg` en el PATH.
4. No se usan las pistas **adaptativas** (audio solo, 1080p sin audio): YouTube las corta con 403 pasado ~1 MB en todos
   los clientes de la librería (verificado 2026-10-29). Por eso la calidad máxima real es la del progresivo (360p/720p).

## 3 - dependencias

- Go 1.26 · `github.com/kkdai/youtube/v2`
- `ffmpeg` solo para `--format=audio` (Windows: `winget install Gyan.FFmpeg`).
