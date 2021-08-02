#### Fazt - Rest API con go

- servidor:
  - módulo go mux
  - go get -u github.com/gorilla/mux
    - descarga y permite importarlo

- refrescador de servidor (hot-reload):
  - go get github.com/githubnemo/CompileDaemon

#### errores
#### 1
- go run main.go 
- main.go:4:2: package mux is not in GOROOT (/usr/local/Cellar/go/1.16.6/libexec/src/mux)
- main.go:6:2: no required module provides package github.com/gorilla/mux: go.mod file not found in current directory or any parent directory; see 'go help modules'
  - habia que iniciar el proyecto con `go mod init go-api`
#### 2
- peticiones entrantes en mac. El servidor tiene que ir con localhost
#### 3
- no me imprime el json
