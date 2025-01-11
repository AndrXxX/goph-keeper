package forms

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/locales"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
)

var masterPassRegFormKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Enter},
	Full: [][]key.Binding{
		{kb.Back, kb.Enter},
		{kb.Up, kb.Down},
	},
}

const (
	mprFormPassword = iota
	mprFormRepeat
)
const minPassLength = 5

type masterPassRegForm struct {
	*baseForm
}

func newMasterPassRegForm() *masterPassRegForm {
	m := masterPassRegForm{
		baseForm: NewBaseForm("Enter a master password to access", make([]textinput.Model, 2), form.FieldsUpdater{}),
	}
	m.baseForm.keys = &masterPassRegFormKeys
	m.baseForm.inputs[mprFormPassword].Prompt = locales.FIPassword
	m.baseForm.inputs[mprFormRepeat].Prompt = locales.FIRepeatPassword
	return &m
}

func (f *masterPassRegForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *masterPassRegForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Back):
			return f, helpers.GenCmd(messages.ChangeView{
				Name: names.AuthMenu,
			})
		case key.Matches(msg, kb.Keys.Enter):
			if len(f.baseForm.inputs[mprFormPassword].Value()) < minPassLength {
				return f, helpers.GenCmd(messages.ShowError{
					Err: fmt.Sprintf("password must be at least %d characters long", minPassLength),
				})
			}
			if f.baseForm.inputs[mprFormPassword].Value() != f.baseForm.inputs[mprFormRepeat].Value() {
				return f, helpers.GenCmd(messages.ShowError{Err: "passwords must be equal"})
			}
			return f, tea.Batch(
				helpers.GenCmd(messages.Auth{MasterPass: f.baseForm.inputs[mprFormPassword].Value()}),
			)
		}
	}
	_, cmd := f.baseForm.Update(msg)
	return f, cmd
}

func (f *masterPassRegForm) View() string {
	return f.baseForm.View()
}
