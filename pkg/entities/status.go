package entities

// Status ...
type Status struct {
	Health    *InstallerStatus
	Install   *InstallerStatus
	Upgrade   *InstallerStatus
	Uninstall *InstallerStatus
}

// NewStatus ...
func NewStatus() *Status {
	return &Status{
		Health:    NewInstallerStatus(),
		Install:   NewInstallerStatus(),
		Upgrade:   NewInstallerStatus(),
		Uninstall: NewInstallerStatus(),
	}
}
