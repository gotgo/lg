package lg

// the default logger
var defaultLogger BasicLogger

// the current logger
var current BasicLogger

func init() {
	defaultLogger = NewNoOpLogger()
	current = defaultLogger
}

func Panic(m string, err error, kv ...*KV) interface{} {
	return current.Panic(m, err, kv)
}

func Error(m string, err error, kv ...*KV) {
	current.Error(m, err, kv)
}

func Warn(m string, kv ...*KV) {
	current.Warn(m, kv)
}

func Inform(m string, kv ...*KV) {
	current.Inform(m, kv)
}

func Verbose(m string, kv ...*KV) {
	current.Verbose(m, kv)
}

func Message(m *LogMessage) {
	current.Message(m)
}
