package lg

import "fmt"

type ConsoleLog struct {
}

func (l *ConsoleLog) Message(m *LogMessage) {
	fmt.Println(m)
}

func (l *ConsoleLog) Levels() []Level {
	return []Level{LevelAll}
}
