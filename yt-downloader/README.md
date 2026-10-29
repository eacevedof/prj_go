# yt-downloader

Descargador de YouTube por consola en Go, con la arquitectura DDD de `task-go` (prj-task-cruds): `cmd/console` →
`console/commands.go` → `src/modules/videos/{application,domain,infrastructure}`. Sin API ni BD.
Sustituye al PoC `youtube-downloader/` (legacy).

## 1 - instalación

### 1.1 Windows (sin `make`)

1. Clonar el repo y abrir una consola (`cmd` o PowerShell) en `yt-downloader\`.
2. Ejecutar `.\install.bat` (en `cmd` y en PowerShell; el `.\` es obligatorio en PowerShell). Hace, en orden: instala **Go** con `winget` si falta · instala **ffmpeg** con `winget` si
   falta (solo hace falta para `--format=audio`) · `go mod download` · crea `.env` desde `.env.example` si no existe ·
   compila `build\console.exe`.
3. Si el `.bat` ha tenido que instalar Go, **cerrar la consola, abrir una nueva** (para que coja el PATH) y volver a
   ejecutar `.\install.bat`.
4. Comprobar: `.\build\console.exe help` debe listar los comandos.
5. (Opcional) Revisar `.env` — ver §2.3; con los valores por defecto ya funciona.

Si no se quiere usar el `.bat`, los mismos pasos a mano:

```
winget install --id GoLang.Go -e
winget install --id Gyan.FFmpeg -e
go mod download
copy .env.example .env
go build -o build\console.exe .\cmd\console
```

### 1.2 Linux / WSL / macOS

6. Instalar Go ≥ 1.26 (`https://go.dev/dl/`) y `ffmpeg` (`apt install ffmpeg` · `brew install ffmpeg`).
7. Clonar el repo y entrar en `yt-downloader/`.
8. Con `make`: `make build` (crea `.env` y compila `build/console`). Sin `make`:
   `cp .env.example .env && go mod download && go build -o build/console ./cmd/console`.
9. Comprobar: `build/console help` (o `make console c="help"`).

## 2 - uso

### 2.1 Ejemplos

Windows (`cmd` / PowerShell, con `.\` delante; en **Git Bash** usar `./build/console.exe`):

```
.\build\console.exe help
.\build\console.exe video-info --url=https://www.youtube.com/watch?v=dQw4w9WgXcQ
.\build\console.exe video-download --url=dQw4w9WgXcQ --format=video --quality=best
.\build\console.exe video-download --url=dQw4w9WgXcQ --format=video --output=..\mis-videos
.\build\console.exe video-download --url=https://youtu.be/dQw4w9WgXcQ --format=audio --filename=cancion
go test -count=1 ./...
```

Linux / WSL / macOS (con `make`; sin él, `build/console <comando>` directamente):

```
make console c="help"
make console c="video-info --url=https://www.youtube.com/watch?v=dQw4w9WgXcQ"
make console c="video-download --url=dQw4w9WgXcQ --format=video --quality=best"
make console c="video-download --url=dQw4w9WgXcQ --format=video --output=../mis-videos"
make console c="video-download --url=https://youtu.be/dQw4w9WgXcQ --format=audio --filename=cancion"
make test
```

Sin `--output` el fichero cae en `storage/downloads` del repo (`APP_DOWNLOADS_PATH`, §2.3). Con `--output` la
carpeta se crea si no existe. Sin ffmpeg el `--format=video` funciona; el `--format=audio` no.

### 2.2 Comandos y opciones

| # | Comando | Opciones |
|---|---|---|
| 1 | `help` | — |
| 2 | `video-info --url=<url o id>` | — · devuelve JSON con título, duración y formatos |
| 3 | `video-download --url=<url o id>` | `--format=video\|audio` (defecto `video`) · `--quality=best\|hd720\|medium` (defecto `best`) · `--output=<carpeta>` (defecto `APP_DOWNLOADS_PATH` del `.env` = `storage/downloads`; se crea si no existe) · `--filename=<nombre sin extensión>` |

`--url` acepta la URL completa (`watch?v=`, `youtu.be/`) o el id a secas. Salida: JSON por stdout (apto para `jq`),
progreso y errores por stderr. Exit 0 ok · 1 error de dominio · 2 inesperado (se loguea en `storage/logs`).

### 2.3 Variables de entorno (`.env`)

Se leen de `<raíz>/.env` (lo crea `install.bat` / `make env` desde `.env.example`, ver §1); una variable **real del
sistema** con el mismo nombre gana sobre el fichero. `%PATH_ROOT%` se sustituye por la raíz del repo. Si falta
una obligatoria el programa no arranca y lo dice (`required env keys not set`).

| # | Variable | Obligatoria | Defecto en `.env.example` | Para qué |
|---|---|---|---|---|
| 1 | `APP_ENV` | **sí** | `local` | entorno (`local`, `prod`) |
| 2 | `APP_DOWNLOADS_PATH` | **sí** | `%PATH_ROOT%/storage/downloads` | carpeta de descargas por defecto; `--output` la sobrescribe |
| 3 | `APP_DEBUG` | no | `0` | `1` muestra el stack de los errores inesperados |
| 4 | `APP_NAME` | no | `yt-downloader-go` | nombre en logs |
| 5 | `APP_LOG_PATH` | no | `%PATH_ROOT%/storage/logs` | carpeta de los logs `YYYY-MM-DD-<tipo>.log` |

## 3 - cómo descarga (y por qué así)

1. Metadatos y formatos con `kkdai/youtube` (cliente ANDROID_VR).
2. `--format=video`: el mp4 **progresivo** (video+audio en un fichero, itag 18/22) con el cliente ANDROID. Con
   ANDROID_VR responde 403.
3. `--format=audio`: el mismo mp4 progresivo y se extrae la pista AAC con **ffmpeg** (`-acodec copy`, sin recodificar)
   a `.m4a`; el mp4 intermedio se borra. Requiere `ffmpeg` en el PATH.
4. No se usan las pistas **adaptativas** (audio solo, 1080p sin audio): YouTube las corta con 403 pasado ~1 MB en todos
   los clientes de la librería (verificado 2026-10-29). Por eso la calidad máxima real es la del progresivo (360p/720p).

## 4 - dependencias

- Go 1.26 · `github.com/kkdai/youtube/v2`
- `ffmpeg` solo para `--format=audio` (Windows: `winget install Gyan.FFmpeg`).
