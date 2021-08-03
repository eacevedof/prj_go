compile:
	go run main.go
server:
	CompileDaemon -command="./go-api"

init:
	# rm go.mod; go mod init eaf
	cd services/crud; rm go.mod; go mod init eaf/services/crud
	go mod tidy