// Package entities holds the extractor module's domain models.
package entities

// Phrase is a re-segmented sentence with precise timing in milliseconds.
type Phrase struct {
	Index   int    `json:"index"`
	Text    string `json:"text"`
	StartMs int    `json:"start_ms"`
	EndMs   int    `json:"end_ms"`
}

// StartSec returns the start time in seconds (for ffmpeg -ss).
func (p Phrase) StartSec() float64 { return float64(p.StartMs) / 1000 }

// EndSec returns the end time in seconds (for ffmpeg -to).
func (p Phrase) EndSec() float64 { return float64(p.EndMs) / 1000 }
