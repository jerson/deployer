package tui

import (
	"bytes"
	"github.com/sirupsen/logrus"
)

// Option ...
type Option func(instance *GUI)

// WithLogOutput ...
func WithLogOutput(w *bytes.Buffer) Option {
	return func(instance *GUI) {
		instance.output = w
	}
}

// WithLog ...
func WithLog(log *logrus.Logger) Option {
	return func(instance *GUI) {
		instance.log = log
	}
}
