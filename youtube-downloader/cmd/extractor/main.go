// Command extractor is the CLI entry point: it parses flags and delegates to the
// extractor module's command controller.
package main

import (
	"flag"
	"os"

	"youtube-downloader/ddd/extractor/infrastructure/commands"
)

func main() {
	subtitlePath := flag.String("vtt", "", "path to the .vtt subtitle file")
	lang := flag.String("lang", "nl", "subtitle language code")
	flag.Parse()

	os.Exit(commands.NewExtractSubtitlePhrasesCommand().Execute(map[string]string{
		"subtitle_path": *subtitlePath,
		"lang":          *lang,
	}))
}
