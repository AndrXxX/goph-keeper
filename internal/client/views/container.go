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
	messages sync.Map
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
			m.errors.Delete(msg.Err)
		}()
	case messages.ShowMessage:
		m.messages.Store(msg.Message, msg.Message)
		go func() {
			time.Sleep(errorsTimeout)
			m.messages.Delete(msg.Message)
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
	mb := strings.Builder{}

	m.errors.Range(func(_, v any) bool {
		eb.WriteString(fmt.Sprintf("%s\n", v.(string)))
		return true
	})
	m.messages.Range(func(_, v any) bool {
		mb.WriteString(fmt.Sprintf("%s\n", v.(string)))
		return true
	})
	err := m.collectMessages(&m.errors)
	if err != "" {
		err = styles.Error.Render(err)
	}
	mes := m.collectMessages(&m.messages)
	if mes != "" {
		mes = styles.Info.Render(mes)
	}
	return m.getStyle().Render(lipgloss.JoinVertical(lipgloss.Left, board, err, mes))
}

func (m *Container) collectMessages(l *sync.Map) string {
	b := strings.Builder{}
	l.Range(func(_, v any) bool {
		b.WriteString(fmt.Sprintf("%s\n", v.(string)))
		return true
	})
	return b.String()
}

func (m *Container) getStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))
}
