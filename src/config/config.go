package config

import "time"

// DEFAULT_RETRY_INTERVAL, DEFAULT_RETRY_MAX, DEFAULT_TIMEOUT, and DEFAULT_CACHE
// are constants representing default timing configurations in the application.
var (
	// DEFAULT_RETRY_INTERVAL defines the default duration to wait between retries.
	// This interval is used in operations where an action is attempted repeatedly after failure.
	DEFAULT_RETRY_INTERVAL time.Duration = 1

	// DEFAULT_RETRY_MAX defines the default maximum duration for retrying operations.
	// This duration limits the total time spent in retry loops.
	DEFAULT_RETRY_MAX time.Duration = 3

	// DEFAULT_TIMEOUT defines the default timeout duration for operations.
	// This duration is used as a standard limit for operations to complete, beyond which
	// they may be aborted or considered failed.
	DEFAULT_TIMEOUT time.Duration = 15

	// DEFAULT_CACHE defines the default duration for caching elements.
	// This duration specifies how long certain data or objects should be kept in cache
	// before being refreshed or invalidated.
	DEFAULT_CACHE time.Duration = 15
)
