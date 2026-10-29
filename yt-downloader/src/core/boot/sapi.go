package boot

// Equivalente a php_sapi_name(). Este repo solo tiene consola, pero se conserva el contrato de task-go.
const (
	SAPI_SERVER = "server"
	SAPI_CLI    = "cli"
)

var sapiName = SAPI_SERVER

func SetSapiNameCli() {
	sapiName = SAPI_CLI
}

func GetSapiName() string {
	return sapiName
}

func IsCli() bool {
	return sapiName == SAPI_CLI
}
