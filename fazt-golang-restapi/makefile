compile:
	go run main.go

server:
	CompileDaemon -command="./eaf"

init:
	rm -f go.mod; rm -f go.sum; rm -f eaf
	go mod init eaf
	# cd services/crud; rm go.mod; go mod init eaf/services/crud
	go mod tidy