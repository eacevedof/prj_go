package boot

// ENV catalogo de entornos. El COMPORTAMIENTO vive en EnvironmentReaderRawRepository.
var ENV = struct {
	LOCAL      string
	DEVELOP    string
	TEST       string
	PRODUCTION string
}{
	LOCAL:      "local",
	DEVELOP:    "develop",
	TEST:       "test",
	PRODUCTION: "production",
}
