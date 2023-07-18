package env

const (
	PROJECT_NAME    = "debug.kitsune"
	PROJECT_VERSION = "local build"
)

var (
	BUILD_COMMIT             string
	BUILD_VERSION            string = PROJECT_VERSION
	BUILD_APP_NAME           string = PROJECT_NAME
	BUILD_SERVICE_SUPERVISOR string
	BUILD_SERVICE_GATEWAY    string
	BUILD_SERVICE_CACHE      string
	BUILD_SERVICE            = map[string]string{
		"supervisor": BUILD_SERVICE_SUPERVISOR,
		"gateway":    BUILD_SERVICE_GATEWAY,
		"cache":      BUILD_SERVICE_CACHE,
	}
)
