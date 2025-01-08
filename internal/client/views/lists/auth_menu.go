package lists

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/forms"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/menuitems"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/internal/client/views/styles"
)

var authMenuKeys = kb.KeyMap{
	Short: []key.Binding{kb.Quit, kb.Enter},
	Full: [][]key.Binding{
		{kb.Quit, kb.Enter},
	},
}

type authMenu struct {
	list list.Model
	help help.Model
	f    *forms.Factory
}

func newAuthMenu() *authMenu {
	defaultList := list.New([]list.Item{
		menuitems.AuthItem{Name: "Register", Code: "register", Desc: "Create a new account"},
		menuitems.AuthItem{Name: "Login", Code: "login", Desc: "Enter an exist account"},
		menuitems.AuthItem{Name: "Enter", Code: "master_pass", Desc: "Enter a master password to access"},
	}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	defaultList.Title = "Goph Keeper"
	defaultList.Styles.Title = styles.Title
	return &authMenu{list: defaultList, help: help.New()}
}

func (m *authMenu) Init() tea.Cmd {
	return nil
}

func (m *authMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width/styles.InnerMargin, msg.Height/2)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Enter):
			selected := m.list.SelectedItem().(menuitems.AuthItem)
			switch selected.Code {
			case "register":
				return m, helpers.GenCmd(messages.ChangeView{Name: names.RegisterForm, View: m.f.RegisterForm()})
			case "login":
				return m, helpers.GenCmd(messages.ChangeView{Name: names.LoginForm, View: m.f.LoginForm()})
			case "master_pass":
				return m, helpers.GenCmd(messages.ChangeView{Name: names.MasterPassAuthForm, View: m.f.MasterPassAuthForm()})
			}
			return m, nil
		}
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *authMenu) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.list.View(), m.help.View(authMenuKeys))
}
