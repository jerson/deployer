package common

import (
	"github.com/jesseduffield/gocui"
)

// ConditionalViewWriter ...
type ConditionalViewWriter struct {
	g         *gocui.Gui
	viewName  string
	condition ConditionalFunction
}

// NewConditionalViewWriter ...
func NewConditionalViewWriter(g *gocui.Gui, viewName string, condition ConditionalFunction) *ConditionalViewWriter {
	return &ConditionalViewWriter{
		g:         g,
		viewName:  viewName,
		condition: condition,
	}
}
func (c *ConditionalViewWriter) Write(p []byte) (n int, err error) {

	result, err := c.condition()
	if err != nil {
		return 0, err
	}
	if !result {
		return len(p), nil
	}

	view, err := c.g.View(c.viewName)
	if err != nil {
		return len(p), nil
	}

	return view.Write(p)

}
