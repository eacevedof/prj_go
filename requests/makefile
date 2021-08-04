compile:
	go run main.go

server:
	CompileDaemon -command="./eaf"

init-main:
	rm -f go.sum; rm -f eaf
	go mod init eaf; go mod tidy

init-request:
	cd request; rm -f go.*; go mod init request; go mod tidy