package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jerson/deployer/pkg/tui/common"
	"github.com/jerson/deployer/pkg/tui/utils"
	"github.com/jesseduffield/gocui"
)

func (gui *GUI) setViewContent(_ *gocui.Gui, v *gocui.View, s string) error {
	v.Clear()
	_, err := fmt.Fprint(v, utils.CleanString(s))
	return err
}

// renderString resets the origin of a view and sets its content
func (gui *GUI) renderString(g *gocui.Gui, viewName, s string) error {
	g.Update(func(*gocui.Gui) error {
		log := gui.log.WithField(common.LogFieldView, viewName).WithField(common.LogFieldMethod, "renderString")
		v, err := g.View(viewName)
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
		return gui.setViewContent(gui.g, v, s)
	})
	return nil
}

// reRenderString sets the view's content, without changing its origin
func (gui *GUI) reRenderString(g *gocui.Gui, viewName, s string) error {
	g.Update(func(*gocui.Gui) error {
		log := gui.log.WithField(common.LogFieldView, viewName).WithField(common.LogFieldMethod, "reRenderString")
		v, err := g.View(viewName)
		if err != nil {
			log.Warn("view not found")
			return nil
		}
		return gui.setViewContent(gui.g, v, s)
	})
	return nil
}

func (gui *GUI) optionsMapToString(optionsMap map[string]string) string {
	optionsArray := make([]string, 0)
	for key, description := range optionsMap {
		optionsArray = append(optionsArray, key+": "+description)
	}
	sort.Strings(optionsArray)
	return strings.Join(optionsArray, ", ")
}

func (gui *GUI) trimmedContent(v *gocui.View) string {
	return strings.TrimSpace(v.Buffer())
}
