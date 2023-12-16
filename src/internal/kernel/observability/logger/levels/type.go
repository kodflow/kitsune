package levels

// TYPE represents the log level type.
type TYPE uint8

// Log level constants.
const (
	OFF     TYPE = iota // OFF represents the log level for disabling logging.
	PANIC               // PANIC represents the log level for critical errors that lead to program termination.
	FATAL               // FATAL represents the log level for fatal errors that may lead to program termination.
	ERROR               // ERROR represents the log level for non-fatal errors.
	WARN                // WARN represents the log level for warning messages.
	SUCCESS             // SUCCESS represents the log level for successful operations.
	MESSAGE             // MESSAGE represents the log level for general messages.
	INFO                // INFO represents the log level for informational messages.
	DEBUG               // DEBUG represents the log level for debugging messages.
	TRACE               // TRACE represents the log level for detailed tracing messages.

	DEFAULT = INFO // DEFAULT represents the default log level.
)

// LABELS maps log level constants to their corresponding labels.
var LABELS = []string{
	OFF:     "OFF",
	PANIC:   "PANIC",
	FATAL:   "FATAL",
	ERROR:   "ERROR",
	SUCCESS: "SUCCESS",
	MESSAGE: "MESSAGE",
	WARN:    "WARN",
	INFO:    "INFO",
	DEBUG:   "DEBUG",
	TRACE:   "TRACE",
}

// COLORS maps log level constants to their corresponding colors.
var COLORS = []string{
	PANIC:   "9",
	FATAL:   "160",
	ERROR:   "1",
	SUCCESS: "2",
	MESSAGE: "7",
	WARN:    "3",
	INFO:    "4",
	DEBUG:   "6",
	TRACE:   "7",
}

// Int returns the integer representation of the log level.
func (t TYPE) Int() uint8 {
	return uint8(t)
}

// String returns the string representation of the log level.
func (t TYPE) String() string {
	if t >= TYPE(len(LABELS)) {
		return "UNKNOWN"
	}
	return LABELS[t]
}

// Color returns the color code associated with the log level.
func (t TYPE) Color() string {
	if t >= TYPE(len(COLORS)) {
		return "UNKNOWN"
	}
	return COLORS[t]
}
