package enums

// PhraseFallbackPaddingMs closes a phrase whose computed end is not after its
// start (no following word to borrow a timestamp from):
// end = lastWordStartMs + PhraseFallbackPaddingMs.
const PhraseFallbackPaddingMs = 400

// SentenceTerminators are the characters that end a phrase during
// re-segmentation of the word stream.
const SentenceTerminators = ".!?"

// PhraseTrailingTrimChars are closing quotes/brackets stripped from a token's
// tail before checking whether it ends a sentence.
const PhraseTrailingTrimChars = `"')]}»`
