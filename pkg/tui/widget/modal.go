package widget

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jerson/deployer/pkg/tui/common"
	"github.com/jerson/deployer/pkg/tui/focus"
	"github.com/jerson/deployer/pkg/tui/layout"

	"github.com/jesseduffield/gocui"
	"github.com/sirupsen/logrus"
)

// Modal ...
type Modal struct {
	g            *gocui.Gui
	log          *logrus.Entry
	viewName     string
	previousView *gocui.View
	focus        *focus.Manager
}

// NewModal ...
func NewModal(name string, g *gocui.Gui, focus *focus.Manager, log *logrus.Entry, previousView *gocui.View) *Modal {
	return &Modal{
		g:            g,
		log:          log.WithField(common.LogFieldWidget, "modal"),
		viewName:     name,
		previousView: previousView,
		focus:        focus,
	}
}

// ShowDefault ...
func (m *Modal) ShowDefault(title, prompt string, handleConfirm, handleClose func(g *gocui.Gui, v *gocui.View) error) error {
	return m.Show(title, prompt, false, handleConfirm, handleClose)
}

// ShowConsole ...
func (m *Modal) ShowConsole(title string, handleClose func(*gocui.Gui, *gocui.View) error) error {
	defer m.log.WithField(common.LogFieldView, m.viewName).WithField(common.LogFieldMethod, "ShowConsole").Debug(title)
	m.onShow()
	m.g.Update(func(g *gocui.Gui) error {
		// delete the existing confirmation panel if it exists
		if view, _ := g.View(m.viewName); view != nil {
			if err := m.Close(); err != nil {
				m.log.Error(err)
			}
		}
		x0, y0, x1, y1 := m.getConsoleDimensions()
		view, err := m.g.SetView(m.viewName, x0, y0, x1, y1, 0)
		if err != nil {
			if err.Error() != "unknown view" {
				return err
			}
		}
		m.g.StartTicking()
		if view == nil {
			return errors.New("view nil")
		}
		view.Title = title
		view.Wrap = true
		view.FgColor = gocui.ColorWhite
		view.Autoscroll = true
		m.g.Update(func(g *gocui.Gui) error {
			return m.focus.Change(m.previousView, view, false)
		})
		view.Editable = false
		return m.setKeyBindings(handleClose, handleClose)
	})
	return nil
}

// Show ...
func (m *Modal) Show(title, prompt string, hasLoader bool, handleConfirm, handleClose func(g *gocui.Gui, v *gocui.View) error) error {
	defer m.log.WithField(common.LogFieldView, m.viewName).Debug("Show:", title)
	m.onShow()
	m.g.Update(func(g *gocui.Gui) error {
		// delete the existing confirmation panel if it exists
		if view, _ := g.View(m.viewName); view != nil {
			if err := m.Close(); err != nil {
				m.log.Error(err)
			}
		}
		view, err := m.prepare(title, prompt, hasLoader)
		if err != nil {
			return err
		}
		view.Editable = false
		if err := m.renderString(prompt); err != nil {
			return err
		}
		return m.setKeyBindings(handleConfirm, handleClose)
	})
	return nil
}

// ShowPrompt ...
func (m *Modal) ShowPrompt(title string, handleConfirm func(g *gocui.Gui, v *gocui.View) error) error {
	defer m.log.WithField(common.LogFieldView, m.viewName).Debug("ShowPrompt:", title)
	m.onShow()
	view, err := m.prepare(title, "", false)
	if err != nil {
		return err
	}
	view.Editable = true
	return m.setKeyBindings(handleConfirm, nil)
}

// ShowErrorWithFocus ...
func (m *Modal) ShowErrorWithFocus(message string, _ *gocui.View) error {
	coloredMessage := layout.Styles.Error.Sprint(strings.TrimSpace(message))
	return m.ShowDefault("Error", coloredMessage, nil, nil)
}

// ShowError ...
func (m *Modal) ShowError(message string) error {
	return m.ShowErrorWithFocus(message, m.g.CurrentView())
}

// Close ...
func (m *Modal) Close() error {
	log := m.log.WithField(common.LogFieldView, m.viewName).WithField(common.LogFieldMethod, "Close")
	view, err := m.g.View(m.viewName)
	if err != nil {
		log.Warn("view not found")
		return nil
	}
	if err := m.focus.ReturnFocus(view); err != nil {
		return err
	}
	m.g.DeleteKeybindings(m.viewName)
	log.Debug("closed")
	return m.g.DeleteView(m.viewName)
}

func (m *Modal) prepare(title, prompt string, hasLoader bool) (*gocui.View, error) {
	x0, y0, x1, y1 := m.getDimensions(true, prompt)
	view, err := m.g.SetView(m.viewName, x0, y0, x1, y1, 0)
	if err != nil {
		if err.Error() != "unknown view" {
			return nil, err
		}
	}
	if view == nil {
		return view, errors.New("view nil")
	}
	view.HasLoader = hasLoader
	m.g.StartTicking()
	view.Title = title
	view.Wrap = true
	view.FgColor = gocui.ColorWhite
	m.g.Update(func(g *gocui.Gui) error {
		return m.focus.Change(m.previousView, view, false)
	})
	return view, nil
}

func (m *Modal) wrapCallback(function func(*gocui.Gui, *gocui.View) error) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		err := m.Close()
		if err != nil {
			return err
		}
		if function != nil {
			if err := function(g, v); err != nil {
				return err
			}
		}
		return nil
	}
}

func (m *Modal) getMessageHeight(wrap bool, message string, width int) int {
	lines := strings.Split(message, "\n")
	lineCount := 0
	// if we need to wrap, calculate height to fit content within previousView's width
	if wrap {
		for _, line := range lines {
			lineCount += len(line)/width + 1
		}
	} else {
		lineCount = len(lines)
	}
	return lineCount
}

func (m *Modal) getDimensions(wrap bool, prompt string) (int, int, int, int) {
	width, height := m.g.Size()
	panelWidth := width / 2
	panelHeight := m.getMessageHeight(wrap, prompt, panelWidth)
	return width/2 - panelWidth/2,
		height/2 - panelHeight/2 - panelHeight%2 - 1,
		width/2 + panelWidth/2,
		height/2 + panelHeight/2
}

func (m *Modal) getConsoleDimensions() (int, int, int, int) {
	width, height := m.g.Size()
	panelWidth := int(float32(width) / 1.4)
	panelHeight := int(float32(height) / 1.4)
	return width/2 - panelWidth/2,
		height/2 - panelHeight/2 - panelHeight%2 - 1,
		width/2 + panelWidth/2,
		height/2 + panelHeight/2
}

func (m *Modal) onShow() {
	//viewNames := []string{"commitMessage",
	//	"credentials",
	//	"menu"}
	//for _, viewName := range viewNames {
	//	_, _ = m.g.SetViewOnBottom(viewName)
	//}
}

func (m *Modal) renderString(s string) error {
	m.g.Update(func(*gocui.Gui) error {
		log := m.log.WithField(common.LogFieldView, m.viewName).WithField(common.LogFieldMethod, "renderString")
		v, err := m.g.View(m.viewName)
		if err != nil {
			log.Warn("view not found")
			return nil
		}
		if err := v.SetOrigin(0, 0); err != nil {
			return err
		}
		if err := v.SetCursor(0, 0); err != nil {
			return err
		}

		v.Clear()
		_, err = fmt.Fprint(v, s)
		return err
	})
	return nil
}

func (m *Modal) setKeyBindings(handleConfirm, handleClose func(g *gocui.Gui, v *gocui.View) error) error {

	defer m.log.WithField(common.LogFieldView, m.viewName).WithField(common.LogFieldMethod, "setKeyBindings").Debug("configured")

	err := m.g.SetKeybinding(m.viewName, nil, gocui.KeyEnter, gocui.ModNone, m.wrapCallback(handleConfirm))
	if err != nil {
		return err
	}
	err = m.g.SetKeybinding(m.viewName, nil, 'y', gocui.ModNone, m.wrapCallback(handleConfirm))
	if err != nil {
		return err
	}
	err = m.g.SetKeybinding(m.viewName, nil, gocui.KeyEsc, gocui.ModNone, m.wrapCallback(handleClose))
	if err != nil {
		return err
	}
	err = m.g.SetKeybinding(m.viewName, nil, 'n', gocui.ModNone, m.wrapCallback(handleClose))
	if err != nil {
		return err
	}

	return nil
}
