compile:
	go run main.go
server:
	CompileDaemon -command="./go-api"

init:
	rm go.mod
	go mod init go-api
	go mod tidy