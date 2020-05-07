package tui

import (
	"fmt"
	"github.com/jesseduffield/gocui"
)

func (gui *GUI) initProject() error {
	gui.g.Update(func(g *gocui.Gui) error {
		view, err := gui.getProjectView()
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(view, gui.config.Name)
		return err
	})
	return nil
}
