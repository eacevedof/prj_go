package components

import (
	"encoding/base64"
	"encoding/json"
)

// Encoder handles JSON and base64 serialization.
type Encoder struct{}

// NewEncoder returns a ready-to-use Encoder.
func NewEncoder() *Encoder { return &Encoder{} }

func (e *Encoder) GetBase64Encoded(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

func (e *Encoder) GetBase64Decoded(encoded string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (e *Encoder) GetArrayAsBase64Encoded(data any) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (e *Encoder) GetDecodedArrayFromBase64(encoded string) (any, error) {
	b, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	var result any
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return result, nil
}
