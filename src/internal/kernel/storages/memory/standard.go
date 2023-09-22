package memory

var instance *Memory = nil

func standard() *Memory {
	if instance == nil {
		instance = New()
	}
	return instance
}

func Clear() {
	standard().Clear()
}

func Store(key string, value interface{}) {
	standard().Store(key, value)
}

func Read(key string) (interface{}, bool) {
	return standard().Read(key)
}

func Delete(key string) {
	standard().Delete(key)
}

func Exists(key string) bool {
	return standard().Exists(key)
}
