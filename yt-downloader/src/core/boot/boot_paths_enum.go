package boot

import (
	"os"
	"path/filepath"
)

// BootPathsEnum rutas absolutas de la app. La raiz es la primera carpeta con go.mod subiendo desde el cwd.
var BootPathsEnum = func() struct {
	PATH_ROOT      string
	PATH_SRC       string
	PATH_LOGS      string
	PATH_DOWNLOADS string
} {
	pathRoot := getPathRoot()
	return struct {
		PATH_ROOT      string
		PATH_SRC       string
		PATH_LOGS      string
		PATH_DOWNLOADS string
	}{
		PATH_ROOT:      pathRoot,
		PATH_SRC:       filepath.Join(pathRoot, "src"),
		PATH_LOGS:      filepath.Join(pathRoot, "storage", "logs"),
		PATH_DOWNLOADS: filepath.Join(pathRoot, "storage", "downloads"),
	}
}()

func getPathRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	for dir := cwd; ; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		if filepath.Dir(dir) == dir {
			return cwd
		}
	}
}
