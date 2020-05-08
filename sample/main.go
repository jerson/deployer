//usr/bin/env go run "$0" "$@"; exit "$?"

package main

import (
	"bytes"
	"context"
	"github.com/jerson/deployer/pkg/deployer"
	"github.com/jerson/deployer/pkg/entities"
	"github.com/jerson/deployer/pkg/formatter"
	"github.com/jerson/deployer/pkg/tui"
	"github.com/jerson/deployer/sample/deploy"
	"github.com/sirupsen/logrus"
)

var deployments = entities.Deployments{
	deploy.Sample1,
	deploy.Sample3,
	deploy.Sample5,
}

func run() error {
	buf := &bytes.Buffer{}

	log := logger()
	log.SetOutput(buf)

	ctx := context.Background()
	ctx = context.WithValue(ctx, deployer.ContextKeyLog, log.WithFields(nil))

	app := deployer.NewDeployer()
	app.AddDeployments(ctx, deployments...)

	tuiConfig := tui.Config{Name: "Sample App"}
	gui := tui.NewGUI(app, tuiConfig, tui.WithLog(log), tui.WithLogOutput(buf))
	return gui.Render()
}

func logger() *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.TraceLevel)
	log.SetFormatter(&formatter.Console{
		FieldsOrder: []string{
			deployer.LogFieldRunner,
			deployer.LogFieldDeployment,
			deployer.LogFieldMethod,
		},
		IgnoreNotFoundFields: false,
	})
	return log
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
