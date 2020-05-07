package state

import (
	"context"
	"sync"

	"github.com/jerson/deployer/pkg/deployer"
	"github.com/jerson/deployer/pkg/formatter"
)

// PanelState ...
type PanelState struct {
	Ctx           context.Context
	SelectedIndex int
	Key           string
	Logger        *LoggerState
	settings      map[string]interface{}
	sync.RWMutex
}

// NewPanelState ...
func NewPanelState(format *formatter.Console) *PanelState {
	instance := &PanelState{
		SelectedIndex: -1,
		Key:           "",
		settings:      map[string]interface{}{},
		Ctx:           context.Background(),
		Logger:        NewLoggerState(format),
	}

	instance.Logger.Log.SetOutput(instance.Logger.Output)
	instance.Ctx = context.WithValue(instance.Ctx, deployer.ContextKeyLog, instance.Logger.Log.WithFields(nil))

	return instance
}

// SetSetting ...
func (p *PanelState) SetSetting(name string, value interface{}) {
	p.Lock()
	defer p.Unlock()
	p.settings[name] = value

}

// GetSetting ...
func (p *PanelState) GetSetting(name string) interface{} {
	p.RLock()
	defer p.RUnlock()
	return p.settings[name]
}

// GetSettingWithDefault ...
func (p *PanelState) GetSettingWithDefault(name string, defaultValue interface{}) interface{} {
	p.Lock()
	defer p.Unlock()
	if p.settings[name] == nil {
		p.settings[name] = defaultValue
	}
	return p.settings[name]
}
