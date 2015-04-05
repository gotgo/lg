package lg

import (
	"fmt"
	"strconv"
	"sync"
)

const stackCallDepth = 2

type MultiLog struct {
	receiveAll       []LogReceiver
	receiversByLevel map[Level][]LogReceiver
	mu               sync.Mutex
	extraStackDepth  int
	skipCallContext  bool
}

func (l *MultiLog) AddReceiver(r LogReceiver) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, level := range r.Levels() {
		if level == LevelAll {
			l.receiveAll = append(l.receiveAll, r)
		} else {
			l.receiversByLevel[level] = append(l.receiversByLevel[level], r)
		}
	}
}

func (l *MultiLog) Message(m *LogMessage) {
	l.log(m)
}

func (l *MultiLog) log(m *LogMessage) {
	const stackCallDepth = 2

	if l.skipCallContext == true {
		caller, err := CallerContext(stackCallDepth + l.extraStackDepth)
		if m.Correlate == nil {
			m.Correlate = make(map[string]string)
		}

		if err != nil {
			m.Correlate["callContextError"] = err.Error()
		} else {
			m.Correlate["lineNumber"] = strconv.Itoa(caller.LineNumber)
			m.Correlate["funcName"] = caller.FuncName
			m.Correlate["file"] = caller.Filename
		}
	}

	all := l.receiveAll
	for _, log := range all {
		log.Message(m)
	}

	receivers := l.receiversByLevel[m.Level]
	for _, log := range receivers {
		log.Message(m)
	}
}

func (l *MultiLog) Panic(m string, err error, kv ...KV) interface{} {
	lm := &LogMessage{
		Message: m,
		Error:   err.Error(),
		Details: collapse(kv),
		Kind:    KindPanic,
	}
	l.log(lm)
	//TODO: format as JSON??
	return lm
}

func (l *MultiLog) Error(m string, err error, kv ...KV) {
	lm := &LogMessage{
		Message: m,
		Error:   err.Error(),
		Details: collapse(kv),
		Kind:    KindError,
	}
	l.log(lm)
}

func (l *MultiLog) Warn(m string, kv ...KV) {
	lm := &LogMessage{
		Message: m,
		Details: collapse(kv),
		Kind:    KindWarn,
	}
	l.log(lm)
}

func (l *MultiLog) Inform(m string, kv ...KV) {
	lm := &LogMessage{
		Message: m,
		Details: collapse(kv),
		Kind:    KindInform,
	}
	l.log(lm)
}

func (l *MultiLog) Verbose(m string, kv ...KV) {
	lm := &LogMessage{
		Message: m,
		Details: collapse(kv),
		Kind:    KindVerbose,
	}
	l.log(lm)
}

/////////////////

func (l *MultiLog) MarshalFail(m string, obj interface{}, err error) {
	lm := &LogMessage{
		Message: m,
		Error:   err.Error(),
		Details: KV{"object": fmt.Sprintf("%+v", obj)},
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
		Details: KV{"rawData": persistedData},
		Kind:    KindUnmarshal,
	}
	l.log(lm)
}

func (l *MultiLog) Timeout(m string, err error, kv ...KV) {
	lm := &LogMessage{
		Message: m,
		Error:   err.Error(),
		Details: collapse(kv),
		Kind:    KindTimeout,
	}
	l.log(lm)
}

func (l *MultiLog) ConnectFail(m string, err error, kv ...KV) {
	lm := &LogMessage{
		Message: m,
		Error:   err.Error(),
		Details: collapse(kv),
		Kind:    KindConnect,
	}
	l.log(lm)
}
