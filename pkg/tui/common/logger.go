package common

import (
	"io"

	"github.com/jerson/deployer/pkg/deployer"
	"github.com/jerson/deployer/pkg/formatter"
	"github.com/sirupsen/logrus"
)

// CustomLogger ...
func CustomLogger(output io.Writer, format *formatter.Console) *logrus.Logger {

	instance := logrus.New()
	if format == nil {
		format = &formatter.Console{
			FieldsOrder:          []string{deployer.LogFieldRunner, deployer.LogFieldDeployment},
			IgnoreNotFoundFields: true,
		}
	}
	instance.SetFormatter(format)
	instance.SetOutput(output)
	instance.SetLevel(logrus.TraceLevel)

	return instance
}
