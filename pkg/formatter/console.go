package formatter

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"sort"
	"strings"
)

// Console - logrus formatter, implements logrus.Main
type Console struct {
	FieldsOrder          []string // default: fields sorted alphabetically
	IgnoreNotFoundFields bool     // ignore not fields on FieldsOrder
	TimestampFormat      string   // default: time.StampMilli = "Jan _2 15:04:05.000"
	HideKeys             bool     // show [fieldValue] instead of [fieldKey:fieldValue]
	NoColors             bool     // disable colors
	NoFieldsColors       bool     // color only level, default is level + fields
	ShowFullLevel        bool     // true to show full level [WARNING] instead [WARN]
	TrimMessages         bool     // true to trim whitespace on messages
}

// Format an log entry
func (f *Console) Format(entry *logrus.Entry) ([]byte, error) {
	levelColor := getColorByLevel(entry.Level)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		//timestampFormat = time.StampMilli
	}

	// output buffer
	b := &bytes.Buffer{}

	// write time
	b.WriteString(entry.Time.Format(timestampFormat))

	// write level
	level := strings.ToUpper(entry.Level.String())

	if !f.NoColors {
		_, _ = fmt.Fprintf(b, "\x1b[%dm", levelColor)
	}

	b.WriteString(" [")
	if f.ShowFullLevel {
		b.WriteString(level)
	} else {
		b.WriteString(level[:1])
	}
	b.WriteString("] ")

	if !f.NoColors && f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	// write fields
	if f.FieldsOrder == nil {
		f.writeFields(b, entry)
	} else {
		f.writeOrderedFields(b, entry)
	}

	if !f.NoColors && !f.NoFieldsColors {
		if entry.Level == logrus.TraceLevel || entry.Level == logrus.DebugLevel {
			_, _ = fmt.Fprintf(b, "\x1b[%dm", levelColor)
		} else {
			b.WriteString("\x1b[0m")
		}
	}

	// write message
	if f.TrimMessages {
		b.WriteString(strings.TrimSpace(entry.Message))
	} else {
		b.WriteString(entry.Message)
	}

	if entry.HasCaller() {
		_, _ = fmt.Fprintf(
			b,
			" (%s:%d %s)",
			entry.Caller.File,
			entry.Caller.Line,
			entry.Caller.Function,
		)
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *Console) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Console) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry, field)
		}
	}

	if length > 0 && !f.IgnoreNotFoundFields {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if foundFieldsMap[field] == false {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Console) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	if f.HideKeys {
		_, _ = fmt.Fprintf(b, "[%v] ", entry.Data[field])
	} else {
		_, _ = fmt.Fprintf(b, "[%s:%v] ", field, entry.Data[field])
	}
}

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 90
	colorWhite  = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.TraceLevel:
		return colorGray
	case logrus.DebugLevel:
		return colorWhite
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}
