package multithread

// IsMaster returns true if the current process is the master process.
func IsMaster() bool {
	return !IsWorker()
}
