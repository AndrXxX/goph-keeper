package views

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/menuitems"
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
	list   list.Model
	help   help.Model
	height int
	width  int
}

func NewMainMenu() *mainMenu {
	defaultList := list.New([]list.Item{
		menuitems.MainMenuItem{Name: "Passwords", Code: datatypes.Passwords, Desc: "Manage passwords"},
		menuitems.MainMenuItem{Name: "Notes", Code: datatypes.Notes, Desc: "Manage notes"},
		menuitems.MainMenuItem{Name: "Bank Cards", Code: datatypes.BankCards, Desc: "Manage bank cards"},
		menuitems.MainMenuItem{Name: "Files", Code: datatypes.Files, Desc: "Manage files"},
	}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	defaultList.Title = "Menu"
	defaultList.Styles.Title = styles.Title
	return &mainMenu{list: defaultList, help: help.New()}
}

func (m *mainMenu) Init() tea.Cmd {
	return nil
}

func (m *mainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width/margin, msg.Height/2)
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
				}
				// TODO
				return nil
			}
		}
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *mainMenu) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.list.View(), m.help.View(mainMenuKeys))
}
