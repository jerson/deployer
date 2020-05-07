package tui

import (
	"bytes"
	"fmt"

	"github.com/jerson/deployer/pkg/deployer"
	"github.com/jerson/deployer/pkg/entities"
	"github.com/jerson/deployer/pkg/tui/common"
	"github.com/jerson/deployer/pkg/tui/layout"

	"github.com/enescakir/emoji"
	"github.com/jesseduffield/gocui"
)

func (gui *GUI) onChangeDeployment() error {
	log := gui.log.WithField(common.LogFieldMethod, "onChangeDeployment")
	log.Debug("changed")

	err := gui.renderOptions()
	if err != nil {
		return err
	}
	err = gui.updateMainTabs()
	if err != nil {
		return err
	}
	err = gui.updateDeployments()
	if err != nil {
		return err
	}

	return nil
}

func (gui *GUI) updateDeployments() error {

	buffer := gui.deploymentsRow()

	gui.g.Update(func(g *gocui.Gui) error {
		view, err := gui.getDeploymentsView()
		if err != nil {
			return err
		}
		view.Clear()
		view.ContainsList = true

		_, err = buffer.WriteTo(view)
		if err != nil {
			return err
		}

		return nil
	})
	return nil
}

func (gui *GUI) deploymentsRow() *bytes.Buffer {
	var output bytes.Buffer
	deployments := gui.cli.Deployments()
	total := len(deployments)
	for i, deployment := range deployments {
		_, err := output.WriteString(gui.deploymentRow(deployment))
		if err != nil {
			return &output
		}
		if i < (total - 1) {
			output.WriteString("\n")
		}
	}
	return &output
}

func (gui *GUI) deploymentRow(deployment entities.Deployment) string {
	status := gui.cli.DeploymentStatus(deployment)
	color := emoji.BlackCircle
	if status != nil {
		color = gui.deploymentStatusIcon(status.Health)
	}

	name := deployer.FunctionName(deployment)
	selected := gui.getSelectedDeployment()

	if selected != nil && deployer.FunctionName(selected) == name {
		return fmt.Sprintf("%v %s", color, layout.Styles.ButtonActive.Sprint(name))
	}

	return fmt.Sprintf("%v %s", color, layout.Styles.Default.Sprint(name))
}

func (gui *GUI) deploymentStatusIcon(status *entities.InstallerStatus) emoji.Emoji {
	color := emoji.BlackCircle
	if status != nil {

		switch status.Status() {
		case entities.InstallerStatusRunning:
			color = emoji.GreenCircle
		case entities.InstallerStatusError:
			color = emoji.RedCircle
		case entities.InstallerStatusIdle:
			color = emoji.YellowCircle
		case entities.InstallerStatusOff:
			color = emoji.BlackCircle
		}

	}
	return color
}
func (gui *GUI) deploymentStatusName(status *entities.InstallerStatus) string {
	name := ""
	if status != nil {

		switch status.Status() {
		case entities.InstallerStatusRunning:
			name = layout.Styles.Success.Sprint("OK")
		case entities.InstallerStatusError:
			name = layout.Styles.Error.Sprint("ERROR")
		case entities.InstallerStatusIdle:
			name = layout.Styles.Warning.Sprint("WAITING")
		case entities.InstallerStatusOff:
			name = layout.Styles.Disabled.Sprint("NONE")
		}

	}
	return name
}

func (gui *GUI) getSelectedDeployment() entities.Deployment {
	selectedLine := gui.state.Panels.Deployments.SelectedIndex
	if selectedLine == -1 {
		return nil
	}
	return gui.cli.Deployments()[selectedLine]
}

func (gui *GUI) shouldRefresh(key string) bool {
	if gui.state.Panels.Main.Key == key {
		return false
	}
	gui.state.Panels.Main.Key = key
	return true
}
