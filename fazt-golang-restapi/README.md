#### Fazt - Rest API con go
- [Código original](https://github.com/FaztWeb/golang-restapi-crud/blob/master/main.go)

- servidor:
  - módulo go mux
  - go get -u github.com/gorilla/mux
    - descarga y permite importarlo

- refrescador de servidor (hot-reload):
  - go get github.com/githubnemo/CompileDaemon

#### errores
#### 1
```
could not import github.com/gorilla/mux (cannot find package "github.com/gorilla/mux" in any of 
	/usr/local/Cellar/go/1.16.6/libexec/src/github.com/gorilla/mux (from $GOROOT)
	/Users/<eaf>/go/src/github.com/gorilla/mux (from $GOPATH))
```
- go run main.go 
- main.go:4:2: package mux is not in GOROOT (/usr/local/Cellar/go/1.16.6/libexec/src/mux)
- main.go:6:2: no required module provides package github.com/gorilla/mux: go.mod file not found in current directory or any parent directory; see 'go help modules'
  - habia que iniciar el proyecto con `go mod init go-api`
#### 2
- peticiones entrantes en mac. El servidor tiene que ir con localhost
#### 3
- no me imprime el json
  - **struct field has json tag but is not exported** (mensaje en vscode)
  - los tags de la estructura deben empezar con mayusuculas, pq de lo contrario los entiende como campos privados

#### 4
- go: crud@v0.0.0 (replaced by ./crud): reading crud/go.mod: open /Users/<user>/projects/prj_go/fazt-golang-restapi/crud/go.mod: no such file or directory
  - Esto pasa porque la carpeta ./crud la entiende como un módulo y para esto debe contar con su fichero go.mod 

#### 5
- cannot refer to unexported name crud.SelectAll
  - el metodo debe empezar con mayusuclas para poder ser exportado 