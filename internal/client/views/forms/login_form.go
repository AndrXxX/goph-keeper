package forms

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/locales"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/state"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
)

const (
	lfLogin = iota
	lfPassword
)

var loginFormKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Enter},
	Full: [][]key.Binding{
		{kb.Back, kb.Enter},
		{kb.Up, kb.Down},
	},
}

type loginForm struct {
	*baseForm
	l Loginer
	s *state.AppState
	f *Factory
}

func newLoginForm() *loginForm {
	m := loginForm{
		baseForm: NewBaseForm("Enter an exist account", make([]textinput.Model, 2), form.FieldsUpdater{}),
	}
	m.baseForm.keys = &loginFormKeys
	m.baseForm.inputs[lfLogin].Prompt = locales.FILogin
	m.baseForm.inputs[lfPassword].Prompt = locales.FIPassword
	return &m
}

func (f *loginForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *loginForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Back):
			return f, helpers.GenCmd(messages.ChangeView{Name: names.AuthMenu})
		case key.Matches(msg, kb.Keys.Enter):
			f.s.User.Login = f.inputs[lfLogin].Value()
			f.s.User.Password = f.inputs[lfPassword].Value()
			token, err := f.l.Login(f.s.User)
			if err != nil {
				return f, helpers.GenCmd(messages.ShowError{Err: err.Error()})
			}
			f.s.User.Token = token
			cmdList := []tea.Cmd{
				helpers.GenCmd(messages.ChangeView{Name: names.MasterPassRegForm, View: f.f.MasterPassRegForm()}),
				helpers.GenCmd(messages.ShowMessage{Message: "Successfully logged in"}),
			}
			return f, tea.Batch(cmdList...)
		}
	}
	_, cmd := f.baseForm.Update(msg)
	return f, cmd
}

func (f *loginForm) View() string {
	return f.baseForm.View()
}
