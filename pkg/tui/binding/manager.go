package binding

import (
	"github.com/jerson/deployer/pkg/tui/common"
	"github.com/jesseduffield/gocui"
	"github.com/sirupsen/logrus"
)

// DefaultBinding ...
type DefaultBinding map[string]struct {
	OnKeyUpPress   func(*gocui.Gui, *gocui.View) error
	OnKeyDownPress func(*gocui.Gui, *gocui.View) error
	OnClick        func(*gocui.Gui, *gocui.View) error
	OnScrollUp     func(*gocui.Gui, *gocui.View) error
	OnScrollDown   func(*gocui.Gui, *gocui.View) error
	OnGoToEnd      func(*gocui.Gui, *gocui.View) error
	OnWrap         func(*gocui.Gui, *gocui.View) error
}

// Manager ...
type Manager struct {
	g              *gocui.Gui
	log            *logrus.Entry
	isModal        func(viewName string) bool
	customBindings DefaultBinding
}

// NewManager ...
func NewManager(g *gocui.Gui, log *logrus.Entry, isModal func(viewName string) bool) *Manager {
	return &Manager{
		g:              g,
		log:            log,
		isModal:        isModal,
		customBindings: DefaultBinding{},
	}
}

func (m *Manager) defaultBindings() []*Binding {

	return []*Binding{
		{
			viewName: "",
			Key:      'q',
			Modifier: gocui.ModNone,
			Handler:  m.quit,
		},
		{
			viewName: "",
			Key:      gocui.KeyCtrlC,
			Modifier: gocui.ModNone,
			Handler:  m.quit,
		},
		{
			viewName: "",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  m.quit,
		},
	}

}

// SetBinding ...
func (m *Manager) SetBinding(bindings DefaultBinding) {
	m.customBindings = bindings
}

// Register ...
func (m *Manager) Register() error {

	bindings := m.defaultBindings()
	for viewName, functions := range m.customBindings {
		if functions.OnScrollUp != nil {
			bindings = append(bindings,
				&Binding{viewName: viewName, Key: gocui.MouseWheelUp, Modifier: gocui.ModNone, Handler: functions.OnScrollUp},
				&Binding{viewName: viewName, Key: gocui.KeyPgup, Modifier: gocui.ModNone, Handler: functions.OnScrollUp},
			)

		}
		if functions.OnScrollDown != nil {
			bindings = append(bindings,
				&Binding{viewName: viewName, Key: gocui.MouseWheelDown, Modifier: gocui.ModNone, Handler: functions.OnScrollDown},
				&Binding{viewName: viewName, Key: gocui.KeyPgdn, Modifier: gocui.ModNone, Handler: functions.OnScrollDown},
			)
		}
		if functions.OnGoToEnd != nil {
			bindings = append(bindings,
				&Binding{viewName: viewName, Key: gocui.KeyEnd, Modifier: gocui.ModNone, Handler: functions.OnGoToEnd},
			)
		}
		if functions.OnWrap != nil {
			bindings = append(bindings,
				&Binding{viewName: viewName, Key: gocui.KeyCtrlW, Modifier: gocui.ModNone, Handler: functions.OnWrap},
			)
		}
		if functions.OnKeyUpPress != nil {
			bindings = append(bindings,
				&Binding{viewName: viewName, Key: 'k', Modifier: gocui.ModNone, Handler: functions.OnKeyUpPress},
				&Binding{viewName: viewName, Key: gocui.KeyArrowUp, Modifier: gocui.ModNone, Handler: functions.OnKeyUpPress},
			)
			if functions.OnScrollUp == nil {
				bindings = append(bindings,
					&Binding{viewName: viewName, Key: gocui.MouseWheelUp, Modifier: gocui.ModNone, Handler: functions.OnKeyUpPress},
				)
			}
		}
		if functions.OnKeyDownPress != nil {
			bindings = append(bindings,
				&Binding{viewName: viewName, Key: 'j', Modifier: gocui.ModNone, Handler: functions.OnKeyDownPress},
				&Binding{viewName: viewName, Key: gocui.KeyArrowDown, Modifier: gocui.ModNone, Handler: functions.OnKeyDownPress},
			)
			if functions.OnScrollDown == nil {
				bindings = append(bindings,
					&Binding{viewName: viewName, Key: gocui.MouseWheelDown, Modifier: gocui.ModNone, Handler: functions.OnKeyDownPress},
				)
			}
		}
		if functions.OnClick != nil {
			bindings = append(bindings,
				&Binding{viewName: viewName, Key: gocui.MouseLeft, Modifier: gocui.ModNone, Handler: functions.OnClick},
			)
		}

	}

	log := m.log.WithField(common.LogFieldMethod, "Register")
	for _, binding := range bindings {
		log.Debugf("%v %v %v", binding.viewName, binding.GetDisplayStrings(false), binding.Modifier)
		err := m.g.SetKeybinding(binding.viewName, nil, binding.Key, binding.Modifier, binding.Handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func (m *Manager) currentViewName() string {
	currentView := m.g.CurrentView()
	if currentView == nil {
		return ""
	}
	return currentView.Name()
}
