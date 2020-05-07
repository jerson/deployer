package entities

import "time"

// InstallerStatusType ...
type InstallerStatusType int

const (
	// InstallerStatusRunning ...
	InstallerStatusRunning InstallerStatusType = iota
	// InstallerStatusError ...
	InstallerStatusError
	// InstallerStatusIdle ...
	InstallerStatusIdle
	// InstallerStatusOff ...
	InstallerStatusOff
)

// InstallerStatus ...
type InstallerStatus struct {
	Error     error
	StartTime *time.Time
	EndTime   *time.Time
}

// NewInstallerStatus ...
func NewInstallerStatus() *InstallerStatus {
	return &InstallerStatus{
		Error:     nil,
		StartTime: nil,
		EndTime:   nil,
	}
}

// Reset ...
func (i *InstallerStatus) Reset() {
	i.Error = nil
	i.StartTime = nil
	i.EndTime = nil
}

// Status ...
func (i *InstallerStatus) Status() InstallerStatusType {

	if i.Error != nil {
		return InstallerStatusError
	}
	if i.StartTime != nil && i.EndTime == nil {
		return InstallerStatusIdle
	}
	if i.StartTime != nil && i.EndTime != nil {
		return InstallerStatusRunning
	}

	return InstallerStatusOff
}

// Start ...
func (i *InstallerStatus) Start() {
	now := time.Now()
	i.StartTime = &now
}

// End ...
func (i *InstallerStatus) End(err error) {
	now := time.Now()
	i.EndTime = &now
	i.Error = err
}

// Completed ...
func (i *InstallerStatus) Completed() bool {
	return i.EndTime != nil && i.StartTime != nil
}

// Took ...
func (i *InstallerStatus) Took() time.Duration {
	if i.EndTime == nil || i.StartTime == nil {
		return -1
	}
	return i.EndTime.Sub(*i.StartTime)
}
