package tui

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/gocui"
)

func (gui *GUI) layout(g *gocui.Gui) error {

	g.Highlight = true

	err := gui.layoutProject()
	if err != nil {
		return err
	}
	err = gui.layoutEngine()
	if err != nil {
		return err
	}
	err = gui.layoutOptions()
	if err != nil {
		return err
	}
	err = gui.layoutDeployments()
	if err != nil {
		return err
	}
	err = gui.layoutMain()
	if err != nil {
		return err
	}
	err = gui.layoutDebug()
	if err != nil {
		return err
	}
	err = gui.layoutStatus()
	if err != nil {
		return err
	}

	if gui.g.CurrentView() == nil {
		view, err := gui.g.View(gui.initiallyFocusedViewName())
		if err != nil {
			return err
		}
		err = gui.focus.Change(nil, view, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gui *GUI) initiallyFocusedViewName() string {
	return "deploy"
}
func (gui *GUI) layoutProject() error {
	maxX, _ := gui.g.Size()
	view, err := gui.g.SetView("project", 0, 0, int(0.22*float32(maxX)), 2, 0)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	view.Title = "Project"
	view.Highlight = true
	return nil
}

func (gui *GUI) layoutEngine() error {
	view, err := gui.g.SetViewBeneath("engine", "project", 7)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	view.Title = fmt.Sprintf("Minikube")
	view.Highlight = true
	return nil
}

func (gui *GUI) layoutOptions() error {
	view, err := gui.g.SetViewBeneath("options", "engine", 6)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	view.Title = "Deployment Options"
	view.Highlight = true

	return nil
}

func (gui *GUI) layoutDeployments() error {
	_, maxY := gui.g.Size()
	height := maxY - 17
	if gui.debug {
		height -= 15
	}
	view, err := gui.g.SetViewBeneath("deploy", "options", height)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	view.Title = "Deployments"
	view.Highlight = true
	return nil
}
func (gui *GUI) layoutMain() error {
	maxX, maxY := gui.g.Size()
	height := maxY - 2
	if gui.debug {
		height -= 15
	}
	view, err := gui.g.SetView("main", int(0.23*float32(maxX)), 0, maxX-1, height, 0)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	view.Tabs = []string{"Deployment", "Console"}
	view.Highlight = true
	return nil
}

func (gui *GUI) layoutDebug() error {
	maxX, maxY := gui.g.Size()
	y := maxY - 2
	if gui.debug {
		y -= 14
	}
	view, err := gui.g.SetView("debug", 0, y, maxX-1, maxY-2, 0)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	view.Highlight = true
	view.Title = "Debug console"
	if gui.debug {
		view.Frame = true
		view.Visible = true
	} else {
		view.Frame = false
		view.Visible = false
	}
	//view.BgColor = Red
	return nil
}

func (gui *GUI) layoutStatus() error {
	maxX, maxY := gui.g.Size()
	view, err := gui.g.SetView("status", 0, maxY-2, maxX-1, maxY-1, 0)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	view.Frame = false
	view.Highlight = true
	return nil
}
