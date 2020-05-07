package state

// PanelsState ...
type PanelsState struct {
	Project     *PanelState
	Engine      *PanelState
	Options     *PanelState
	Deployments *PanelState
	Main        *MainState
	Status      *PanelState
}

// State ...
type State struct {
	Panels     *PanelsState
	Background *PanelState
}

// NewState ...
func NewState() *State {
	return &State{
		Background: NewPanelState(nil),
		Panels: &PanelsState{
			Project:     NewPanelState(nil),
			Engine:      NewPanelState(nil),
			Options:     NewPanelState(nil),
			Deployments: NewPanelState(nil),
			Main:        NewMainState(),
			Status:      NewPanelState(nil),
		},
	}
}
