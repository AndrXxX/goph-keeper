package lists

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/contract"
	"github.com/AndrXxX/goph-keeper/internal/client/views/forms"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/internal/client/views/styles"
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
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
	sm   contract.SyncManager
	lr   refresher
}

func newPasswordList() *passwordList {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	defaultList.Title = "Passwords"
	defaultList.Styles.Title = styles.Title
	return &passwordList{list: defaultList, help: help.New()}
}

func (l *passwordList) Init() tea.Cmd {
	l.lr.Refresh()
	return nil
}

func (l *passwordList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(l.list.Items()) == 0 {
		l.lr.Refresh()
	}
	l.lr.RefreshIn(refreshListInterval)
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.list.SetSize(msg.Width/styles.InnerMargin, msg.Height/2)
	case messages.AddPassword:
		err := l.sm.Sync(datatypes.Passwords, []any{*msg.Item})
		if err != nil {
			return l, helpers.GenCmd(messages.ShowError{Err: "Ошибка при обновлении"})
		}
		return l, helpers.GenCmd(messages.ShowMessage{Message: "Изменения сохранены"})
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Edit, kb.Keys.Enter):
			if len(l.list.VisibleItems()) != 0 {
				e := l.list.SelectedItem().(*entities.PasswordItem)
				f := forms.NewPasswordForm(e)
				return f, helpers.GenCmd(messages.ChangeView{Name: names.PasswordForm, View: f})
			}
		case key.Matches(msg, kb.Keys.New):
			f := forms.NewPasswordForm(nil)
			return f, helpers.GenCmd(messages.ChangeView{Name: names.PasswordForm, View: f})
		case key.Matches(msg, kb.Keys.Back):
			return l, helpers.GenCmd(messages.ChangeView{Name: names.MainMenu})
		case key.Matches(msg, kb.Keys.Delete):
			// TODO: approve + action
			return l, l.DeleteCurrent()
		}
	}
	l.list, cmd = l.list.Update(msg)
	return l, cmd
}

func (l *passwordList) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, l.list.View(), l.help.View(passwordListKeys))
}

func (l *passwordList) DeleteCurrent() tea.Cmd {
	if len(l.list.VisibleItems()) > 0 {
		l.list.RemoveItem(l.list.Index())
	}

	var cmd tea.Cmd
	l.list, cmd = l.list.Update(nil)
	return cmd
}
