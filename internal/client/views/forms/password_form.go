package forms

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiagomelo/go-clipboard/clipboard"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/locales"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/messages"
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
	*baseForm
}

func NewPasswordForm(item *entities.PasswordItem) *passwordForm {
	m := passwordForm{
		baseForm: NewBaseForm("Create a new password", make([]textinput.Model, 3), form.FieldsUpdater{}),
		item:     item,
	}
	m.baseForm.keys = &passwordFormKeys
	if m.item == nil {
		m.item = &entities.PasswordItem{}
	}

	m.baseForm.inputs[pfLogin].Prompt = locales.FILogin
	m.baseForm.inputs[pfLogin].SetValue(m.item.Login)

	m.baseForm.inputs[pfPass].Prompt = locales.FIPassword
	m.baseForm.inputs[pfPass].SetValue(m.item.Password)
	//m.baseForm.inputs[pfPass].EchoMode = textinput.EchoPassword
	//m.baseForm.inputs[pfPass].EchoCharacter = '•'

	m.baseForm.inputs[pfDesc].Prompt = locales.FIDescription
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
			item := f.getPasswordItem()
			if _, err := govalidator.ValidateStruct(item); err != nil {
				return f, helpers.GenCmd(messages.ValidityError{Error: err})
			}
			return f, tea.Batch(
				helpers.GenCmd(messages.UploadItemUpdates{Type: datatypes.Passwords, Items: []any{*item}}),
				helpers.GenCmd(messages.ChangeView{Name: names.PasswordList}),
				helpers.GenCmd(messages.AddPassword{Item: item}),
				helpers.GenCmd(messages.ShowMessage{Message: "Выполняется синхронизация изменений"}),
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
	f.item.Login = f.baseForm.inputs[pfLogin].Value()
	f.item.Password = f.baseForm.inputs[pfPass].Value()
	f.item.Desc = f.baseForm.inputs[pfDesc].Value()
	f.item.UpdatedAt = time.Now()
	return f.item
}

func (f *passwordForm) View() string {
	return f.baseForm.View()
}
