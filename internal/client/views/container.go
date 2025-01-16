package views

import (
	"fmt"
	"sync/atomic"
	"time"

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
	help           help.Model
	loaded         bool
	current        names.ViewName
	views          Map
	quitting       atomic.Bool
	errors         helpers.MsgList
	messages       helpers.MsgList
	uo             map[tea.Msg]UpdateOption
	bi             *contract.BuildInfo
	updateInterval time.Duration
}

func (m *container) Init() tea.Cmd {
	return m.Tick()
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
	case messages.Tick:
		cmdList = append(cmdList, m.Tick())
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
		if f, ok := m.uo[helpers.GenMsgKey(msg)]; ok {
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
	if err := m.errors.Join("\n"); err != "" {
		items = append(items, styles.Error.Render(err))
	}
	if mes := m.messages.Join("\n"); mes != "" {
		items = append(items, styles.Info.Render(mes))
	}
	if m.bi != nil {
		ver := fmt.Sprintf("ver. %s [%s]", m.bi.Version, m.bi.Date)
		items = []string{lipgloss.JoinVertical(lipgloss.Left, items...)}
		items = append(items, styles.Blurred.Margin(1, 0, 0).Render(ver))
	}
	return styles.Border.Render(lipgloss.JoinVertical(lipgloss.Center, items...))
}

func (m *container) Tick() tea.Cmd {
	return tea.Tick(m.updateInterval, func(t time.Time) tea.Msg {
		return messages.Tick(t)
	})
}
