package tui

import (
	"github.com/jesseduffield/gocui"
)

func (gui *GUI) getProjectView() (*gocui.View, error) {
	return gui.g.View("project")
}

func (gui *GUI) getEngineView() (*gocui.View, error) {
	return gui.g.View("engine")
}
func (gui *GUI) getOptionsView() (*gocui.View, error) {
	return gui.g.View("options")
}

func (gui *GUI) getMainView() (*gocui.View, error) {
	return gui.g.View("main")
}

func (gui *GUI) getDeploymentsView() (*gocui.View, error) {
	return gui.g.View("deployments")
}

func (gui *GUI) getDebugView() (*gocui.View, error) {
	return gui.g.View("debug")
}

func (gui *GUI) getStatusView() (*gocui.View, error) {
	return gui.g.View("status")
}
