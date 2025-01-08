package forms

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/state"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
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
	s *state.AppState
}

func newMasterPassAuthForm() *masterPassAuthForm {
	m := masterPassAuthForm{
		baseForm: NewBaseForm("Enter master password to access", make([]textinput.Model, 1), form.FieldsUpdater{}),
	}
	m.baseForm.keys = &masterPassAuthFormKeys
	m.baseForm.inputs[mpaFormPassword].Prompt = "Password: "
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
			f.s.User = &entities.User{}
			f.s.User.MasterPassword = f.baseForm.inputs[mprFormPassword].Value()
			err := f.s.Auth()
			if err != nil {
				return f, tea.Batch(helpers.GenCmd(messages.ShowError{Err: err.Error()}))
			}
			return f, tea.Batch(helpers.GenCmd(messages.ChangeView{Name: names.MainMenu}))
		}
	}
	_, cmd := f.baseForm.Update(msg)
	return f, cmd
}

func (f *masterPassAuthForm) View() string {
	return f.baseForm.View()
}
