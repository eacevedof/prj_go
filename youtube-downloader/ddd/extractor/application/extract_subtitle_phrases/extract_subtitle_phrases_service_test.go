package extractsubtitlephrases

import (
	"os"
	"path/filepath"
	"testing"
)

// TestExtractSubtitlePhrasesServiceReal runs the use case against the real probe
// caption (if present) and sanity-checks the re-segmentation. Run with -v to
// eyeball the first phrases.
func TestExtractSubtitlePhrasesServiceReal(t *testing.T) {
	path := filepath.Join("..", "..", "..", "..", "_probe", "-MsBisZFGn4.nl.vtt")
	if _, err := os.Stat(path); err != nil {
		t.Skipf("probe vtt not present: %v", err)
	}

	resultDto, err := GetInstance().Invoke(
		FromPrimitives(map[string]string{"subtitle_path": path, "lang": "nl"}),
	)
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	if resultDto.Count() == 0 {
		t.Fatal("expected phrases, got none")
	}

	t.Logf("parsed %d phrases", resultDto.Count())
	phrases := resultDto.Phrases()
	for i := 0; i < 10 && i < len(phrases); i++ {
		p := phrases[i]
		idx, _ := p["index"].(int)
		text, _ := p["text"].(string)
		startMs, _ := p["start_ms"].(int)
		endMs, _ := p["end_ms"].(int)
		t.Logf("%3d [%7.2fs -> %7.2fs] %s", idx, float64(startMs)/1000, float64(endMs)/1000, text)
	}

	for _, p := range phrases {
		startMs, _ := p["start_ms"].(int)
		endMs, _ := p["end_ms"].(int)
		idx, _ := p["index"].(int)
		text, _ := p["text"].(string)
		if endMs <= startMs {
			t.Errorf("phrase %d has non-positive duration: %d..%d", idx, startMs, endMs)
		}
		if text == "" {
			t.Errorf("phrase %d has empty text", idx)
		}
	}
}
