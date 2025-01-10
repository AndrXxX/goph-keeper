package forms

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/locales"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
)

const (
	rfLogin = iota
	rfPassword
)

var registerFormKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Enter},
	Full: [][]key.Binding{
		{kb.Back, kb.Enter},
		{kb.Up, kb.Down},
	},
}

type registerForm struct {
	*baseForm
	r Registerer
	f *Factory
}

func newRegisterForm() *registerForm {
	m := registerForm{
		baseForm: NewBaseForm("Create a new account", make([]textinput.Model, 2), form.FieldsUpdater{}),
	}
	m.baseForm.keys = &registerFormKeys
	m.baseForm.inputs[rfLogin].Prompt = locales.FILogin
	m.baseForm.inputs[rfPassword].Prompt = locales.FIPassword
	return &m
}

func (f *registerForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *registerForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Back):
			return f, helpers.GenCmd(messages.ChangeView{
				Name: names.AuthMenu,
			})
		case key.Matches(msg, kb.Keys.Enter):
			u := &entities.User{Login: f.inputs[rfLogin].Value(), Password: f.inputs[rfPassword].Value()}
			token, err := f.r.Register(u)
			if err != nil {
				return f, helpers.GenCmd(messages.ShowError{Err: err.Error()})
			}
			u.Token = token
			cmdList := []tea.Cmd{
				helpers.GenCmd(messages.UpdateUser{User: u}),
				helpers.GenCmd(messages.ChangeView{Name: names.MasterPassRegForm, View: f.f.MasterPassRegForm()}),
				helpers.GenCmd(messages.ShowMessage{Message: "Successfully registered"}),
			}
			return f, tea.Batch(cmdList...)
		}
	}
	_, cmd := f.baseForm.Update(msg)
	return f, cmd
}

func (f *registerForm) View() string {
	return f.baseForm.View()
}
