package views

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/views/contract"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/internal/client/views/styles"
)

type container struct {
	help     help.Model
	loaded   bool
	current  names.ViewName
	views    Map
	quitting atomic.Bool
	errors   sync.Map
	messages sync.Map
	uo       map[tea.Msg]UpdateOption
	bi       *contract.BuildInfo
}

func (m *container) Init() tea.Cmd {
	return nil
}

func (m *container) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmdList := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmdList []tea.Cmd
		m.help.Width = msg.Width - styles.InnerMargin
		for i := range m.views {
			_, cmd := m.views[i].Update(msg)
			cmdList = append(cmdList, cmd)
		}
		m.loaded = true
		return m, tea.Batch(cmdList...)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Quit):
			return m, helpers.GenCmd(messages.Quit{})
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
		//cmdList = append(cmdList, helpers.GenCmd(tea.WindowSize()))
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
	if m.quitting.Load() {
		return ""
	}
	if !m.loaded {
		return styles.Border.Render("loading...")
	}
	var items []string
	items = append(items, lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.views[m.current].View(),
	))
	if err := m.collectMessages(&m.errors); err != "" {
		items = append(items, styles.Error.Render(err))
	}
	if mes := m.collectMessages(&m.messages); mes != "" {
		items = append(items, styles.Info.Render(mes))
	}
	var bottom string
	if m.bi != nil {
		bottom = lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.Blurred.Render(fmt.Sprintf("ver. %s [%s]", m.bi.Version, m.bi.Date)),
		)
	}
	return styles.Border.Render(lipgloss.JoinVertical(lipgloss.Left, items...), bottom)
}

func (m *container) collectMessages(l *sync.Map) string {
	b := strings.Builder{}
	l.Range(func(_, v any) bool {
		b.WriteString(fmt.Sprintf("%s\n", v.(string)))
		return true
	})
	return b.String()
}
