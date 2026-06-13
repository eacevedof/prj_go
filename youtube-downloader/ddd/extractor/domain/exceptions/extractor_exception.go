// Package exceptions holds the extractor module's controlled domain errors.
package exceptions

// ExtractorException is a controlled domain error. The CLI controller maps it
// to exit code 1 (a web controller would map it to 4xx).
type ExtractorException struct {
	message string
}

// NewExtractorException builds a domain error with a human-readable message.
func NewExtractorException(message string) *ExtractorException {
	return &ExtractorException{message: message}
}

func (e *ExtractorException) Error() string { return e.message }
