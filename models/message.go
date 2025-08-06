package models

import (
	"fmt"
	"strings"
)

// LogLevel represents typical log output levels
type LogLevel int

const (
	InfoLevel  LogLevel = iota
	WarnLevel  LogLevel = iota
	ErrorLevel LogLevel = iota
	DebugLevel LogLevel = iota
	FatalLevel LogLevel = iota
	OtherLevel LogLevel = iota
)

func (l LogLevel) String() string {
	switch l {
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case DebugLevel:
		return "DEBUG"
	case FatalLevel:
		return "FATAL"
	case OtherLevel:
		return "OTHER"
	}

	return ""
}

// LogMessage is meant to represent a single logical log message
// Typically denoted by an all-caps indicator near the start, such as INFO or WARN
type LogMessage struct {
	index   int
	level   LogLevel
	message string
}

func (m LogMessage) String() string {
	var styleMsg string
	switch m.level {
	case InfoLevel:
		styleMsg = infoStyle.Render(m.message)
	case WarnLevel:
		styleMsg = warnStyle.Render(m.message)
	case ErrorLevel:
		styleMsg = errorStyle.Render(m.message)
	case DebugLevel:
		styleMsg = debugStyle.Render(m.message)
	case FatalLevel:
		styleMsg = errorStyle.Render(m.message)
	case OtherLevel:
		styleMsg = m.message
	}

	return fmt.Sprintf("%4d. %s", m.index+1, styleMsg)
}

func NewLogMessage(i int, message string) LogMessage {
	m := LogMessage{
		index:   i,
		message: message,
	}

	switch {
	case strings.Contains(message, "INFO"):
		m.level = InfoLevel

	case strings.Contains(message, "WARN"):
		m.level = WarnLevel

	case strings.Contains(message, "ERRO"):
		m.level = ErrorLevel

	case strings.Contains(message, "DEBU"):
		m.level = DebugLevel

	case strings.Contains(message, "FATA"):
		m.level = FatalLevel

	default:
		m.level = OtherLevel
	}

	return m
}

type Filters struct {
	ShowInfo  bool
	ShowWarn  bool
	ShowError bool
	ShowDebug bool
	ShowFatal bool
	ShowOther bool

	FilterText string
}

func (f Filters) IncludeMessage(msg LogMessage) bool {
	if f.FilterText != "" && !strings.Contains(msg.message, f.FilterText) {
		return false
	}

	switch msg.level {
	case InfoLevel:
		return f.ShowInfo
	case WarnLevel:
		return f.ShowWarn
	case ErrorLevel:
		return f.ShowError
	case DebugLevel:
		return f.ShowDebug
	case FatalLevel:
		return f.ShowFatal
	case OtherLevel:
		return f.ShowOther
	}

	return false
}
