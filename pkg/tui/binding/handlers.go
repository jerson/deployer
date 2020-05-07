package binding

import (
	"errors"

	"github.com/jerson/deployer/pkg/tui/common"
	"github.com/jesseduffield/gocui"
)

// OnClick ...
func (m *Manager) OnClick(v *gocui.View, selectedLine *int, itemCount int) error {

	log := m.log.WithField(common.LogFieldView, v.Name()).WithField(common.LogFieldMethod, "OnClick")

	if m.isModal(m.currentViewName()) && v != nil {
		log.Debug("skipped: isModal")
		return nil
	}
	if v == nil {
		return errors.New("view nil")
	}

	if _, err := m.g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	newSelectedLine := v.SelectedLineIdx()

	if newSelectedLine < 0 {
		newSelectedLine = 0
	}

	if newSelectedLine > itemCount-1 {
		newSelectedLine = itemCount - 1
	}

	*selectedLine = newSelectedLine

	log.Debugf("(%v/%v)", *selectedLine, itemCount)

	return nil
}

// OnChangeSelectedLine ...
func (m *Manager) OnChangeSelectedLine(v *gocui.View, line *int, total int, up bool) error {
	defer func() {
		m.log.WithField(common.LogFieldView, v.Name()).WithField(common.LogFieldMethod, "OnChangeSelectedLine").Debugf("(%v/%v) up:%v", *line, total, up)
	}()
	if up {
		if *line == -1 || *line == 0 {
			return nil
		}
		*line--
	} else {
		if *line == -1 || *line == total-1 {
			return nil
		}
		*line++
	}
	return nil
}
