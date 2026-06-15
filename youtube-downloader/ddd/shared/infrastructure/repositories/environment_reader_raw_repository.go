package repositories

import "os"

// EnvironmentReaderRawRepository reads environment variables.
type EnvironmentReaderRawRepository struct{}

// GetInstance returns a ready-to-use environment reader.
func GetInstance() *EnvironmentReaderRawRepository {
	return &EnvironmentReaderRawRepository{}
}

// Get returns the env var value or defaultVal if not set.
func (r *EnvironmentReaderRawRepository) Get(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// GetRequired returns the env var value or panics if not set.
func (r *EnvironmentReaderRawRepository) GetRequired(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("required env var not set: " + key)
	}
	return v
}
