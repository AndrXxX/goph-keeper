package keyboard

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Full  [][]key.Binding
	Short []key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return k.Short
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return k.Full
}

type KeyList struct {
	New    key.Binding
	Edit   key.Binding
	Delete key.Binding
	Up     key.Binding
	Down   key.Binding
	Right  key.Binding
	Left   key.Binding
	Enter  key.Binding
	Help   key.Binding
	Quit   key.Binding
	Back   key.Binding
	Copy   key.Binding
	Save   key.Binding
}

var Keys = KeyList{
	New:    New,
	Edit:   Edit,
	Delete: Delete,
	Up:     Up,
	Down:   Down,
	Right:  Right,
	Left:   Left,
	Enter:  Enter,
	Help:   Help,
	Quit:   Quit,
	Back:   Back,
	Copy:   Copy,
	Save:   Save,
}

var New = key.NewBinding(
	key.WithKeys("n"),
	key.WithHelp("n", "new"),
)
var Edit = key.NewBinding(
	key.WithKeys("e"),
	key.WithHelp("e", "edit"),
)
var Delete = key.NewBinding(
	key.WithKeys("d"),
	key.WithHelp("d", "delete"),
)
var Up = key.NewBinding(
	key.WithKeys("up"),
	key.WithHelp("↑", "move up"),
)
var Down = key.NewBinding(
	key.WithKeys("down"),
	key.WithHelp("↓", "move down"),
)
var Right = key.NewBinding(
	key.WithKeys("right"),
	key.WithHelp("→", "move right"),
)
var Left = key.NewBinding(
	key.WithKeys("left"),
	key.WithHelp("←", "move left"),
)
var Enter = key.NewBinding(
	key.WithKeys("enter"),
	key.WithHelp("enter", "enter"),
)
var Help = key.NewBinding(
	key.WithKeys("?"),
	key.WithHelp("?", "toggle help"),
)
var Quit = key.NewBinding(
	key.WithKeys("q"),
	key.WithHelp("q", "quit"),
)
var Back = key.NewBinding(
	key.WithKeys("esc"),
	key.WithHelp("esc", "back"),
)
var Copy = key.NewBinding(
	key.WithKeys("ctrl+c"),
	key.WithHelp("ctrl+c", "copy to clipboard"),
)
var Save = key.NewBinding(
	key.WithKeys("ctrl+s"),
	key.WithHelp("ctrl+s", "save form"),
)
