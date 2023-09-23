package config

const (
	PATH_SERVICES = "/etc/kitsune/"
	PATH_BIN      = "/usr/local/bin/"
	PATH_RUN      = "/var/run/kitsune/"
	PATH_LOGS     = "/var/log/kitsune/"
)

var PATHS = []string{
	PATH_SERVICES,
	PATH_RUN,
	PATH_BIN,
	PATH_LOGS,
}
