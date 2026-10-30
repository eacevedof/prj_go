# yt-downloader

## 1 - instalación

### 1.1 Solo usarlo en Windows (sin instalar nada)

1. Descargar `yt-x.y.z.exe` (la `x.y.z` es la versión, p. ej. `yt-0.5.3.exe`) de la página de
   **Releases** del repo: <https://github.com/eacevedof/prj_go/releases> (el asset de la release `yt-downloader-vx.y.z`),
   y guardarlo en cualquier carpeta, p. ej. `C:\eaf-ytdl\`. Si no hay release, pedírselo a quien compila (§1.3).
2. Doble clic (o ejecutarlo sin argumentos): se abre un **menú** que pregunta qué hacer (video, info, playlist, audio,
   playlist en audio) y la URL; para un vídeo lee su info y lista las resoluciones reales para elegir. No hace falta `.env`: los vídeos caen en `downloads\` junto al `.exe` y los
   logs en `logs\`. La cabecera del menú muestra la versión y su fecha.
3. La primera descarga baja sola el motor **yt-dlp** (binario oficial, ~24 MB, verificado por SHA-256) a `tools\`
   junto al `.exe`, y una vez al día comprueba si hay versión nueva y la actualiza. Para vídeo por encima de 360p hace falta `ffmpeg` en el PATH (une vídeo y audio):
   `winget install Gyan.FFmpeg`. Sin ffmpeg: vídeo a 360p y audio `m4a` sí funcionan.

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
    cambio en `changelog.md`, `make dist` (genera `build/dist/yt-<versión>.exe`, Windows amd64, ~14 MB) y
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
del repo (`APP_DOWNLOADS_PATH`, §2.3). El nombre es siempre un **slug ASCII con fecha-hora**: el título "La Niña de
mis Ojos" (o el `--filename` que se pase) se guarda como `la-nina-de-mis-ojos-20261030-153045.mp4`; la carpeta de una
playlist es el slug de su título. Todo error queda además en `<downloads>/error-YYYYMMDD.log`. Con `--output` la
carpeta se crea si no existe. Sin ffmpeg: vídeo solo hasta 360p (progresivo) y audio `m4a` sí; con ffmpeg, hasta 1080p.

### 2.2 Comandos y opciones

| # | Comando | Opciones |
|---|---|---|
| 1 | `menu` | — · menú interactivo; es lo que arranca sin argumentos |
| 2 | `help` | — · lista los comandos; la cabecera lleva la versión |
| 3 | `version` | — · versión y fecha (`eaf-yt-dl v0.5.3 (2026-10-30)`) |
| 4 | `video-info --url=<url o id>` | — · devuelve JSON con título, duración y formatos |
| 5 | `video-download --url=<url o id>` | `--format=video\|audio` (defecto `video`) · `--quality=best\|hd1080\|hd720\|medium\|<N>p` (defecto `best` = la mayor que exista; `1440p`, `2160p`… también valen) · `--output=<carpeta>` (defecto `APP_DOWNLOADS_PATH`, §2.3; se crea si no existe) · `--filename=<nombre sin extensión>` (se pasa a slug + fecha-hora) |
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
| 7 | `APP_ENGINE` | `ytdlp` | motor de descarga: `ytdlp` (hasta 1080p) o `native` (solo mp4 progresivo 360p/720p, sin binarios externos) |

## 3 - cómo descarga (y por qué así)

1. Metadatos y formatos (título, duración, itags) siempre con la librería Go `kkdai/youtube` (cliente ANDROID_VR).
2. **Motor `ytdlp`** (defecto): la descarga la hace el binario oficial [yt-dlp](https://github.com/yt-dlp/yt-dlp).
   Si no está en `tools\`, junto al `.exe` ni en el PATH, se descarga de su última release de GitHub junto con el
   fichero `SHA2-256SUMS` y **solo se guarda si el SHA-256 coincide** con el publicado. Una vez al día (sello
   `tools\yt-dlp.last-check`) compara la versión instalada con la última release y, si cambió, la baja igual de
   verificada: YouTube rompe yt-dlp cada pocas semanas y sin esto el motor caduca. yt-dlp resuelve lo que la
   librería Go no puede: el reto JavaScript y el PO token que YouTube exige para las pistas **adaptativas**.
3. Qué son las adaptativas: YouTube publica cada vídeo como un mp4 **progresivo** (vídeo+audio ya mezclados, solo
   360p y a veces 720p: itag 18/22) y como pistas **adaptativas** separadas (vídeo solo 144p…4K, audio solo). El
   reproductor mezcla las adaptativas al vuelo; para tener 720p/1080p en un fichero hay que bajar las dos y unirlas.
4. `--format=video` con ffmpeg: mejor vídeo H.264 hasta la altura de `--quality` (`best` sin tope · `hd1080` ·
   `hd720` · `medium` 360p · `<N>p`) + mejor audio m4a, unidos en mp4 por ffmpeg. Por encima de 1080p YouTube solo
   tiene VP9/AV1: se bajan igual y se meten en el mp4. Sin ffmpeg: solo el progresivo (360p/720p) y se avisa.
5. `--format=audio`: la mejor pista de audio m4a (AAC) directamente; no hace falta ffmpeg.
6. **Motor `native`** (`APP_ENGINE=native`, o si yt-dlp no se pudo obtener): solo el mp4 progresivo con el cliente
   ANDROID; las adaptativas dan 403 pasado ~1 MB en todos los clientes de la librería (verificado 2026-10-29). Audio
   = progresivo + extracción AAC con ffmpeg. Hay progresivos **sin `content_length`** que YouTube recorta a ~200 KB:
   se detecta comparando con `bitrate × duración` y falla con `was cut by YouTube`.
7. `playlist-download` lee la lista (id + título) con la librería Go y delega cada vídeo en `video-download`,
   secuencialmente, con el motor que toque.

## 4 - dependencias

- Go 1.26 · `github.com/kkdai/youtube/v2` (metadatos y motor `native`).
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) (motor `ytdlp`): se descarga solo a `tools\` (gitignored), verificado por SHA-256.
- `ffmpeg` para vídeo por encima de 360p con `ytdlp` y para audio con `native` (Windows: `winget install Gyan.FFmpeg`).
