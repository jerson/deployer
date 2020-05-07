package binding

import (
	"github.com/jesseduffield/gocui"
)

// Binding ...
type Binding struct {
	viewName    string
	Handler     func(*gocui.Gui, *gocui.View) error
	Key         interface{}
	Modifier    gocui.Modifier
	Description string
}

// GetDisplayStrings ...
func (b *Binding) GetDisplayStrings(_ bool) []string {
	return []string{b.GetKey(), b.Description}
}

// GetKey ...
func (b *Binding) GetKey() string {
	key := 0

	switch b.Key.(type) {
	case rune:
		key = int(b.Key.(rune))
	case gocui.Key:
		key = int(b.Key.(gocui.Key))
	}

	switch key {
	case 27:
		return "esc"
	case 13:
		return "enter"
	case 32:
		return "space"
	case 65514:
		return "►"
	case 65515:
		return "◄"
	case 65517:
		return "▲"
	case 65516:
		return "▼"
	case 65508:
		return "PgUp"
	case 65507:
		return "PgDn"
	}

	return string(key)
}
