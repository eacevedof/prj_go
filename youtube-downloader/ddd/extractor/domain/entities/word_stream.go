package entities

// WordStream is the raw output of reading a subtitle file: the word sequence
// (with rolling duplicates still present) plus the maximum cue end time, used
// to close the final phrase during segmentation.
type WordStream struct {
	Words      []Word
	FinalEndMs int
}
