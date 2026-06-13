// Package repositories accesses the extractor module's data sources.
package repositories

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"

	"youtube-downloader/ddd/extractor/domain/entities"
	"youtube-downloader/ddd/extractor/domain/enums"
)

// YouTube auto-captions use a "rolling" format: each spoken word carries an
// inline timestamp (<00:00:00.599>) inside cues whose own start/end times are
// line-wrap boundaries, not sentence boundaries, and the same text is re-emitted
// as plain carry-over lines. This repository ignores the misleading cue
// boundaries and reconstructs the raw word stream from word-level timestamps;
// re-segmentation into phrases is the domain's job.
var (
	reCue   = regexp.MustCompile(`^(\d{2}):(\d{2}):(\d{2})\.(\d{3})\s+-->\s+(\d{2}):(\d{2}):(\d{2})\.(\d{3})`)
	reTS    = regexp.MustCompile(`<(\d{2}):(\d{2}):(\d{2})\.(\d{3})>`)
	reCtags = regexp.MustCompile(`</?c[^>]*>`)
)

// SubtitlesReaderVttRepository reads YouTube auto-caption WebVTT files.
type SubtitlesReaderVttRepository struct{}

// NewSubtitlesReaderVttRepository returns a ready-to-use reader.
func NewSubtitlesReaderVttRepository() *SubtitlesReaderVttRepository {
	return &SubtitlesReaderVttRepository{}
}

// GetWordStream reads a .vtt file and returns its raw word stream with inline
// timestamps, plus the maximum cue end time. Only markup lines (with inline word
// timestamps) contribute words; plain carry-over lines are skipped. No exception
// handling — I/O errors propagate untouched to the caller.
func (r *SubtitlesReaderVttRepository) GetWordStream(path string) (*entities.WordStream, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, enums.VttScannerBufferBytes), enums.VttScannerBufferBytes)

	var (
		curStart int
		maxEnd   int
		words    []entities.Word
	)
	for sc.Scan() {
		line := sc.Text()
		if m := reCue.FindStringSubmatch(line); m != nil {
			curStart = hmsToMs(m[1], m[2], m[3], m[4])
			if e := hmsToMs(m[5], m[6], m[7], m[8]); e > maxEnd {
				maxEnd = e
			}
			continue
		}
		if !reTS.MatchString(line) {
			continue
		}
		words = append(words, parseMarkupLine(line, curStart)...)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return &entities.WordStream{Words: words, FinalEndMs: maxEnd}, nil
}

// parseMarkupLine extracts word/timestamp pairs from a caption line such as
// "Hallo<00:00:00.599><c> en</c><00:00:00.840><c> welkom</c>". Text before the
// first inline timestamp inherits the cue start time; each subsequent word takes
// the timestamp immediately preceding it.
func parseMarkupLine(line string, cueStartMs int) []entities.Word {
	line = reCtags.ReplaceAllString(line, "")
	locs := reTS.FindAllStringIndex(line, -1)
	subs := reTS.FindAllStringSubmatch(line, -1)

	var out []entities.Word
	add := func(seg string, onset int) {
		for _, w := range strings.Fields(seg) {
			out = append(out, entities.Word{Text: w, StartMs: onset})
		}
	}

	if len(locs) == 0 {
		add(line, cueStartMs)
		return out
	}
	add(line[:locs[0][0]], cueStartMs)
	for i := range locs {
		onset := hmsToMs(subs[i][1], subs[i][2], subs[i][3], subs[i][4])
		end := len(line)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		add(line[locs[i][1]:end], onset)
	}
	return out
}

func hmsToMs(h, m, s, ms string) int {
	hi, _ := strconv.Atoi(h)
	mi, _ := strconv.Atoi(m)
	si, _ := strconv.Atoi(s)
	msi, _ := strconv.Atoi(ms)
	return ((hi*60+mi)*60+si)*1000 + msi
}
