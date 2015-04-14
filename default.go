package lg

import "sync"

// the default logger
var defaultLogger LevelLogger

// the current logger
var current LevelLogger
var mu sync.Mutex

func init() {
	defaultLogger = &MultiLog{}
	current = defaultLogger
}

func Use(l LevelLogger) {
	mu.Lock()
	defer mu.Unlock()
	current = l
}

func AddReceiver(r LogReceiver) {
	mu.Lock()
	defer mu.Unlock()
	current.AddReceiver(r)
}

func Panic(m string, err error, kv ...KV) interface{} {
	return current.Panic(m, err, kv...)
}

func Error(m string, err error, kv ...KV) {
	current.Error(m, err, kv...)
}

func Warn(m string, kv ...KV) {
	current.Warn(m, kv...)
}

func Inform(m string, kv ...KV) {
	current.Inform(m, kv...)
}

func Verbose(m string, kv ...KV) {
	current.Verbose(m, kv...)
}

func Message(m *LogMessage) {
	current.Message(m)
}
