package deployer

import (
	"context"
	"github.com/sirupsen/logrus"
	"reflect"
	"runtime"
	"strings"
)

// FunctionName ...
func FunctionName(dependency interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(dependency).Pointer()).Name()
	parts := strings.Split(name, ".")
	return strings.Join(parts[len(parts)-1:], ".")
}

// PrintTitle ...
func PrintTitle(ctx context.Context, title string) {
	log := Log(ctx)
	log.Info(title)
}

// Log ...
func Log(ctx context.Context) *logrus.Entry {

	var log *logrus.Entry
	if ctx.Value(ContextKeyLog) != nil {
		log = ctx.Value(ContextKeyLog).(*logrus.Entry)
	} else {
		instance := logrus.New()
		instance.SetFormatter(&logrus.TextFormatter{
			DisableTimestamp: true,
			FullTimestamp:    false,
			ForceColors:      true,
		})
		instance.SetLevel(logrus.TraceLevel)
		if ctx.Value(ContextKeyLogLevel) != nil {
			level := ctx.Value(ContextKeyLogLevel).(logrus.Level)
			instance.SetLevel(level)
		}
		log = instance.WithFields(nil)
	}
	return log
}
