#compile:
	# go build -buildmode=plugin eaf/main.go => package eaf/main.go is not in GOROOT (/usr/local/Cellar/go/1.16.6/libexec/src/eaf/main.go)
	# go run main.go #cannot find module for path _/Users/<user>/projects/prj_go/fazt-golang-restapi/crud
	# GO111MODULE=off go run main.go # no me vale pq no usa los módulos instalados por go  cannot find package "github.com/gorilla/mux" in any of:

compile:
	go run main.go

server:
	CompileDaemon -command="./eaf"

init-main:
	# rm -f go.mod; no lo borro pq tiene los imports manuales
	rm -f go.sum; rm -f eaf
	go mod init eaf; go mod tidy

init-crud:
	cd crud; rm -f go.*; go mod init crud; go mod tidy