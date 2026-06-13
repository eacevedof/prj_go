// Package commands holds the extractor module's CLI controllers. A command has
// the same architecture as a web controller; only the IO boundary differs
// (console args in, stdout/exit code out).
package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	extractsubtitlephrases "youtube-downloader/ddd/extractor/application/extract_subtitle_phrases"
	"youtube-downloader/ddd/extractor/domain/exceptions"
)

// ExtractSubtitlePhrasesCommand is the CLI controller for the use case.
type ExtractSubtitlePhrasesCommand struct {
	extractSubtitlePhrasesService *extractsubtitlephrases.ExtractSubtitlePhrasesService
}

// NewExtractSubtitlePhrasesCommand wires the command to its service.
func NewExtractSubtitlePhrasesCommand() *ExtractSubtitlePhrasesCommand {
	return &ExtractSubtitlePhrasesCommand{
		extractSubtitlePhrasesService: extractsubtitlephrases.NewExtractSubtitlePhrasesService(),
	}
}

// Execute parses args → builds DTO → calls service → prints JSON → returns an
// exit code. It is the only layer that catches: domain errors → 1, others → 2.
func (c *ExtractSubtitlePhrasesCommand) Execute(args map[string]string) int {
	extractSubtitlePhrasesResultDto, err := c.extractSubtitlePhrasesService.Invoke(
		extractsubtitlephrases.FromPrimitives(args),
	)
	if err != nil {
		var extractorException *exceptions.ExtractorException
		if errors.As(err, &extractorException) {
			fmt.Fprintln(os.Stderr, "Error:", extractorException.Error())
			return 1
		}
		fmt.Fprintln(os.Stderr, "Internal error:", err)
		return 2
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(extractSubtitlePhrasesResultDto.ToPrimitives()); err != nil {
		fmt.Fprintln(os.Stderr, "Internal error:", err)
		return 2
	}
	return 0
}
