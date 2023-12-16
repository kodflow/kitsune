package config

import "time"

var (
	DEFAULT_RETRY_INTERVAL time.Duration = 1
	DEFAULT_RETRY_MAX      time.Duration = 3
	DEFAULT_TIMEOUT        time.Duration = 15
	DEFAULT_CACHE          time.Duration = 15
)
