package tui

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/alexeyco/simpletable"
	"github.com/jerson/deployer/pkg/deployer"
	"github.com/jerson/deployer/pkg/entities"
	"github.com/jerson/deployer/pkg/tui/common"
	"github.com/jerson/deployer/pkg/tui/layout"
	"github.com/jesseduffield/gocui"
	"github.com/xeonx/timeago"
)

func (gui *GUI) onMainTabChange(i int) error {
	log := gui.log.WithField(common.LogFieldMethod, "onMainTabChange").WithField("tabIndex", i)

	gui.g.Update(func(g *gocui.Gui) error {

		view, err := gui.getMainView()
		if err != nil {
			return err
		}
		if i == view.TabIndex {
			log.Trace("skipped")
			return nil
		}
		view.TabIndex = i
		gui.state.Panels.Main.TabIndex = view.TabIndex

		log.Debug("changed")
		return gui.updateMainTabs()
	})
	return nil

}

func (gui *GUI) updateMain() error {
	gui.g.Update(func(g *gocui.Gui) error {
		switch gui.state.Panels.Main.TabIndex {
		case 0:
			return gui.renderMainInfo()
		}
		return nil
	})
	return nil
}

func (gui *GUI) updateMainTabs() error {

	// Should be cleared here because is called from onChangeDeployment
	view, err := gui.getMainView()
	if err != nil {
		return err
	}
	_ = view.SetOrigin(0, 0)

	switch gui.state.Panels.Main.TabIndex {
	case 0:
		return gui.renderMainInfo()
	case 1:
		return gui.renderMainConsole()
	}
	return nil
}

func (gui *GUI) renderMainConsole() error {
	log := gui.log.WithField(common.LogFieldMethod, "renderMainConsole")

	view, err := gui.getMainView()
	if err != nil {
		return err
	}
	view.Clear()

	deployment := gui.getSelectedDeployment()
	if deployment == nil {
		log.Warn("deployment not selected")
		return nil
	}

	key := deployer.FunctionName(deployment)
	log = log.WithField("deployment", key)

	state := gui.state.Panels.Main.Console[key]
	if state == nil {
		log.Warn("console state not found")
		return nil
	}
	data := []byte(state.Output.String())
	tee := bytes.NewReader(data)

	go func() {
		total, err := io.Copy(view, tee)
		if err != nil {
			log.Error(err)
			return
		}
		log.Infof("copied to view: %v", total)
	}()

	view.Autoscroll = true
	return nil
}

func (gui *GUI) renderMainInfo() error {

	view, err := gui.getMainView()
	if err != nil {
		return err
	}
	view.Clear()

	deployment := gui.getSelectedDeployment()
	if deployment == nil {
		_, _ = gui.bufferErrorBanner(errors.New("first select a deployment")).WriteTo(view)
		return nil
	}
	_, _ = gui.bufferDeploymentResume(deployment).WriteTo(view)

	return nil
}

func (gui *GUI) bufferDeploymentResume(deployment entities.Deployment) *bytes.Buffer {

	w := &bytes.Buffer{}

	_, _ = gui.bufferTitle(deployment).WriteTo(w)

	ctx := context.Background()
	dependsOn, _ := deployment(ctx)
	if dependsOn != nil {
		_, _ = gui.bufferDependsOn("Depends on", dependsOn).WriteTo(w)
	}

	usedBy := gui.cli.UsedBy(deployment)
	if usedBy != nil {
		_, _ = gui.bufferUsedBy("Used by", usedBy).WriteTo(w)
	}

	status := gui.cli.DeploymentStatus(deployment)
	if status != nil {
		if status.Health != nil {
			_, _ = gui.bufferStatus("Health", status.Health).WriteTo(w)
		}
		if status.Install != nil {
			_, _ = gui.bufferStatus("Install", status.Install).WriteTo(w)
		}
		if status.Upgrade != nil {
			_, _ = gui.bufferStatus("Upgrade", status.Upgrade).WriteTo(w)
		}
		if status.Uninstall != nil {
			_, _ = gui.bufferStatus("Uninstall", status.Uninstall).WriteTo(w)
		}

	}
	return w
}

func (gui *GUI) bufferTitle(deployment entities.Deployment) *bytes.Buffer {

	var output bytes.Buffer

	output.WriteString("\n\n")
	output.WriteString(fmt.Sprintf("Name: %s\n", layout.Styles.Title.Sprint(deployer.FunctionName(deployment))))
	output.WriteString("\n")
	return &output
}

func (gui *GUI) bufferDependsOn(title string, dependencies []entities.Deployment) *bytes.Buffer {

	var output bytes.Buffer
	if len(dependencies) < 1 {
		return &output
	}

	output.WriteString(fmt.Sprintf("%s:\n", title))

	for _, deployment := range dependencies {
		output.WriteString(fmt.Sprintf("\t - %s\n", gui.deploymentRow(deployment)))
	}

	output.WriteString("\n")
	return &output
}

func (gui *GUI) bufferUsedBy(title string, dependencies []entities.Deployment) *bytes.Buffer {

	var output bytes.Buffer
	if len(dependencies) < 1 {
		return &output
	}

	output.WriteString(fmt.Sprintf("%s:\n", title))

	for _, deployment := range dependencies {
		output.WriteString(fmt.Sprintf("\t - %s\n", gui.deploymentRow(deployment)))
	}

	output.WriteString("\n")
	return &output
}

func (gui *GUI) bufferStatus(title string, status *entities.InstallerStatus) *bytes.Buffer {

	var output bytes.Buffer

	output.WriteString(fmt.Sprintf("%s:\n", title))
	output.WriteString(fmt.Sprintf("\tStatus: %v %s\n",
		gui.deploymentStatusIcon(status),
		gui.deploymentStatusName(status),
	))
	if status.StartTime != nil {
		output.WriteString(
			fmt.Sprintf("\tStart time: %s %s %s\n",
				status.StartTime.Format("2 Jan 2006 15:04:05"),
				layout.Styles.Disabled.Sprint("-"),
				timeago.English.Format(*status.StartTime),
			),
		)
	}
	if status.EndTime != nil {
		output.WriteString(
			fmt.Sprintf("\tEnd time: %s %s %s\n",
				status.StartTime.Format("2 Jan 2006 15:04:05"),
				layout.Styles.Disabled.Sprint("-"),
				timeago.English.Format(*status.StartTime),
			),
		)
	}
	if status.Completed() {
		output.WriteString(fmt.Sprintf("\tTook: %s\n", status.Took()))
	}
	if status.Error != nil {
		output.WriteString("\tError: \n")
		_, _ = gui.bufferErrorBanner(status.Error).WriteTo(&output)
		output.WriteString("\n")
	}
	output.WriteString("\n\n")
	return &output
}

func (gui *GUI) bufferErrorBanner(err error) *bytes.Buffer {
	log := gui.log.WithField(common.LogFieldMethod, "bufferErrorBanner")

	if err == nil {
		log.Warn("empty error")
		return nil
	}
	var output bytes.Buffer
	table := simpletable.New()

	row := []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Text: layout.Styles.Error.Sprint(err.Error())},
	}

	table.Body.Cells = append(table.Body.Cells, row)

	table.SetStyle(simpletable.StyleRounded)

	output.WriteString(table.String())
	return &output
}
