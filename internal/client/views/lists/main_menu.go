package lists

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/tools/container/intsets"

	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	kb "github.com/AndrXxX/goph-keeper/internal/client/views/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/views/menuitems"
	"github.com/AndrXxX/goph-keeper/internal/client/views/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/internal/client/views/styles"
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
)

var mainMenuKeys = kb.KeyMap{
	Short: []key.Binding{kb.Quit, kb.Enter},
	Full: [][]key.Binding{
		{kb.Quit, kb.Enter},
	},
}

type mainMenu struct {
	list list.Model
	help help.Model
}

type mmOption func(a *mainMenu)

func withMenuItem(i menuitems.MainMenuItem) mmOption {
	return func(a *mainMenu) {
		a.list.InsertItem(intsets.MaxInt, i)
	}
}

func newMainMenu(opts ...mmOption) *mainMenu {
	m := &mainMenu{
		list: list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		help: help.New(),
	}
	m.list.SetShowStatusBar(false)
	m.list.SetShowHelp(false)
	m.list.Title = "Menu"
	m.list.Styles.Title = styles.Title
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *mainMenu) Init() tea.Cmd {
	return nil
}

func (m *mainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height/2)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Enter):
			return m, func() tea.Msg {
				selected := m.list.SelectedItem().(menuitems.MainMenuItem)
				switch selected.Code {
				case datatypes.Passwords:
					return messages.ChangeView{
						Name: names.PasswordList,
					}
				case datatypes.Notes:
					return messages.ChangeView{
						Name: names.NotesList,
					}
				case datatypes.BankCards:
					return messages.ChangeView{
						Name: names.BankCardList,
					}
				case datatypes.Files:
					return messages.ChangeView{
						Name: names.FileList,
					}
				}
				return nil
			}
		case key.Matches(msg, kb.Keys.Back):
			return m, helpers.GenCmd(messages.Quit{})
		}
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *mainMenu) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.list.View(), styles.Help.Render(m.help.View(mainMenuKeys)))
}
