package views

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiagomelo/go-clipboard/clipboard"

	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/pkg/entities"
)

var passwordFormKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Save, kb.Copy},
	Full: [][]key.Binding{
		{kb.Back, kb.Save, kb.Copy, kb.Quit},
		{kb.Up, kb.Down, kb.Enter},
	},
}

type passwordForm struct {
	item     *entities.PasswordItem
	creating bool
	fu       form.FieldsUpdater
	*baseForm
}

func NewPasswordForm(item *entities.PasswordItem) *passwordForm {
	m := passwordForm{
		baseForm: NewBaseForm("Create a new password", make([]textinput.Model, 3), form.FieldsUpdater{}),
		creating: item == nil,
		item:     item,
	}
	m.baseForm.keys = &passwordFormKeys
	if m.creating {
		m.item = &entities.PasswordItem{}
	}

	m.baseForm.inputs[0].Prompt = "Login: "
	m.baseForm.inputs[0].SetValue(m.item.Login)

	m.baseForm.inputs[1].Prompt = "Password: "
	m.baseForm.inputs[1].SetValue(m.item.Password)
	//m.baseForm.inputs[1].EchoMode = textinput.EchoPassword
	//m.baseForm.inputs[1].EchoCharacter = '•'

	m.baseForm.inputs[2].Prompt = "Description: "
	m.baseForm.inputs[2].SetValue(m.item.Desc)
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
			println("kb.Keys.Back")
			return f, func() tea.Msg {
				return messages.ChangeView{
					Name: names.PasswordList,
				}
			}
		case key.Matches(msg, kb.Keys.Save):
			// TODO: сделать уведомление
			var nMsg tea.Msg
			if f.creating {
				nMsg = messages.AddPassword{
					Item: f.getPasswordItem(),
				}
			}
			return f, func() tea.Msg {
				return messages.ChangeView{
					Name: names.PasswordList,
					Msg:  nMsg,
				}
			}
		case key.Matches(msg, kb.Keys.Copy):
			c := clipboard.New()
			err := c.CopyText(f.baseForm.inputs[f.baseForm.focusIndex].Value())
			if err != nil {
				println(err.Error())
			}
			// TODO: process error
			return f, nil
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
