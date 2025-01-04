package messages

import tea "github.com/charmbracelet/bubbletea"

type ChangeView struct {
	Name string
	View tea.Model
	Msg  tea.Msg
}
