package tui

import (
	"github.com/jerson/deployer/pkg/tui/binding"
	"github.com/jesseduffield/gocui"
)

func (gui *GUI) registerBindings() error {

	gui.binding.SetBinding(gui.customBinding())
	err := gui.binding.Register()
	if err != nil {
		return err
	}

	err = gui.g.SetTabClickBinding("main", gui.onMainTabChange)
	if err != nil {
		return err
	}

	err = gui.g.SetKeybinding("", nil, gocui.KeyF1, gocui.ModNone, gui.toggleDebug)
	if err != nil {
		return err
	}
	return nil
}

func (gui *GUI) toggleDebug(_ *gocui.Gui, _ *gocui.View) error {
	gui.debug = !gui.debug
	return nil
}

func (gui *GUI) customBinding() binding.DefaultBinding {
	return binding.DefaultBinding{
		"main": {
			OnGoToEnd: func(g *gocui.Gui, v *gocui.View) error {
				return gui.binding.AutoScroll(v)
			},
			OnWrap: func(g *gocui.Gui, v *gocui.View) error {
				return gui.binding.OnWrap(v)
			},
			OnScrollUp: func(g *gocui.Gui, v *gocui.View) error {
				return gui.binding.ScrollUp(v)
			},
			OnScrollDown: func(g *gocui.Gui, v *gocui.View) error {
				return gui.binding.ScrollDown(v)
			},
			OnClick: func(g *gocui.Gui, v *gocui.View) error {
				state := gui.state.Panels.Main
				err := gui.onClick(v, &state.SelectedIndex, v.LinesHeight())
				if err != nil {
					return err
				}
				return gui.onChangeEngine()
			},
		},
		"debug": {
			OnGoToEnd: func(g *gocui.Gui, v *gocui.View) error {
				return gui.binding.AutoScroll(v)
			},
			OnWrap: func(g *gocui.Gui, v *gocui.View) error {
				return gui.binding.OnWrap(v)
			},
			OnScrollUp: func(g *gocui.Gui, v *gocui.View) error {
				return gui.binding.ScrollUp(v)
			},
			OnScrollDown: func(g *gocui.Gui, v *gocui.View) error {
				return gui.binding.ScrollDown(v)
			},
			OnClick: func(g *gocui.Gui, v *gocui.View) error {
				state := gui.state.Panels.Main
				err := gui.onClick(v, &state.SelectedIndex, v.LinesHeight())
				if err != nil {
					return err
				}
				return gui.onChangeEngine()
			},
		},
		"project": {
			OnClick: func(g *gocui.Gui, v *gocui.View) error {
				return gui.initProject()
			},
		},
		"options": {
			OnClick: func(g *gocui.Gui, v *gocui.View) error {
				state := gui.state.Panels.Options
				err := gui.onClick(v, &state.SelectedIndex, 4)
				if err != nil {
					return err
				}
				return gui.onChangeOption()
			},
		},
		"engine": {
			OnClick: func(g *gocui.Gui, v *gocui.View) error {
				state := gui.state.Panels.Engine
				err := gui.onClick(v, &state.SelectedIndex, 6)
				if err != nil {
					return err
				}
				return gui.onChangeEngine()
			},
		},
		"deployments": {
			OnKeyUpPress: func(g *gocui.Gui, v *gocui.View) error {
				state := gui.state.Panels.Deployments
				err := gui.onChangeSelectedLine(v, &state.SelectedIndex, len(gui.cli.Deployments()), true)
				if err != nil {
					return err
				}
				return gui.onChangeDeployment()
			},
			OnKeyDownPress: func(g *gocui.Gui, v *gocui.View) error {
				state := gui.state.Panels.Deployments
				err := gui.onChangeSelectedLine(v, &state.SelectedIndex, len(gui.cli.Deployments()), false)
				if err != nil {
					return err
				}
				return gui.onChangeDeployment()
			},
			OnClick: func(g *gocui.Gui, v *gocui.View) error {
				state := gui.state.Panels.Deployments
				err := gui.onClick(v, &state.SelectedIndex, len(gui.cli.Deployments()))
				if err != nil {
					return err
				}
				return gui.onChangeDeployment()
			},
		},
	}
}

func (gui *GUI) onChangeSelectedLine(v *gocui.View, line *int, total int, up bool) error {
	err := gui.binding.OnChangeSelectedLine(v, line, total, up)
	if err != nil {
		return err
	}
	v.FocusPoint(0, *line)
	return nil
}
func (gui *GUI) onClick(v *gocui.View, selectedLine *int, itemCount int) error {

	err := gui.binding.OnClick(v, selectedLine, itemCount)
	if err != nil {
		return err
	}
	v.FocusPoint(0, *selectedLine)
	return nil
}
