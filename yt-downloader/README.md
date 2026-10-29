# yt-downloader

Descargador de YouTube por consola en Go, con la arquitectura DDD de `task-go` (prj-task-cruds): `cmd/console` →
`console/commands.go` → `src/modules/videos/{application,domain,infrastructure}`. Sin API ni BD.
Sustituye al PoC `youtube-downloader/` (legacy).

## 1 - arrancar

```
make build                                   # .env desde .env.example + build/console
make console c="help"
make console c="video-info --url=https://www.youtube.com/watch?v=<id>"
make console c="video-download --url=<url o id> --format=video --quality=best --output=<carpeta>"
make console c="video-download --url=<url o id> --format=audio --filename=<nombre>"
make test
```

Salida: JSON por stdout (apto para `jq`), progreso y errores por stderr. Exit 0 ok · 1 error de dominio · 2 inesperado
(se loguea en `storage/logs`).

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
