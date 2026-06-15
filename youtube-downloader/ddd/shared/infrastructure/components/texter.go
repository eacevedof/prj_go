package components

import (
	"regexp"
	"strings"
)

var (
	texterRexSpaces   = regexp.MustCompile(`\s+`)
	texterRexNonAlpha = regexp.MustCompile(`[^a-zA-Z0-9áéíóúüñÁÉÍÓÚÜÑ@.,_\- \n#|+/:]`)
)

// Texter provides text sanitization utilities.
type Texter struct{}

// NewTexter returns a ready-to-use Texter.
func NewTexter() *Texter { return &Texter{} }

// GetSanitizedPrimitives sanitizes all string values in a map.
func (t *Texter) GetSanitizedPrimitives(data map[string]string) map[string]string {
	result := make(map[string]string, len(data))
	for k, v := range data {
		result[k] = t.GetSanitizedText(v)
	}
	return result
}

// GetTrimLoweredText trims, collapses whitespace, and lowercases text.
func (t *Texter) GetTrimLoweredText(text string) string {
	return strings.ToLower(strings.TrimSpace(texterRexSpaces.ReplaceAllString(text, " ")))
}

// GetSanitizedText strips unsafe characters, retaining safe ones only.
func (t *Texter) GetSanitizedText(text string) string {
	text = texterRexSpaces.ReplaceAllString(strings.TrimSpace(text), " ")
	return texterRexNonAlpha.ReplaceAllString(text, "")
}
