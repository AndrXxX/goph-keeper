package views

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/pkg/entities"
)

var authFormKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Enter},
	Full: [][]key.Binding{
		{kb.Back, kb.Enter},
		{kb.Up, kb.Down},
	},
}

type authForm struct {
	focusIndex int
	help       help.Model
	inputs     []textinput.Model
	item       *entities.PasswordItem
	creating   bool
	fu         form.FieldsUpdater
	*baseForm
}

func NewAuthForm() *authForm {
	m := authForm{
		baseForm: NewBaseForm("Create a new password", make([]textinput.Model, 2), form.FieldsUpdater{}),
	}

	m.baseForm.inputs[0].Prompt = "Login: "
	m.baseForm.inputs[1].Prompt = "Password: "

	return &m
}

func (f *authForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *authForm) kbKeys() kb.KeyMap {
	return authFormKeys
}

func (f *authForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (f *authForm) View() string {
	return f.baseForm.View()
}
