package lg

type KV map[string]interface{}

type Kind string

const (
	KindPanic   = "panic"
	KindError   = "error"
	KindWarn    = "warn"
	KindInform  = "inform"
	KindVerbose = "verbose"

	KindTimeout   = "timeout"
	KindConnect   = "connect"
	KindMarshal   = "marshal"
	KindUnmarshal = "unmarshal"
)

type Level uint8

const (
	LevelAll Level = iota
	LevelVerbose
	LevelInform
	LevelWarn
	LevelError
	LevelPanic
)

type LevelLogger interface {
	// usage panic(lg.Panic("crap!", err))
	Panic(m string, err error, kv ...KV) interface{}
	Error(m string, err error, kv ...KV) //do we have a message here?
	Warn(m string, kv ...KV)

	// lg.Inform("Server Started", lg.KV{"config", config, "port", port}
	Inform(m string, kv ...KV)

	Verbose(m string, kv ...KV) //debug

	Message(m *LogMessage)
}

type Logger interface {
	LevelLogger
	// MarshalFail occurs when an object fails to marshal.
	// Solving a Marshal failure requires discovering which object type and what data was
	// in that instance that could have caused the failure. This is why the interface requires
	// the object
	MarshalFail(m string, obj interface{}, err error)
	// UnmarshalFail occures when a stream is unable to be unmarshalled.
	// Solving a unmarshal failure requires knowing what object type, which field, and
	// what's wrong with the source data that causes the problem
	UnmarshalFail(m string, data []byte, err error)

	Timeout(m string, err error, kv ...KV)
	ConnectFail(m string, err error, kv ...KV)
}

type LogReceiver interface {
	Message(m *LogMessage)
	Levels() []Level
}

type LogMessage struct {
	Message string `json:"message"`
	Details KV     `json:"details,omitempty"`
	Kind    Kind   `json:"kind,omitempty"`
	Level   Level  `json:"level,omitempty"`

	Error string `json:"error,omitempty"` //TODO

	// Options
	// 1. traceuid
	// 2. spanuid
	// 3. line #
	// 4. file name
	// 5. func name
	Correlate map[string]string `json:"correlate,omitempty"`
}

func collapse(kvs []KV) KV {
	if kvs == nil {
		return nil
	} else if len(kvs) == 0 {
		return KV{}
	} else if len(kvs) == 1 {
		return kvs[0]
	} else {
		data := KV{}
		for _, mp := range kvs {
			for k, v := range mp {
				data[k] = v
			}
		}
		return data
	}
}
