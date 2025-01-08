package forms

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiagomelo/go-clipboard/clipboard"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/contract"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
)

const (
	pfLogin = iota
	pfPass
	pfDesc
)

var passwordFormKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Save, kb.Copy},
	Full: [][]key.Binding{
		{kb.Back, kb.Save, kb.Copy, kb.Quit},
		{kb.Up, kb.Down, kb.Enter},
	},
}

type passwordForm struct {
	item *entities.PasswordItem
	fu   form.FieldsUpdater
	*baseForm
	sm contract.SyncManager
}

func NewPasswordForm(item *entities.PasswordItem, sm contract.SyncManager) *passwordForm {
	m := passwordForm{
		baseForm: NewBaseForm("Create a new password", make([]textinput.Model, 3), form.FieldsUpdater{}),
		item:     item,
		sm:       sm,
	}
	m.baseForm.keys = &passwordFormKeys
	if m.item == nil {
		m.item = &entities.PasswordItem{}
	}

	m.baseForm.inputs[pfLogin].Prompt = "Login: "
	m.baseForm.inputs[pfLogin].SetValue(m.item.Login)

	m.baseForm.inputs[pfPass].Prompt = "Password: "
	m.baseForm.inputs[pfPass].SetValue(m.item.Password)
	//m.baseForm.inputs[pfPass].EchoMode = textinput.EchoPassword
	//m.baseForm.inputs[pfPass].EchoCharacter = '•'

	m.baseForm.inputs[pfDesc].Prompt = "Description: "
	m.baseForm.inputs[pfDesc].SetValue(m.item.Desc)
	return &m
}

func (f *passwordForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *passwordForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Back):
			return f, helpers.GenCmd(messages.ChangeView{Name: names.PasswordList})
		case key.Matches(msg, kb.Keys.Save):
			err := f.sm.Sync(datatypes.Passwords, []any{*f.getPasswordItem()})
			if err != nil {
				return f, helpers.GenCmd(messages.ShowError{Err: fmt.Sprintf("Ошибка при обновлении: %s", err)})
			}
			return f, tea.Batch(
				helpers.GenCmd(messages.ChangeView{Name: names.PasswordList}),
				helpers.GenCmd(messages.AddPassword{Item: f.getPasswordItem()}),
				helpers.GenCmd(messages.ShowMessage{Message: "Изменения сохранены"}),
			)
		case key.Matches(msg, kb.Keys.Copy):
			c := clipboard.New()
			err := c.CopyText(f.baseForm.inputs[f.baseForm.focusIndex].Value())
			if err != nil {
				return f, helpers.GenCmd(messages.ShowError{
					Err: fmt.Sprintf("failed to copy: %s", err.Error()),
				})
			}
			return f, helpers.GenCmd(messages.ShowMessage{Message: "value copied to clipboard"})
		}
	}
	_, cmd := f.baseForm.Update(msg)
	return f, cmd
}

func (f *passwordForm) getPasswordItem() *entities.PasswordItem {
	f.item.Login = f.baseForm.inputs[0].Value()
	f.item.Password = f.baseForm.inputs[1].Value()
	f.item.Desc = f.baseForm.inputs[2].Value()
	return f.item
}

func (f *passwordForm) View() string {
	return f.baseForm.View()
}
