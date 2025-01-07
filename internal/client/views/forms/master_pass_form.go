package forms

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/state"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
)

var masterPassFormKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Enter},
	Full: [][]key.Binding{
		{kb.Back, kb.Enter},
		{kb.Up, kb.Down},
	},
}

const (
	mpFormPassword = iota
	mpFormRepeat
)
const minPassLength = 5

type masterPassForm struct {
	*baseForm
	s *state.AppState
}

func newMasterPassForm() *masterPassForm {
	m := masterPassForm{
		baseForm: NewBaseForm("Enter a master password to access", make([]textinput.Model, 2), form.FieldsUpdater{}),
	}
	m.baseForm.keys = &masterPassFormKeys
	m.baseForm.inputs[mpFormPassword].Prompt = "Password: "
	m.baseForm.inputs[mpFormRepeat].Prompt = "Repeat password: "
	return &m
}

func (f *masterPassForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *masterPassForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Back):
			return f, helpers.GenCmd(messages.ChangeView{
				Name: names.AuthMenu,
			})
		case key.Matches(msg, kb.Keys.Enter):
			if len(f.baseForm.inputs[mpFormPassword].Value()) < minPassLength {
				return f, helpers.GenCmd(messages.ShowError{
					Err: fmt.Sprintf("password must be at least %d characters long", minPassLength),
				})
			}
			if f.baseForm.inputs[mpFormPassword].Value() != f.baseForm.inputs[mpFormRepeat].Value() {
				return f, helpers.GenCmd(messages.ShowError{Err: "passwords must be equal"})
			}
			f.s.User.MasterPassword = f.baseForm.inputs[mpFormPassword].Value()
			// TODO: access to local DB
			cmdList := []tea.Cmd{
				helpers.GenCmd(messages.ChangeView{Name: names.MainMenu}),
				helpers.GenCmd(messages.ShowMessage{Message: "Successfully!"}), // TODO
			}
			return f, tea.Batch(cmdList...)
		}
	}
	_, cmd := f.baseForm.Update(msg)
	return f, cmd
}

func (f *masterPassForm) View() string {
	return f.baseForm.View()
}
