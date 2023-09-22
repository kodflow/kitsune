package config

const (
	DEFAULT_APP_NAME = "kitsune"
	DEFAULT_VERSION  = "local"
)

var (
	BUILD_COMMIT   string
	BUILD_VERSION  string = DEFAULT_VERSION
	BUILD_APP_NAME string = DEFAULT_APP_NAME

	BUILD_ARCH string
	BUILD_OS   string

	BUILD_SERVICE_SUPERVISOR string
	BUILD_SERVICE_GATEWAY    string
	BUILD_SERVICE_CACHE      string

	BUILD_SERVICES = map[string]string{
		"supervisor": BUILD_SERVICE_SUPERVISOR,
		"gateway":    BUILD_SERVICE_GATEWAY,
		"cache":      BUILD_SERVICE_CACHE,
	}
)
