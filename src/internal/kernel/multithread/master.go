package multithread

func IsMaster() bool {
	return !IsWorker()
}
