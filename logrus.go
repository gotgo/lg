package lg

import "github.com/Sirupsen/logrus"

// Create a new instance of the logger. You can have any number of instances.
//var log = logrus.New()

type Logrus struct {
}

func (l *Logrus) getFields(m *LogMessage) logrus.Fields {
	d := Fields{}

	if m.Details != nil {
		for k, v := range m.Details {
			d[k] = v
		}
	}

	d["kind"] = m.Kind
	d["error"] = m.Error
	if m.Correlate != nil {
		d["correlate"] = m.Correlate //copy instead??
	}
	return d
}

func (l *Logrus) Message(m *LogMessage) {
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
		d.Info(m.Message)
	}
}

func (l *Logrus) Levels() []Level {
	return []Level{LevelAll}
}
