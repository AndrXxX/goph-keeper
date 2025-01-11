package forms

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/locales"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
)

var masterPassAuthFormKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Enter},
	Full: [][]key.Binding{
		{kb.Back, kb.Enter},
		{kb.Up, kb.Down},
	},
}

const (
	mpaFormPassword = iota
)

type masterPassAuthForm struct {
	*baseForm
}

func newMasterPassAuthForm() *masterPassAuthForm {
	m := masterPassAuthForm{
		baseForm: NewBaseForm("Enter master password to access", make([]textinput.Model, 1), form.FieldsUpdater{}),
	}
	m.baseForm.keys = &masterPassAuthFormKeys
	m.baseForm.inputs[mpaFormPassword].Prompt = locales.FIPassword
	return &m
}

func (f *masterPassAuthForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *masterPassAuthForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Back):
			return f, helpers.GenCmd(messages.ChangeView{
				Name: names.AuthMenu,
			})
		case key.Matches(msg, kb.Keys.Enter):
			return f, tea.Batch(
				helpers.GenCmd(messages.UpdateUser{User: &entities.User{}}),
				helpers.GenCmd(messages.Auth{MasterPass: f.baseForm.inputs[mpaFormPassword].Value()}),
			)
		}
	}
	_, cmd := f.baseForm.Update(msg)
	return f, cmd
}

func (f *masterPassAuthForm) View() string {
	return f.baseForm.View()
}
