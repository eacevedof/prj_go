// Package extractsubtitlephrases is the "extract subtitle phrases" use case:
// read a subtitle file and re-segment it into timed phrases.
package extractsubtitlephrases

// ExtractSubtitlePhrasesDto is the immutable input DTO for the use case.
type ExtractSubtitlePhrasesDto struct {
	subtitlePath string
	lang         string
}

// FromPrimitives builds the input DTO from a primitives map (CLI args / request).
func FromPrimitives(primitives map[string]string) *ExtractSubtitlePhrasesDto {
	return &ExtractSubtitlePhrasesDto{
		subtitlePath: primitives["subtitle_path"],
		lang:         primitives["lang"],
	}
}

// SubtitlePath is the path to the subtitle (.vtt) file to read.
func (d *ExtractSubtitlePhrasesDto) SubtitlePath() string { return d.subtitlePath }

// Lang is the subtitle language code (e.g. "nl").
func (d *ExtractSubtitlePhrasesDto) Lang() string { return d.lang }
