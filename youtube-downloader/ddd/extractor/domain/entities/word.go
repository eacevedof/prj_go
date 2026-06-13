package entities

// Word is a single spoken token with its onset timestamp in milliseconds,
// reconstructed from a subtitle's inline word-level timestamps.
type Word struct {
	Text    string
	StartMs int
}
