package extractsubtitlephrases

import "youtube-downloader/ddd/extractor/domain/entities"

// ExtractSubtitlePhrasesResultDto is returned by the use case service.
type ExtractSubtitlePhrasesResultDto struct {
	phrases []entities.Phrase
}

// NewExtractSubtitlePhrasesResultDto builds the result from segmented phrases.
func NewExtractSubtitlePhrasesResultDto(phrases []entities.Phrase) *ExtractSubtitlePhrasesResultDto {
	return &ExtractSubtitlePhrasesResultDto{phrases: phrases}
}

// Phrases returns the re-segmented phrases.
func (d *ExtractSubtitlePhrasesResultDto) Phrases() []entities.Phrase { return d.phrases }

// Count returns how many phrases were extracted.
func (d *ExtractSubtitlePhrasesResultDto) Count() int { return len(d.phrases) }

// ToPrimitives serializes the result so a controller can print it.
func (d *ExtractSubtitlePhrasesResultDto) ToPrimitives() map[string]any {
	return map[string]any{
		"count":   len(d.phrases),
		"phrases": d.phrases,
	}
}
