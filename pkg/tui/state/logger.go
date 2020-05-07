package state

import (
	"bytes"

	"github.com/jerson/deployer/pkg/formatter"
	"github.com/jerson/deployer/pkg/tui/common"
	"github.com/sirupsen/logrus"
)

// LoggerState ...
type LoggerState struct {
	Output *bytes.Buffer
	Log    *logrus.Logger
}

// NewLoggerState ...
func NewLoggerState(format *formatter.Console) *LoggerState {
	instance := &LoggerState{
		Output: &bytes.Buffer{},
	}
	instance.Log = common.CustomLogger(instance.Output, format)
	return instance
}

// Reset ...
func (l *LoggerState) Reset() {
	l.Output = &bytes.Buffer{}
}
