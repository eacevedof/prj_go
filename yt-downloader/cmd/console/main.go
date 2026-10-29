// Consola de yt-downloader-go (mismo esqueleto que task-go). go run ./cmd/console help
package main

import (
	"fmt"
	"os"

	"yt-downloader-go/console"
	"yt-downloader-go/src/core/config"
)

func main() {
	if err := config.AppbootLoadenv(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	ConsoleMainGetInstance(os.Args, console.GetCommands()).Exec()
}
