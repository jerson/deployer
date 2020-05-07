package layout

import (
	"github.com/fatih/color"
)

// TextStyles ...
type TextStyles struct {
	Default      *color.Color
	Button       *color.Color
	ButtonActive *color.Color
	Disabled     *color.Color
	Error        *color.Color
	Warning      *color.Color
	Success      *color.Color
	Title        *color.Color
}

// Styles ...
var Styles = TextStyles{
	Default:      color.New(color.FgHiWhite),
	Button:       color.New(color.Bold),
	ButtonActive: color.New(color.Bold, color.FgGreen),
	Disabled:     color.New(color.FgWhite),
	Error:        color.New(color.FgRed),
	Warning:      color.New(color.FgYellow),
	Success:      color.New(color.FgGreen),
	Title:        color.New(color.FgGreen),
}
