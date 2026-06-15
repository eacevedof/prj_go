package components

import (
	"regexp"
	"strings"
)

var sluggerRexNonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

// Slugger converts text to URL-safe hyphenated slugs.
type Slugger struct{}

// NewSlugger returns a ready-to-use Slugger.
func NewSlugger() *Slugger { return &Slugger{} }

// GetSluggedText converts text to a lowercase hyphenated slug.
func (s *Slugger) GetSluggedText(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))
	text = strings.NewReplacer(
		"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u", "ñ", "n",
		"Á", "a", "É", "e", "Í", "i", "Ó", "o", "Ú", "u", "Ü", "u", "Ñ", "n",
	).Replace(text)
	return strings.Trim(sluggerRexNonAlnum.ReplaceAllString(text, "-"), "-")
}
