// Package enums holds domain-layer constants for the extractor module.
// No magic numbers/strings leak into services, repositories or commands.
package enums

// VttScannerBufferBytes is the bufio.Scanner buffer used to read VTT lines.
// Auto-caption lines with inline word timestamps can be very long, so 1 MB
// avoids "token too long" on dense cues.
const VttScannerBufferBytes = 1024 * 1024 // 1 MB
