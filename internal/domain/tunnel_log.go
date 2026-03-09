package domain

import "fmt"

type LogLevel string

const (
	LogLevelInfo    LogLevel = "info"
	LogLevelWarning LogLevel = "warning"
	LogLevelError   LogLevel = "error"
)

type TunnelLog struct {
	ID        int
	TunnelID  int
	Level     LogLevel
	Message   string
	CreatedAt string
}

func NewTunnelLog(tunnelID int, level LogLevel, message string) (TunnelLog, error) {
	if tunnelID <= 0 {
		return TunnelLog{}, fmt.Errorf("tunnel id must be greater than 0")
	}

	if message == "" {
		return TunnelLog{}, fmt.Errorf("message must not be empty")
	}

	if level != LogLevelInfo && level != LogLevelWarning && level != LogLevelError {
		return TunnelLog{}, fmt.Errorf("invalid log level: %s", level)
	}

	return TunnelLog{
		TunnelID: tunnelID,
		Level:    level,
		Message:  message,
	}, nil
}
