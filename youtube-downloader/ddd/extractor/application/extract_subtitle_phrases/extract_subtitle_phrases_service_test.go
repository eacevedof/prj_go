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

	resultDto, err := NewExtractSubtitlePhrasesService().Invoke(
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
		t.Logf("%3d [%7.2fs -> %7.2fs] %s", p.Index, p.StartSec(), p.EndSec(), p.Text)
	}

	for _, p := range phrases {
		if p.EndMs <= p.StartMs {
			t.Errorf("phrase %d has non-positive duration: %d..%d", p.Index, p.StartMs, p.EndMs)
		}
		if p.Text == "" {
			t.Errorf("phrase %d has empty text", p.Index)
		}
	}
}
