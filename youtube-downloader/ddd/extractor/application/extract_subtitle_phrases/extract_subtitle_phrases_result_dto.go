package extractsubtitlephrases

// ExtractSubtitlePhrasesResultDto is returned by the use case service.
type ExtractSubtitlePhrasesResultDto struct {
	count   int
	phrases []map[string]any
}

// NewExtractSubtitlePhrasesResultDto builds the result from primitive phrase maps.
func NewExtractSubtitlePhrasesResultDto(count int, phrases []map[string]any) *ExtractSubtitlePhrasesResultDto {
	return &ExtractSubtitlePhrasesResultDto{count: count, phrases: phrases}
}

// Phrases returns the re-segmented phrases as primitive maps.
func (d *ExtractSubtitlePhrasesResultDto) Phrases() []map[string]any { return d.phrases }

// Count returns how many phrases were extracted.
func (d *ExtractSubtitlePhrasesResultDto) Count() int { return d.count }

// ToPrimitives serializes the result so a controller can print it.
func (d *ExtractSubtitlePhrasesResultDto) ToPrimitives() map[string]any {
	return map[string]any{
		"count":   d.count,
		"phrases": d.phrases,
	}
}
