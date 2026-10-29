@echo off
setlocal
rem yt-downloader — instalador para Windows (sin make). Ejecutar desde la raiz del repo:  install.bat
rem Instala Go y ffmpeg con winget si faltan, descarga los modulos, crea .env y compila build\console.exe
cd /d "%~dp0"

echo [1/5] Go
where go >nul 2>&1
if errorlevel 1 (
    echo   no encontrado: instalando con winget...
    winget install --id GoLang.Go -e --accept-source-agreements --accept-package-agreements
    if errorlevel 1 goto :fail
    echo   Go instalado. ABRE UNA CONSOLA NUEVA para que coja el PATH y vuelve a lanzar install.bat
    goto :end
)
for /f "tokens=3" %%v in ('go version') do echo   %%v

echo [2/5] ffmpeg (solo para --format=audio)
where ffmpeg >nul 2>&1
if errorlevel 1 (
    echo   no encontrado: instalando con winget...
    winget install --id Gyan.FFmpeg -e --accept-source-agreements --accept-package-agreements
    if errorlevel 1 echo   AVISO: ffmpeg no se pudo instalar; el video funcionara, el audio no.
) else (
    echo   ok
)

echo [3/5] modulos Go
go mod download
if errorlevel 1 goto :fail

echo [4/5] .env
if not exist .env (
    copy /y .env.example .env >nul
    echo   creado desde .env.example
) else (
    echo   ya existe, no se toca
)

echo [5/5] compilando build\console.exe
go build -o build\console.exe .\cmd\console
if errorlevel 1 goto :fail

echo.
echo Listo. Prueba:
echo   build\console.exe help
echo   build\console.exe video-info --url=https://www.youtube.com/watch?v=dQw4w9WgXcQ
goto :end

:fail
echo.
echo ERROR: la instalacion ha fallado en el paso anterior.
exit /b 1

:end
endlocal
