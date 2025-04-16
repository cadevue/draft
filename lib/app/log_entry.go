package app

import "fmt"

type CommandType int

const (
	Set CommandType = iota
	Del
	Append
)

type LogEntry struct {
	term    int // term when entry was received by leader
	command CommandType
	key     string
	value   string
}

func NewLogEntry(term int, command CommandType, key string, value string) *LogEntry {
	return &LogEntry{
		term:    term,
		command: command,
		key:     key,
		value:   value,
	}
}

func (le *LogEntry) GetTerm() int {
	return le.term
}

func (le *LogEntry) GetCommand() CommandType {
	return le.command
}

func (le *LogEntry) GetKey() string {
	return le.key
}

func (le *LogEntry) GetValue() string {
	return le.value
}

func (le *LogEntry) String() string {
	commandStr := ""
	switch le.command {
	case Set:
		commandStr = "set"
	case Del:
		commandStr = "del"
	case Append:
		commandStr = "append"
	default:
		commandStr = "unknown"
	}
	return fmt.Sprintf("LogEntry{Command: %s, Key: %s, Value: %s}", commandStr, le.key, le.value)
}
