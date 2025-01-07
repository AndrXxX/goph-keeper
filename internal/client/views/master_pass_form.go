package views

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
)

var masterPassFormKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Enter},
	Full: [][]key.Binding{
		{kb.Back, kb.Enter},
		{kb.Up, kb.Down},
	},
}

type masterPassForm struct {
	*baseForm
}

func NewMasterPassForm() *masterPassForm {
	m := masterPassForm{
		baseForm: NewBaseForm("Enter a master password to access", make([]textinput.Model, 1), form.FieldsUpdater{}),
	}
	m.baseForm.keys = &masterPassFormKeys
	m.baseForm.inputs[0].Prompt = "Password: "
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
			println("kb.Keys.Back")
			return f, func() tea.Msg {
				return messages.ChangeView{
					Name: names.AuthMenu,
				}
			}
		case key.Matches(msg, kb.Keys.Enter):
			// TODO: check login/pass
			return f, func() tea.Msg {
				return messages.ChangeView{
					Name: names.MainMenu,
				}
			}
		}
	}
	_, cmd := f.baseForm.Update(msg)
	return f, cmd
}

func (f *masterPassForm) View() string {
	return f.baseForm.View()
}
