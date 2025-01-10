package views

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/AndrXxX/goph-keeper/internal/client/jobs"
	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/contract"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/internal/client/views/styles"
)

type container struct {
	help     help.Model
	loaded   bool
	current  names.ViewName
	views    Map
	quitting bool
	errors   sync.Map
	messages sync.Map
	uo       map[tea.Msg]UpdateOption
	sm       contract.SyncManager
	qr       contract.QueueRunner
}

func (m *container) Init() tea.Cmd {
	m.current = names.AuthMenu
	return nil
}

func (m *container) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmdList := make([]tea.Cmd, 0)
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
	case messages.UploadItemUpdates:
		err := m.qr.AddJob(&jobs.UploadItemsUpdatesJob{
			Type:        msg.Type,
			Items:       msg.Items,
			SyncManager: m.sm,
		})
		if err != nil {
			return m, helpers.GenCmd(messages.ShowError{Err: fmt.Sprintf("Ошибка при обновлении: %s", err)})
		}
	case messages.ChangeView:
		m.current = msg.Name
		if msg.View != nil {
			m.views[m.current] = msg.View
		}
		if msg.Msg != nil {
			_, cmd := m.views[m.current].Update(msg.Msg)
			cmdList = append(cmdList, cmd)
		}
	default:
		mgsK := fmt.Sprintf("%T", msg)
		if f, ok := m.uo[mgsK]; ok {
			_, cmd := f(msg)
			cmdList = append(cmdList, cmd)
		}
	}
	_, cmd := m.views[m.current].Update(msg)
	cmdList = append(cmdList, cmd)
	return m, tea.Batch(cmdList...)
}

func (m *container) View() string {
	if m.quitting {
		return styles.Border.Render("")
	}
	if !m.loaded {
		return styles.Border.Render("loading...")
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
	return styles.Border.Render(lipgloss.JoinVertical(lipgloss.Left, board, err, mes))
}

func (m *container) collectMessages(l *sync.Map) string {
	b := strings.Builder{}
	l.Range(func(_, v any) bool {
		b.WriteString(fmt.Sprintf("%s\n", v.(string)))
		return true
	})
	return b.String()
}
