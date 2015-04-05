package lg

import (
	"fmt"
	"sync"
)

const stackCallDepth = 2

type MultiLog struct {
	receiveAll      []LogReceiver
	receiversByKind map[Kind][]LogReceiver
	mu              sync.Mutex
	extraStackDepth int
	skipCallContext bool
}

func (l *MultiLog) AddReceiver(r LogReceiver) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if r.Kind() == KindAll {
		l.receivers.Add(r)
	}

	for _, k := range r.Kinds() {
		l.receiversByKind[k] = r
	}
}

func (l *MultiLog) Message(m *LogMessage) {
	l.log(m)
}

func (l *MultiLog) log(m *LogMessage) {
	const stackCallDepth = 2

	if l.skipCallContext == true {
		caller, err := CallContext(stackCallDepth + l.extraStackDepth)
		if m.Correlate == nil {
			m.Correlate = make(map[string]string)
		}

		if err != nil {
			m.Correlate["callContextError"] = err.Error()
		} else if caller.IsValid() {
			m.Correlate["lineNumber"] = caller.Line()
			m.Correlate["funcName"] = caller.Func()
			m.Correlate["file"] = caller.FileName()
		} else {
			m.Correlate["codeLocation"] = "inValid"
		}
	}

	all := l.receiveAll
	for _, log := range all {
		log.log(m)
	}

	receivers := l.receiversByKind[m.Kind]
	for _, log := range receivers {
		log.log(m)
	}
}

func (l *MultiLog) Panic(m string, err error, kv ...*KV) interface{} {
	lm := &LogMessage{
		Message: m,
		Error:   err.Error(),
		Details: kv,
		Kind:    KindPanic,
	}
	l.log(lm)
}

func (l *MultiLog) Error(m string, err error, kv ...*KV) {
	lm := &LogMessage{
		Message: m,
		Error:   err.Error(),
		Details: kv,
		Kind:    KindError,
	}
	l.log(lm)
}

func (l *MultiLog) Warn(m string, kv ...*KV) {
	lm := &LogMessage{
		Message: m,
		Details: kv,
		Kind:    KindWarn,
	}
	l.log(lm)
}

func (l *MultiLog) Inform(m string, kv ...*KV) {
	lm := &LogMessage{
		Message: m,
		Details: kv,
		Kind:    KindInform,
	}
	l.log(lm)
}

func (l *MultiLog) Verbose(m string, kv ...*KV) {
	lm := &LogMessage{
		Message: m,
		Details: kv,
		Kind:    KindVerbose,
	}
	l.log(lm)
}

/////////////////

func (l *MultiLog) MarshalFail(m string, obj interface{}, err error) {
	lm := &LogMessage{
		Message: m,
		Error:   err.Error(),
		Details: KV{"object", fmt.Sprintf("%+v", obj)},
		Kind:    KindMarshal,
	}
	l.log(lm)
}

func (l *MultiLog) UnmarshalFail(m string, data []byte, err error) {
	const arbitraryCutoffSize = 5000
	persistedData := data[:arbitraryCutoffSize]
	lm := &LogMessage{
		Message: m,
		Error:   err.Error(),
		Details: KV{"rawData", persistedData},
		Kind:    KindUnmarshal,
	}
	l.log(lm)
}

func (l *MultiLog) Timeout(m string, err error, kv ...*KV) {
	lm := &LogMessage{
		Message: m,
		Error:   err.Error(),
		Details: kv,
		Kind:    KindTimeout,
	}
	l.log(lm)
}

func (l *MultiLog) ConnectFail(m string, err error, kv ...*KV) {
	lm := &LogMessage{
		Message: m,
		Error:   err.Error(),
		Details: kv,
		Kind:    KindConnect,
	}
	l.log(lm)
}
