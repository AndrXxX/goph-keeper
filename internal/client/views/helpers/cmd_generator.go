package helpers

import tea "github.com/charmbracelet/bubbletea"

func GenCmd(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
