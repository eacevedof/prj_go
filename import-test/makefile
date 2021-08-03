compile:
	cd hello; go mod edit -replace eaf/greetings=../greetings; go mod tidy; go run hello.go

init:
	cd greetings; rm go.mod; go mod init eaf/greetings; go mod tidy;
	cd hello; rm go.mod; go mod init eaf/hello; go mod tidy
