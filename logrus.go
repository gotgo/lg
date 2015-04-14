package lg

//TODO: break out into seperate repo

import (
	"sync"

	"github.com/Sirupsen/logrus"
)

var log = logrus.New()

type LogrusReceiver struct {
	log  *logrus.Logger
	once sync.Once
}

func (l *LogrusReceiver) getFields(m *LogMessage) logrus.Fields {
	d := logrus.Fields{}

	if m.Details != nil {
		for k, v := range m.Details {
			d[k] = v
		}
	}

	d["kind"] = m.Kind
	if m.Correlate != nil {
		d["correlate"] = m.Correlate //copy instead??
	}
	return d
}

// Current - return the current Logrus instance to customize
func (l *LogrusReceiver) Current() *logrus.Logger {
	l.once.Do(func() {
		log := logrus.New()
		log.Formatter = new(logrus.JSONFormatter)
		log.Level = logrus.DebugLevel
		l.log = log
	})
	return l.log
}

func (l *LogrusReceiver) Message(m *LogMessage) {
	e := l.Current().WithFields(l.getFields(m))

	switch m.Level {
	case LevelVerbose:
		e.Debug(m.Message)
	case LevelInform:
		e.Info(m.Message)
	case LevelWarn:
		e.Warn(m.Message)
	case LevelError:
		e.Error(m.Message)
	case LevelPanic:
		e.Panic(m.Message)
	default:
		e.Info(m.Message)
	}
}

func (l *LogrusReceiver) Levels() []Level {
	return []Level{LevelAll}
}
