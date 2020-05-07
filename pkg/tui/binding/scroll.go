package binding

import (
	"math"

	"github.com/jerson/deployer/pkg/tui/common"
	"github.com/jesseduffield/gocui"
)

const scrollOffset = 2
const scrollPastBottom = true

// ScrollUp ...
func (m *Manager) ScrollUp(view *gocui.View) error {
	defer func() {
		x, y := view.Origin()
		m.log.WithField(common.LogFieldView, view.Name()).WithField(common.LogFieldMethod, "ScrollUp").Tracef("(%v,%v)", x, y)
	}()

	view.Autoscroll = false
	ox, oy := view.Origin()
	newOy := int(math.Max(0, float64(oy-scrollOffset)))
	return view.SetOrigin(ox, newOy)
}

// ScrollDown ...
func (m *Manager) ScrollDown(view *gocui.View) error {
	log := m.log.WithField(common.LogFieldView, view.Name()).WithField(common.LogFieldMethod, "ScrollDown")
	defer func() {
		x, y := view.Origin()
		log.Tracef("(%v,%v)", x, y)
	}()
	view.Autoscroll = false
	ox, oy := view.Origin()

	reservedLines := 0
	if scrollPastBottom {
		_, sizeY := view.Size()
		reservedLines = sizeY
	}

	totalLines := view.ViewLinesHeight()
	if oy+reservedLines >= totalLines {
		log.Warn("not more to scroll")
		return nil
	}

	return view.SetOrigin(ox, oy+scrollOffset)
}

// ScrollLeft ...
func (m *Manager) ScrollLeft(view *gocui.View) error {
	log := m.log.WithField(common.LogFieldView, view.Name()).WithField(common.LogFieldMethod, "ScrollLeft")
	defer func() {
		x, y := view.Origin()
		log.Tracef("(%v,%v)", x, y)
	}()
	ox, oy := view.Origin()
	newOx := int(math.Max(0, float64(ox-scrollOffset)))

	return view.SetOrigin(newOx, oy)
}

// ScrollRight ...
func (m *Manager) ScrollRight(view *gocui.View) error {
	log := m.log.WithField(common.LogFieldView, view.Name()).WithField(common.LogFieldMethod, "ScrollRight")
	defer func() {
		x, y := view.Origin()
		log.Tracef("(%v,%v)", x, y)
	}()
	ox, oy := view.Origin()

	content := view.ViewBufferLines()
	var largestNumberOfCharacters int
	for _, txt := range content {
		if len(txt) > largestNumberOfCharacters {
			largestNumberOfCharacters = len(txt)
		}
	}

	sizeX, _ := view.Size()
	if ox+sizeX >= largestNumberOfCharacters {
		log.Warn("not more to scroll")
		return nil
	}

	return view.SetOrigin(ox+scrollOffset, oy)
}

// AutoScroll ...
func (m *Manager) AutoScroll(view *gocui.View) error {
	defer func() {
		m.log.WithField("view", view.Name()).WithField("method", "AutoScroll").Trace("updated")
	}()
	view.Autoscroll = true

	return nil
}

// OnWrap ...
func (m *Manager) OnWrap(view *gocui.View) error {
	defer func() {
		m.log.WithField("view", view.Name()).WithField("method", "OnWrap").Trace("toggled")
	}()
	view.Wrap = !view.Wrap

	return nil
}
