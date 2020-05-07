package tui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/enescakir/emoji"
	"github.com/jerson/deployer/modules/kubectl"
	"github.com/jerson/deployer/modules/minikube"
	"github.com/jerson/deployer/pkg/deployer"
	"github.com/jerson/deployer/pkg/formatter"
	"github.com/jerson/deployer/pkg/tui/common"
	"github.com/jerson/deployer/pkg/tui/layout"
	"github.com/jerson/deployer/pkg/tui/widget"
	"github.com/jesseduffield/gocui"
	"github.com/sirupsen/logrus"
)

func (gui *GUI) onChangeEngine() error {

	view, err := gui.getEngineView()
	if err != nil {
		return err
	}

	state := gui.state.Panels.Engine
	if state.SelectedIndex == 4 {
		var err error
		modal := widget.NewModal("engine_modal", gui.g, gui.focus, gui.log.WithFields(nil), view)
		running := state.GetSettingWithDefault("running", false).(bool)

		if running {
			err = modal.ShowDefault("Stop Minikube", "Are you sure?",
				func(g *gocui.Gui, view *gocui.View) error {
					return gui.stopEngine()
				},
				nil,
			)
		} else {
			err = modal.ShowDefault("Start Minikube", "Are you sure?",
				func(g *gocui.Gui, view *gocui.View) error {
					return gui.startEngine()
				},
				nil,
			)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (gui *GUI) startEngine() error {

	view, err := gui.getEngineView()
	if err != nil {
		return err
	}

	viewModalName := "engine_modal_start"
	modal := widget.NewModal(viewModalName, gui.g, gui.focus, gui.log.WithFields(nil), view)
	err = modal.ShowConsole("Starting Minikube", nil)
	if err != nil {
		return err
	}

	gui.g.Update(func(g *gocui.Gui) error {
		viewModal, err := gui.g.View(viewModalName)
		if err != nil {
			return err
		}

		log := common.CustomLogger(viewModal, &formatter.Console{}).WithField(common.LogFieldEngine, "minikube")
		ctx := context.WithValue(context.Background(), deployer.ContextKeyLog, log)
		go func() {
			_, _ = minikube.Run(ctx, "start")
			log.Info("finished")
			_ = gui.countDown(3, log, func() error {
				return modal.Close()
			})
		}()
		return nil
	})

	return nil
}

func (gui *GUI) stopEngine() error {

	view, err := gui.getEngineView()
	if err != nil {
		return err
	}

	viewModalName := "engine_modal_stop"
	modal := widget.NewModal(viewModalName, gui.g, gui.focus, gui.log.WithFields(nil), view)
	err = modal.ShowConsole("Stopping Minikube", nil)
	if err != nil {
		return err
	}

	gui.g.Update(func(g *gocui.Gui) error {
		viewModal, err := gui.g.View(viewModalName)
		if err != nil {
			return err
		}

		log := common.CustomLogger(viewModal, &formatter.Console{}).WithField(common.LogFieldEngine, "minikube")
		ctx := context.WithValue(context.Background(), deployer.ContextKeyLog, log)
		go func() {
			_, _ = minikube.Run(ctx, "stop")
			log.Info("finished")
			_ = gui.countDown(3, log, func() error {
				return modal.Close()
			})
		}()

		return nil
	})

	return nil
}

func (gui *GUI) countDown(times int, log *logrus.Entry, callback func() error) error {

	for i := times; i >= 0; i-- {
		time.Sleep(time.Second)
		log.Warn(i)
	}

	return callback()
}

func (gui *GUI) updateEngine() error {

	state := gui.state.Panels.Engine
	buffer := gui.bufferEngine(state.Ctx)

	gui.g.Update(func(g *gocui.Gui) error {
		view, err := gui.getEngineView()
		if err != nil {
			return err
		}

		view.Clear()

		_, err = buffer.WriteTo(view)
		if err != nil {
			return err
		}

		return nil
	})
	return nil
}

func (gui *GUI) bufferEngine(ctx context.Context) *bytes.Buffer {

	var output bytes.Buffer
	status := ""

	state := gui.state.Panels.Engine
	running := false
	ip, err := minikube.Run(ctx, "ip")
	if err != nil {
		ip = "-"
		status = layout.Styles.Error.Sprint("NOT RUNNING")
		state.SetSetting("running", false)
	} else {
		running = true
		status = layout.Styles.Success.Sprint("RUNNING")
		state.SetSetting("running", true)
	}

	output.WriteString(fmt.Sprintf("Status: %s\n", status))
	versionString, _ := minikube.Run(ctx, "version", "-o", "json")
	version := MinikubeVersion{
		MinikubeVersion: "NOT FOUND",
	}
	_ = json.Unmarshal([]byte(versionString), &version)

	output.WriteString(fmt.Sprintf("Version: %s\n", version.MinikubeVersion))

	versionKubeCTLString, _ := kubectl.Run(ctx, "version", "-o", "json")
	versionKubeCTL := KubeCTLVersion{
		ClientVersion: KubeCTLModeVersion{
			GitVersion: "NOT FOUND",
		},
		ServerVersion: KubeCTLModeVersion{
			GitVersion: "NOT FOUND",
		},
	}
	_ = json.Unmarshal([]byte(versionKubeCTLString), &versionKubeCTL)

	output.WriteString(fmt.Sprintf("Kubernetes: %s\n", versionKubeCTL.ServerVersion.GitVersion))
	output.WriteString(fmt.Sprintf("IP: %s\n", ip))

	if running {
		_, _ = gui.bufferMenuItem(emoji.StopButton, " Stop Minikube").WriteTo(&output)
	} else {
		_, _ = gui.bufferMenuItem(emoji.PlayButton, " Start Minikube").WriteTo(&output)
	}

	output.WriteString("\n")
	return &output
}

// MinikubeVersion ...
type MinikubeVersion struct {
	Commit          string
	MinikubeVersion string
}

// KubeCTLVersion ...
type KubeCTLVersion struct {
	ClientVersion KubeCTLModeVersion
	ServerVersion KubeCTLModeVersion
}

// KubeCTLModeVersion ...
type KubeCTLModeVersion struct {
	GitVersion string
}
