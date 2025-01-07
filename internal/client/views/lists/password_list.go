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
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/internal/client/views/styles"
	"github.com/AndrXxX/goph-keeper/pkg/entities"
)

var passwordListKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Edit, kb.Delete, kb.New},
	Full: [][]key.Binding{
		{kb.Edit, kb.Delete, kb.New, kb.Quit},
		{kb.Up, kb.Down, kb.Enter, kb.Back},
	},
}

type passwordList struct {
	list list.Model
	help help.Model
}

func newPasswordList() *passwordList {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	defaultList.Title = "Passwords"
	defaultList.Styles.Title = styles.Title
	return &passwordList{list: defaultList, help: help.New()}
}

func (pl *passwordList) Init() tea.Cmd {
	return nil
}

func (pl *passwordList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		pl.list.SetSize(msg.Width/styles.InnerMargin, msg.Height/2)
	case messages.AddPassword:
		// TODO: action
		pl.list.InsertItem(-1, msg.Item)
		pl.View()
		return pl, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Edit, kb.Keys.Enter):
			if len(pl.list.VisibleItems()) != 0 {
				e := pl.list.SelectedItem().(*entities.PasswordItem)
				f := forms.NewPasswordForm(e)
				return f, helpers.GenCmd(messages.ChangeView{Name: names.PasswordForm, View: f})
			}
		case key.Matches(msg, kb.Keys.New):
			f := forms.NewPasswordForm(nil)
			return f, helpers.GenCmd(messages.ChangeView{Name: names.PasswordForm, View: f})
		case key.Matches(msg, kb.Keys.Back):
			return pl, helpers.GenCmd(messages.ChangeView{Name: names.MainMenu})
		case key.Matches(msg, kb.Keys.Delete):
			// TODO: approve + action
			return pl, pl.DeleteCurrent()
		}
	}
	pl.list, cmd = pl.list.Update(msg)
	return pl, cmd
}

func (pl *passwordList) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, pl.list.View(), pl.help.View(passwordListKeys))
}

func (pl *passwordList) DeleteCurrent() tea.Cmd {
	if len(pl.list.VisibleItems()) > 0 {
		pl.list.RemoveItem(pl.list.Index())
	}

	var cmd tea.Cmd
	pl.list, cmd = pl.list.Update(nil)
	return cmd
}
