package views

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/internal/client/views/styles"
)

const errorsTimeout = 2 * time.Second

type Container struct {
	help     help.Model
	loaded   bool
	current  names.ViewName
	views    Map
	quitting bool
	errors   sync.Map
}

func NewContainer(cols Map) *Container {
	return &Container{help: help.New(), views: cols, errors: sync.Map{}}
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
		m.help.Width = msg.Width - styles.InnerMargin
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
	case messages.ShowError:
		m.errors.Store(msg.Err, msg.Err)
		go func() {
			time.Sleep(errorsTimeout)
			m.errors.Clear()
		}()
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
	eb := strings.Builder{}

	m.errors.Range(func(_, v any) bool {
		eb.WriteString(fmt.Sprintf("%s\n", v.(string)))
		return true
	})
	err := eb.String()
	if err != "" {
		err = styles.Error.Render(err)
	}
	return m.getStyle().Render(lipgloss.JoinVertical(lipgloss.Left, board, err))
}

func (m *Container) getStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))
}
