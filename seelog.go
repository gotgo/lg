package lg

//TODO: break out into seperate repo

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/cihub/seelog"
)

//even though Logger is an instance, all instances are using the same
//singleton logger. it's unclear if this should stay this way
var logger seelog.LoggerInterface

func init() {
	DisableLog()
}

// DisableLog disables all library log output
func DisableLog() {
	logger = seelog.Disabled
}

// UseLogger uses a specified seelog.LoggerInterface to output library log.
// Use this func if you are using Seelog logging system in your app.
func UseLogger(newLogger seelog.LoggerInterface) {
	logger = newLogger
	newLogger.SetAdditionalStackDepth(2)
}

// SetLogWriter uses a specified io.Writer to output library log.
// Use this func if you are not using Seelog logging system in your app.
func SetLogWriter(writer io.Writer) error {
	if writer == nil {
		return errors.New("Nil writer")
	}

	newLogger, err := seelog.LoggerFromWriterWithMinLevel(writer, seelog.TraceLvl)
	if err != nil {
		return err
	}

	UseLogger(newLogger)
	return nil
}

// Call this before app shutdown
func FlushLog() {
	logger.Flush()
}

type RequstLogger interface {
	Capture(service string, name string, args map[string]string, duration int, outcome bool)
}

type SeeLog struct {
}

func (l *SeeLog) Message(m *LogMessage) {
	bytes, err := json.Marshal(m)
	if err != nil {
		panic("could not log event because of marshal fail")
	}

	str := string(bytes)
	switch m.Level {
	case LevelInform:
		logger.Info(str)
	case LevelVerbose:
		logger.Debug(str)
	case LevelWarn:
		logger.Warn(str)
	case LevelError:
		logger.Error(str)
	case LevelPanic:
		logger.Critical(str)
	default:
		logger.Info(str)
	}
}

func (l *SeeLog) Levels() []Level {
	return []Level{LevelAll}
}
