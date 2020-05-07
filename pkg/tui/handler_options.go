package tui

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/enescakir/emoji"
	"github.com/jerson/deployer/pkg/deployer"
	"github.com/jerson/deployer/pkg/entities"
	"github.com/jerson/deployer/pkg/formatter"
	"github.com/jerson/deployer/pkg/tui/common"
	"github.com/jerson/deployer/pkg/tui/layout"
	"github.com/jerson/deployer/pkg/tui/state"
	"github.com/jerson/deployer/pkg/tui/widget"
	"github.com/jesseduffield/gocui"
	"github.com/sirupsen/logrus"
	padUtf8 "github.com/willf/pad/utf8"
)

var deploymentFormatter = &formatter.Console{
	FieldsOrder:          []string{deployer.LogFieldRunner},
	IgnoreNotFoundFields: true,
}

func (gui *GUI) getDeploymentOptions(deployment entities.Deployment) []string {
	deploymentOptions := []string{
		"Uninstall",
		"Re-Deploy",
		"Upgrade",
		"Health Check",
	}
	installerStatus := gui.deploymentStatusType(deployment)

	if installerStatus == entities.InstallerStatusOff || installerStatus == entities.InstallerStatusError {
		deploymentOptions[0] = "Install"
	}
	return deploymentOptions
}
func (gui *GUI) onChangeOption() error {

	view, err := gui.getOptionsView()
	if err != nil {
		return err
	}

	deployment := gui.getSelectedDeployment()
	if deployment == nil {
		return nil
	}
	optionsState := gui.state.Panels.Options

	if optionsState.SelectedIndex < 0 {
		return nil
	}

	modal := widget.NewModal("options_modal", gui.g, gui.focus, gui.log.WithFields(nil), view)
	installerStatus := gui.deploymentStatusType(deployment)
	deploymentOptions := gui.getDeploymentOptions(deployment)
	err = modal.ShowDefault(fmt.Sprintf("%s deployment", deploymentOptions[optionsState.SelectedIndex]), "Are you sure?",
		func(g *gocui.Gui, v *gocui.View) error {

			switch optionsState.SelectedIndex {
			case 0:
				if installerStatus == entities.InstallerStatusOff || installerStatus == entities.InstallerStatusError {
					return gui.changeToConsoleTab(deployment, gui.startDeployment)
				}
				return gui.changeToConsoleTab(deployment, gui.stopDeployment)
			case 1:
				return gui.changeToConsoleTab(deployment, gui.restartDeployment)
			case 2:
				return gui.changeToConsoleTab(deployment, gui.upgradeDeployment)
			case 3:
				return gui.changeToConsoleTab(deployment, gui.healthDeployment)
			}
			return errors.New("not implemented")
		},
		nil,
	)

	return nil
}

func (gui *GUI) changeToConsoleTab(deployment entities.Deployment, callback func(deployment entities.Deployment) error) (err error) {

	err = gui.onMainTabChange(1)
	if err != nil {
		return err
	}
	return callback(deployment)
}
func (gui *GUI) stopDeployment(deployment entities.Deployment) error {
	gui.log.WithField(common.LogFieldMethod, "stopDeployment").Info(deployer.FunctionName(deployment))

	view, err := gui.getMainView()
	if err != nil {
		return err
	}

	log := gui.deploymentLog(deployment, view)
	ctx := context.WithValue(context.Background(), deployer.ContextKeyLog, log.WithFields(nil))

	go func() {
		_ = gui.cli.Uninstall(ctx, deployment)
	}()

	return nil
}

func (gui *GUI) startDeployment(deployment entities.Deployment) error {
	gui.log.WithField(common.LogFieldMethod, "startDeployment").Info(deployer.FunctionName(deployment))

	view, err := gui.getMainView()
	if err != nil {
		return err
	}

	log := gui.deploymentLog(deployment, view)
	ctx := context.WithValue(context.Background(), deployer.ContextKeyLog, log.WithFields(nil))

	go func() {
		_ = gui.cli.Install(ctx, deployment)

	}()

	return nil
}

func (gui *GUI) upgradeDeployment(deployment entities.Deployment) error {
	gui.log.WithField(common.LogFieldMethod, "upgradeDeployment").Info(deployer.FunctionName(deployment))

	view, err := gui.getMainView()
	if err != nil {
		return err
	}

	log := gui.deploymentLog(deployment, view)
	ctx := context.WithValue(context.Background(), deployer.ContextKeyLog, log.WithFields(nil))
	go func() {
		_ = gui.cli.Upgrade(ctx, deployment)
	}()

	return nil

}

func (gui *GUI) restartDeployment(deployment entities.Deployment) error {
	gui.log.WithField(common.LogFieldMethod, "restartDeployment").Info(deployer.FunctionName(deployment))

	view, err := gui.getMainView()
	if err != nil {
		return err
	}

	log := gui.deploymentLog(deployment, view)
	ctx := context.WithValue(context.Background(), deployer.ContextKeyLog, log.WithFields(nil))

	go func() {
		err := gui.cli.Uninstall(ctx, deployment)
		if err != nil {
			return
		}
		_ = gui.cli.Install(ctx, deployment)
	}()

	return nil

}

func (gui *GUI) healthDeployment(deployment entities.Deployment) error {
	gui.log.WithField(common.LogFieldMethod, "healthDeployment").Info(deployer.FunctionName(deployment))

	view, err := gui.getMainView()
	if err != nil {
		return err
	}

	log := gui.deploymentLog(deployment, view)
	ctx := context.WithValue(context.Background(), deployer.ContextKeyLog, log.WithFields(nil))

	go func() {
		_ = gui.cli.Health(ctx, deployment)
	}()

	return nil

}

func (gui *GUI) deploymentLog(deployment entities.Deployment, w io.Writer) *logrus.Logger {
	key := deployer.FunctionName(deployment)
	gui.state.Panels.Main.Console[key] = state.NewLoggerState(deploymentFormatter)
	deploymentState := gui.state.Panels.Main.Console[key]

	log := deploymentState.Log
	writer := common.NewConditionalWriter(func() (success bool, err error) {
		currentDeployment := gui.getSelectedDeployment()
		if currentDeployment == nil {
			return false, nil
		}
		success = gui.state.Panels.Main.TabIndex == 1 && deployer.FunctionName(currentDeployment) == key
		return success, nil
	}, w)
	log.SetOutput(io.MultiWriter(writer, deploymentState.Output))
	return log
}

func (gui *GUI) updateOptions() error {
	gui.g.Update(func(g *gocui.Gui) error {
		err := gui.renderOptions()
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}
func (gui *GUI) renderOptions() error {
	gui.g.Update(func(g *gocui.Gui) error {

		view, err := gui.getOptionsView()
		if err != nil {
			return err
		}

		deployment := gui.getSelectedDeployment()
		if deployment == nil {
			return nil
		}
		view.Clear()

		deploymentOptions := gui.getDeploymentOptions(deployment)
		for _, option := range deploymentOptions {
			menu := gui.bufferMenuItem(emoji.BackhandIndexPointingRight, option)
			_, err := fmt.Fprintln(view, menu)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return nil
}

func (gui *GUI) deploymentStatusType(deployment entities.Deployment) entities.InstallerStatusType {
	installerStatus := entities.InstallerStatusOff
	status := gui.cli.DeploymentStatus(deployment)
	if status != nil && status.Health != nil {
		installerStatus = status.Health.Status()
	}
	return installerStatus
}

func (gui *GUI) bufferMenuItem(icon interface{}, name string) *bytes.Buffer {
	var output bytes.Buffer
	menu := fmt.Sprintf("[%v %s]", icon, padUtf8.Right(name, 14, " "))
	output.WriteString(layout.Styles.Button.Sprint(menu))
	return &output
}
