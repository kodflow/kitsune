package levels

type TYPE uint8

const (
	OFF TYPE = iota
	PANIC
	FATAL
	ERROR
	SUCCESS
	MESSAGE
	WARN
	INFO
	DEBUG
	TRACE

	DEFAULT = INFO
)

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

func (t TYPE) String() string {
	if t >= TYPE(len(LABELS)) {
		return "UNKNOWN"
	}
	return LABELS[t]
}

func (t TYPE) Color() string {
	if t >= TYPE(len(COLORS)) {
		return "UNKNOWN"
	}
	return COLORS[t]
}
