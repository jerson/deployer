package state

import "sync"

// MainState ...
type MainState struct {
	Console  map[string]*LoggerState
	TabIndex int
	PanelState
	sync.Mutex
}

// NewMainState ...
func NewMainState() *MainState {
	return &MainState{
		Console:    map[string]*LoggerState{},
		TabIndex:   0,
		PanelState: *NewPanelState(nil),
		Mutex:      sync.Mutex{},
	}
}

// GetConsoleLogger ...
func (m *MainState) GetConsoleLogger(name string) *LoggerState {
	m.Lock()
	defer m.Unlock()
	if m.Console[name] == nil {
		m.Console[name] = NewLoggerState(nil)
	}
	return m.Console[name]
}

// SetLogger ...
func (m *MainState) SetLogger(name string, logger *LoggerState) {
	m.Lock()
	defer m.Unlock()
	m.Console[name] = logger

}
