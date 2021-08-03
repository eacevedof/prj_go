- [call-module](https://golang.org/doc/tutorial/call-module-code)

#### 1
```
cd hello; go run hello.go
hello.go:6:2: package eaf/greetings is not in GOROOT (/usr/local/Cellar/go/1.16.6/libexec/src/eaf/greetings)
make: *** [compile] Error 1
```
- cd hello; go mod edit -replace eaf/greetings=../greetings; go mod tidy; go run hello.go