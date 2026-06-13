// Package services holds the extractor module's pure domain logic (no I/O).
package services

import (
	"sort"
	"strconv"
	"strings"

	"youtube-downloader/ddd/extractor/domain/entities"
	"youtube-downloader/ddd/extractor/domain/enums"
)

// PhraseSegmenter turns a raw word stream into clean phrases. It normalizes the
// stream (dedupe of rolling duplicates + sort by onset) and splits it on
// sentence-ending punctuation, giving precise per-phrase start/end times.
type PhraseSegmenter struct{}

// NewPhraseSegmenter returns a ready-to-use segmenter.
func NewPhraseSegmenter() *PhraseSegmenter { return &PhraseSegmenter{} }

// Segment normalizes words and splits them into phrases. finalEndMs closes the
// last phrase (and any phrase whose next-word timestamp is missing).
func (s *PhraseSegmenter) Segment(words []entities.Word, finalEndMs int) []entities.Phrase {
	words = s.dedupe(words)
	sort.SliceStable(words, func(i, j int) bool { return words[i].StartMs < words[j].StartMs })

	var phrases []entities.Phrase
	var cur []entities.Word

	flush := func(nextStart int) {
		if len(cur) == 0 {
			return
		}
		parts := make([]string, len(cur))
		for i, w := range cur {
			parts[i] = w.Text
		}
		start := cur[0].StartMs
		end := nextStart
		if end <= start {
			end = cur[len(cur)-1].StartMs + enums.PhraseFallbackPaddingMs
		}
		phrases = append(phrases, entities.Phrase{
			Index:   len(phrases) + 1,
			Text:    strings.Join(parts, " "),
			StartMs: start,
			EndMs:   end,
		})
		cur = nil
	}

	for i, w := range words {
		cur = append(cur, w)
		if s.endsSentence(w.Text) {
			next := finalEndMs
			if i+1 < len(words) {
				next = words[i+1].StartMs
			}
			flush(next)
		}
	}
	flush(finalEndMs)
	return phrases
}

// dedupe removes exact (onset, text) duplicates produced by the rolling format.
func (s *PhraseSegmenter) dedupe(words []entities.Word) []entities.Word {
	seen := make(map[string]bool, len(words))
	out := words[:0:0]
	for _, w := range words {
		key := strconv.Itoa(w.StartMs) + "\x00" + w.Text
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, w)
	}
	return out
}

// endsSentence reports whether a token ends a sentence (terminator possibly
// wrapped in closing quotes/brackets).
func (s *PhraseSegmenter) endsSentence(tok string) bool {
	tok = strings.TrimRight(tok, enums.PhraseTrailingTrimChars)
	if tok == "" {
		return false
	}
	return strings.IndexByte(enums.SentenceTerminators, tok[len(tok)-1]) >= 0
}
