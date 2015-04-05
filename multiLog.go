package lg

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/gotgo/lg/er"
)

const stackCallDepth = 2

type MultiLog struct {
	receiveAll       []LogReceiver
	receiversByLevel map[Level][]LogReceiver
	mu               sync.Mutex
	extraStackDepth  int
	skipCallContext  bool
	count            int64
}

func (l *MultiLog) AddReceiver(r LogReceiver) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, level := range r.Levels() {
		if level == LevelAll {
			l.receiveAll = append(l.receiveAll, r)
			break
		} else {
			l.receiversByLevel[level] = append(l.receiversByLevel[level], r)
		}
	}
}

func (l *MultiLog) Message(m *LogMessage) {
	l.log(m)
}

func getError(err error) KV {
	if err == nil {
		return nil
	}

	var generic interface{}
	generic = err
	derr, ok := generic.(*er.Error)
	if ok {
		return KV{"error": derr}
	} else {
		return KV{"error": err.Error()}
	}
}

func (l *MultiLog) log(m *LogMessage) {
	//do we want this??
	m.index = atomic.AddInt64(&l.count, 1)

	const stackCallDepth = 2

	if l.skipCallContext == true {
		caller, err := er.CallerContext(stackCallDepth + l.extraStackDepth)
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
	kerr := getError(err)

	lm := &LogMessage{
		Message: m,
		Details: collapse(kv, kerr),
		Level:   LevelPanic,
		Kind:    KindPanic,
	}
	l.log(lm)
	return lm
}

func (l *MultiLog) Error(m string, err error, kv ...KV) {
	ekv := getError(err)
	lm := &LogMessage{
		Message: m,
		Details: collapse(kv, ekv),
		Level:   LevelError,
		Kind:    KindError,
	}
	l.log(lm)
}

func (l *MultiLog) Warn(m string, kv ...KV) {
	lm := &LogMessage{
		Message: m,
		Details: collapse(kv),
		Level:   LevelWarn,
		Kind:    KindWarn,
	}
	l.log(lm)
}

func (l *MultiLog) Inform(m string, kv ...KV) {
	lm := &LogMessage{
		Message: m,
		Details: collapse(kv),
		Level:   LevelInform,
		Kind:    KindInform,
	}
	l.log(lm)
}

func (l *MultiLog) Verbose(m string, kv ...KV) {
	lm := &LogMessage{
		Message: m,
		Details: collapse(kv),
		Level:   LevelVerbose,
		Kind:    KindVerbose,
	}
	l.log(lm)
}

////////////////////////
// SPECIALIZED ERRORS
////////////////////////

func (l *MultiLog) MarshalFail(m string, obj interface{}, err error) {
	kerr := getError(err)
	lm := &LogMessage{
		Message: m,
		Details: collapse([]KV{}, KV{"object": fmt.Sprintf("%+v", obj)}, kerr),
		Level:   LevelError,
		Kind:    KindMarshal,
	}
	l.log(lm)
}

func (l *MultiLog) UnmarshalFail(m string, data []byte, err error) {
	const arbitraryCutoffSize = 5000
	persistedData := data[:arbitraryCutoffSize]
	kerr := getError(err)
	lm := &LogMessage{
		Message: m,
		Details: collapse([]KV{}, KV{"rawData": persistedData}, kerr),
		Level:   LevelError,
		Kind:    KindUnmarshal,
	}
	l.log(lm)
}

func (l *MultiLog) Timeout(m string, err error, kv ...KV) {
	kerr := getError(err)
	lm := &LogMessage{
		Message: m,
		Details: collapse(kv, kerr),
		Level:   LevelWarn,
		Kind:    KindTimeout,
	}
	l.log(lm)
}

func (l *MultiLog) ConnectFail(m string, err error, kv ...KV) {
	kerr := getError(err)
	lm := &LogMessage{
		Message: m,
		Details: collapse(kv, kerr),
		Kind:    KindConnect,
		Level:   LevelWarn,
	}
	l.log(lm)
}
