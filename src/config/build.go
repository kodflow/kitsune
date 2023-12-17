package config

// Constants for default application name and version.
const (
	DEFAULT_APP_NAME = "kitsune" // The default name of the application.
	DEFAULT_VERSION  = "local"   // The default version, typically used for local or development builds.
)

// Global variables for build information.
var (
	BUILD_COMMIT   string                    // The commit hash of the build, useful for tracking specific builds in version control.
	BUILD_VERSION  string = DEFAULT_VERSION  // The version of the build, defaults to the value in DEFAULT_VERSION.
	BUILD_APP_NAME string = DEFAULT_APP_NAME // The name of the build, defaults to the value in DEFAULT_APP_NAME.

	BUILD_ARCH string // The architecture for which the build was compiled (e.g., amd64, arm).
	BUILD_OS   string // The operating system for which the build was compiled (e.g., linux, windows).

	BUILD_SERVICE_SUPERVISOR string // Build identifier for the supervisor service.
	BUILD_SERVICE_GATEWAY    string // Build identifier for the gateway service.
	BUILD_SERVICE_CACHE      string // Build identifier for the cache service.

	// BUILD_SERVICES is a map associating service names with their respective build identifiers.
	// This map is useful for dynamically referencing build information for different services.
	BUILD_SERVICES = map[string]string{
		"supervisor": BUILD_SERVICE_SUPERVISOR,
		"gateway":    BUILD_SERVICE_GATEWAY,
		"cache":      BUILD_SERVICE_CACHE,
	}
)
