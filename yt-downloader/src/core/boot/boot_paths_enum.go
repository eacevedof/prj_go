package boot

import (
	"os"
	"path/filepath"
)

// BootPathsEnum rutas absolutas de la app. Dos modos:
//   - REPO (hay go.mod subiendo desde el cwd): raiz = esa carpeta; descargas y logs en storage/.
//   - EXE EMPAQUETADO (no hay go.mod): raiz = carpeta del ejecutable; descargas en <exe>/downloads y logs en
//     <exe>/logs, para que quien solo tiene el .exe no configure nada.
//   - PATH_TOOLS (<raiz>/tools): binarios de terceros que la app descarga (yt-dlp). Gitignored.
var BootPathsEnum = func() struct {
	PATH_ROOT      string
	PATH_SRC       string
	PATH_LOGS      string
	PATH_DOWNLOADS string
	PATH_TOOLS     string
	IS_PACKAGED    bool
} {
	pathRoot, isPackaged := getPathRoot()
	pathLogs := filepath.Join(pathRoot, "storage", "logs")
	pathDownloads := filepath.Join(pathRoot, "storage", "downloads")
	if isPackaged {
		pathLogs = filepath.Join(pathRoot, "logs")
		pathDownloads = filepath.Join(pathRoot, "downloads")
	}
	return struct {
		PATH_ROOT      string
		PATH_SRC       string
		PATH_LOGS      string
		PATH_DOWNLOADS string
		PATH_TOOLS     string
		IS_PACKAGED    bool
	}{
		PATH_ROOT:      pathRoot,
		PATH_SRC:       filepath.Join(pathRoot, "src"),
		PATH_LOGS:      pathLogs,
		PATH_DOWNLOADS: pathDownloads,
		PATH_TOOLS:     filepath.Join(pathRoot, "tools"),
		IS_PACKAGED:    isPackaged,
	}
}()

// getPathRoot (raiz, empaquetado): la primera carpeta con go.mod subiendo desde el cwd; si no hay, la del ejecutable.
func getPathRoot() (string, bool) {
	if cwd, err := os.Getwd(); err == nil {
		for dir := cwd; ; dir = filepath.Dir(dir) {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				return dir, false
			}
			if filepath.Dir(dir) == dir {
				break
			}
		}
	}
	if executable, err := os.Executable(); err == nil {
		return filepath.Dir(executable), true
	}
	return ".", true
}
