package focus

import (
	"github.com/jerson/deployer/pkg/tui/common"
	"github.com/jerson/deployer/pkg/tui/utils"

	"github.com/golang-collections/collections/stack"
	"github.com/jesseduffield/gocui"
	"github.com/sirupsen/logrus"
)

// Manager ...
type Manager struct {
	g            *gocui.Gui
	log          *logrus.Entry
	mainViewName string
	views        *stack.Stack
	isModal      func(viewName string) bool
}

// NewManager ...
func NewManager(g *gocui.Gui, log *logrus.Entry, mainViewName string, isModal func(viewName string) bool) *Manager {
	return &Manager{
		g:            g,
		log:          log,
		mainViewName: mainViewName,
		isModal:      isModal,
		views:        stack.New(),
	}
}

// Layout returns a manager function for when view gain and lose focus
func (m *Manager) Layout() func(g *gocui.Gui) error {
	var previousView *gocui.View
	return func(g *gocui.Gui) error {
		newView := m.g.CurrentView()
		if err := m.onChange(); err != nil {
			return err
		}
		// for now we don't consider losing focus to a popup panel as actually losing focus
		if newView != previousView && !m.isModal(newView.Name()) {
			if err := m.onFocusLost(previousView, newView); err != nil {
				return err
			}
			if err := m.onFocus(newView); err != nil {
				return err
			}
			previousView = newView
		}
		return nil
	}
}

// Change ...
func (m *Manager) Change(oldView, newView *gocui.View, returning bool) error {

	log := m.log.WithField(common.LogFieldMethod, "Change").WithField("returning", returning)
	defer func() {
		oldName := ""
		if oldView != nil {
			oldName = oldView.Name()
		}
		log.Infof("%v => %v", oldName, newView.Name())
	}()

	// we assume we'll never want to return focus to a popup panel i.e.
	// we should never stack popup panels
	if oldView != nil && !m.isModal(oldView.Name()) && !returning {
		m.pushPreviousView(oldView.Name())
	}

	//m.Log.Info("setting highlight to true for view " + newView.Name())
	//m.Log.Info("new focused view is " + newView.Name())
	if _, err := m.g.SetCurrentView(newView.Name()); err != nil {
		return err
	}
	if _, err := m.g.SetViewOnTop(newView.Name()); err != nil {
		return err
	}

	m.g.Cursor = newView.Editable

	//if err := m.renderPanelOptions(); err != nil {
	//	return err
	//}

	//return m.newLineFocused(newView)
	return nil
}

func (m *Manager) onChange() error {
	currentView := m.g.CurrentView()
	for _, view := range m.g.Views() {
		view.Highlight = view == currentView && view.Name() != m.mainViewName
		//	view.Highlight = view == currentView
	}
	return nil
}

func (m *Manager) onFocusLost(v *gocui.View, newView *gocui.View) error {
	log := m.log.WithField(common.LogFieldMethod, "onFocusLost")
	if v == nil {
		log.Warn("view not found")
		return nil
	}

	if !m.isModal(newView.Name()) {
		v.ParentView = nil
	}

	log.Info(v.Name())

	return nil
}

func (m *Manager) onFocus(view *gocui.View) error {
	log := m.log.WithField(common.LogFieldMethod, "onFocus")
	if view == nil {
		log.Warn("view not found")
		return nil
	}
	log.Info(view.Name())
	return nil
}

func (m *Manager) pushPreviousView(name string) {
	m.views.Push(name)
}

// ReturnFocus ...
func (m *Manager) ReturnFocus(v *gocui.View) error {
	previousViewName := m.popPreviousView()
	previousView, err := m.g.View(previousViewName)
	if err != nil {
		previousView, err = m.g.View(m.mainViewName)
		if err != nil {
			return err
		}
	}
	return m.Change(v, previousView, true)
}

func (m *Manager) popPreviousView() string {
	if m.views.Len() > 0 {
		return m.views.Pop().(string)
	}
	return ""
}

func (m *Manager) peekPreviousView() string {
	if m.views.Len() > 0 {
		return m.views.Peek().(string)
	}
	return ""
}

// FocusPoint ...
func (m *Manager) FocusPoint(_ int, selectedY int, lineCount int, v *gocui.View) error {
	log := m.log.WithField(common.LogFieldView, v.Name()).WithField(common.LogFieldMethod, "FocusPoint")
	if selectedY < 0 || selectedY > lineCount {
		log.WithField("selectedY", selectedY).WithField("lineCount", lineCount).Warn("invalid range")
		return nil
	}
	ox, oy := v.Origin()
	originalOy := oy
	cx, cy := v.Cursor()
	originalCy := cy
	_, height := v.Size()

	ly := utils.Max(height-1, 0)

	windowStart := oy
	windowEnd := oy + ly

	if selectedY < windowStart {
		oy = utils.Max(oy-(windowStart-selectedY), 0)
	} else if selectedY > windowEnd {
		oy += selectedY - windowEnd
	}

	if windowEnd > lineCount-1 {
		shiftAmount := windowEnd - (lineCount - 1)
		oy = utils.Max(oy-shiftAmount, 0)
	}

	if originalOy != oy {
		_ = v.SetOrigin(ox, oy)
	}

	cy = selectedY - oy
	if originalCy != cy {
		_ = v.SetCursor(cx, selectedY-oy)
	}

	fx, fy := v.Cursor()
	log.Tracef("(%v,%v)", fx, fy)

	return nil
}
