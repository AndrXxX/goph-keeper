package views

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
)

var margin = 5

type Container struct {
	help     help.Model
	loaded   bool
	current  names.ViewName
	views    Map
	quitting bool
}

func NewContainer(cols Map) *Container {
	return &Container{help: help.New(), views: cols}
}

func (m *Container) Init() tea.Cmd {
	m.current = names.AuthMenu
	return nil
}

func (m *Container) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmdList []tea.Cmd
		m.help.Width = msg.Width - margin
		for i := range m.views {
			_, cmd = m.views[i].Update(msg)
			cmdList = append(cmdList, cmd)
		}
		m.loaded = true
		return m, tea.Batch(cmdList...)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	case messages.ChangeView:
		m.current = msg.Name
		if msg.View != nil {
			m.views[m.current] = msg.View
		}
		if msg.Msg != nil {
			_, cmd = m.views[m.current].Update(msg.Msg)
		}
	}
	_, cmd = m.views[m.current].Update(msg)
	return m, cmd
}

func (m *Container) View() string {
	if m.quitting {
		return m.getStyle().Render("")
	}
	if !m.loaded {
		return m.getStyle().Render("loading...")
	}
	board := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.views[m.current].View(),
	)
	return lipgloss.JoinVertical(lipgloss.Left, m.getStyle().Render(board))
}

func (m *Container) getStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))
}
