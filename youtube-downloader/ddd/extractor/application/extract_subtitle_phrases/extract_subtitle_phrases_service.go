package extractsubtitlephrases

import (
	"os"
	"strings"

	"youtube-downloader/ddd/extractor/domain/exceptions"
	"youtube-downloader/ddd/extractor/domain/services"
	"youtube-downloader/ddd/extractor/infrastructure/repositories"
)

// ExtractSubtitlePhrasesService orchestrates the use case: validate input, read
// the subtitle word stream and re-segment it into phrases. One public method.
type ExtractSubtitlePhrasesService struct {
	subtitlesReaderVttRepository *repositories.SubtitlesReaderVttRepository
	phraseSegmenter              *services.PhraseSegmenter
}

// NewExtractSubtitlePhrasesService wires its collaborators. This single explicit
// place declaring dependencies makes a future move to a DI container mechanical.
func NewExtractSubtitlePhrasesService() *ExtractSubtitlePhrasesService {
	return &ExtractSubtitlePhrasesService{
		subtitlesReaderVttRepository: repositories.NewSubtitlesReaderVttRepository(),
		phraseSegmenter:              services.NewPhraseSegmenter(),
	}
}

// Invoke runs the use case. Domain problems return an ExtractorException;
// infrastructure errors (I/O) propagate unchanged for the controller to map.
func (s *ExtractSubtitlePhrasesService) Invoke(extractSubtitlePhrasesDto *ExtractSubtitlePhrasesDto) (*ExtractSubtitlePhrasesResultDto, error) {
	if err := s.failIfWrongInput(extractSubtitlePhrasesDto); err != nil {
		return nil, err
	}

	wordStream, err := s.subtitlesReaderVttRepository.GetWordStream(extractSubtitlePhrasesDto.SubtitlePath())
	if err != nil {
		return nil, err
	}

	return NewExtractSubtitlePhrasesResultDto(
		s.phraseSegmenter.Segment(wordStream.Words, wordStream.FinalEndMs),
	), nil
}

// failIfWrongInput validates the input before any repository is touched.
func (s *ExtractSubtitlePhrasesService) failIfWrongInput(extractSubtitlePhrasesDto *ExtractSubtitlePhrasesDto) error {
	if strings.TrimSpace(extractSubtitlePhrasesDto.SubtitlePath()) == "" {
		return exceptions.NewExtractorException("subtitle path is required")
	}
	if _, err := os.Stat(extractSubtitlePhrasesDto.SubtitlePath()); err != nil {
		return exceptions.NewExtractorException("subtitle file not found: " + extractSubtitlePhrasesDto.SubtitlePath())
	}
	return nil
}
