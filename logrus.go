package lg

//TODO: break out into seperate repo

import "github.com/Sirupsen/logrus"

// Create a new instance of the logger. You can have any number of instances.
//var log = logrus.New()

type LogrusReceiver struct {
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

func (l *LogrusReceiver) Message(m *LogMessage) {
	e := logrus.WithFields(l.getFields(m))
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
