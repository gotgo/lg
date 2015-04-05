package lg

var noOpLogger *NoOpLogger

func init() {
	noOpLogger = new(NoOpLogger)
}

func NewNoOpLogger() *NoOpLogger {
	return noOpLogger
}

type NoOpLogger struct {
}

func (l *NoOpLogger) Panic(m string, err error, kv ...*KV) interface{} {
	return &LogMessage{Message: m, Error: err.Error(), Details: kv, Kind: Panic}
}

func (l *NoOpLogger) Error(m string, err error, kv ...*KV) {}
func (l *NoOpLogger) Warn(m string, kv ...*KV)             {}
func (l *NoOpLogger) Inform(m string, kv ...*KV)           {}
func (l *NoOpLogger) Verbose(m string, kv ...*KV)          {}
func (l *NoOpLogger) Message(m *LogMessage)                {}
