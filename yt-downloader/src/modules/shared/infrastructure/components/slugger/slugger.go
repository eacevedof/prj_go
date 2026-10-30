package slugger

import (
	"regexp"
	"strings"
	"time"
)

const SLUG_MAX_LENGTH = 100

var (
	nonSlugChars    = regexp.MustCompile(`[^a-z0-9]+`)
	repeatedDashes  = regexp.MustCompile(`-{2,}`)
	latinToAsciiMap = map[rune]string{
		'á': "a", 'à': "a", 'ä': "a", 'â': "a", 'ã': "a", 'å': "a", 'æ': "ae",
		'é': "e", 'è': "e", 'ë': "e", 'ê': "e",
		'í': "i", 'ì': "i", 'ï': "i", 'î': "i",
		'ó': "o", 'ò': "o", 'ö': "o", 'ô': "o", 'õ': "o", 'ø': "o", 'œ': "oe",
		'ú': "u", 'ù': "u", 'ü': "u", 'û': "u",
		'ñ': "n", 'ç': "c", 'ß': "ss", 'ý': "y", 'ÿ': "y",
		'Á': "a", 'À': "a", 'Ä': "a", 'Â': "a", 'Ã': "a", 'Å': "a", 'Æ': "ae",
		'É': "e", 'È': "e", 'Ë': "e", 'Ê': "e",
		'Í': "i", 'Ì': "i", 'Ï': "i", 'Î': "i",
		'Ó': "o", 'Ò': "o", 'Ö': "o", 'Ô': "o", 'Õ': "o", 'Ø': "o", 'Œ': "oe",
		'Ú': "u", 'Ù': "u", 'Ü': "u", 'Û': "u",
		'Ñ': "n", 'Ç': "c", 'Ý': "y",
		'&': " and ", '@': " at ", '+': " plus ",
	}
)

// Slugger nombres de fichero/carpeta en slug ASCII: "La Niña de mis Ojos" -> "la-nina-de-mis-ojos". Acentos y eñes se
// transliteran; lo que no sea latino (cirilico, kanji, emojis) se descarta. Nunca devuelve vacio: fallback.
type Slugger struct{}

func GetInstance() *Slugger {
	return &Slugger{}
}

func (s *Slugger) GetSlug(text string, fallback string) string {
	var ascii strings.Builder
	for _, char := range text {
		if replacement, isMapped := latinToAsciiMap[char]; isMapped {
			ascii.WriteString(replacement)
			continue
		}
		if char < 128 {
			ascii.WriteRune(char)
		}
	}
	slug := strings.ToLower(ascii.String())
	slug = nonSlugChars.ReplaceAllString(slug, "-")
	slug = repeatedDashes.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > SLUG_MAX_LENGTH {
		slug = strings.Trim(slug[:SLUG_MAX_LENGTH], "-")
	}
	if slug == "" {
		return s.GetSlug(fallback, "file")
	}
	return slug
}

// GetSlugWithNow "La Niña de mis Ojos" -> "la-nina-de-mis-ojos-20261030-153045" (hora local): nombres unicos por descarga.
func (s *Slugger) GetSlugWithNow(text string, fallback string) string {
	return s.GetSlug(text, fallback) + "-" + time.Now().Format("20060102-150405")
}
