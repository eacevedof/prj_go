compile:
	go run main.go

server:
	CompileDaemon -command="./crud"

init:
	rm go.mod; rm go.sum; rm -f crud
	go mod init crud
	go mod tidy