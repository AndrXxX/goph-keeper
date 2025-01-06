package views

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tiagomelo/go-clipboard/clipboard"

	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/internal/client/views/styles"
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
	focusIndex int
	help       help.Model
	inputs     []textinput.Model
	cursorMode cursor.Mode
	item       *entities.PasswordItem
	creating   bool
	fu         form.FieldsUpdater
}

func NewPasswordForm(item *entities.PasswordItem) *passwordForm {
	m := passwordForm{
		inputs:   make([]textinput.Model, 3),
		item:     item,
		creating: item == nil,
		fu:       form.FieldsUpdater{},
	}
	if m.creating {
		m.item = &entities.PasswordItem{}
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = styles.Cursor
		t.CharLimit = 150

		switch i {
		case 0:
			t.Prompt = "Login: "
			t.SetValue(m.item.Login)
		case 1:
			t.Prompt = "Password: "
			t.SetValue(m.item.Password)
			//t.EchoMode = textinput.EchoPassword
			//t.EchoCharacter = '•'
		case 2:
			t.Prompt = "Description: "
			t.SetValue(m.item.Desc)
		}

		m.inputs[i] = t
	}
	m.fu.Set(m.inputs, &m.focusIndex)

	return &m
}

func (f *passwordForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *passwordForm) kbKeys() kb.KeyMap {
	return passwordFormKeys
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
			err := c.CopyText(f.inputs[f.focusIndex].Value())
			if err != nil {
				println(err.Error())
			}
			// TODO: process error
			return f, nil
		case key.Matches(msg, kb.Keys.Up):
			f.focusIndex--
			if f.focusIndex < 0 {
				f.focusIndex = len(f.inputs) - 1
			}
			return f, tea.Batch(f.fu.Set(f.inputs, &f.focusIndex)...)
		case key.Matches(msg, kb.Keys.Down, kb.Keys.Enter):
			f.focusIndex++
			if f.focusIndex >= len(f.inputs) {
				f.focusIndex = 0
			}
			return f, tea.Batch(f.fu.Set(f.inputs, &f.focusIndex)...)
		}
	}
	cmd := f.updateInputs(msg)
	return f, cmd
}

func (f *passwordForm) getPasswordItem() *entities.PasswordItem {
	f.item.Login = f.inputs[0].Value()
	f.item.Password = f.inputs[1].Value()
	f.item.Desc = f.inputs[2].Value()
	return f.item
}

func (f *passwordForm) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(f.inputs))
	for i := range f.inputs {
		f.inputs[i], cmds[i] = f.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (f *passwordForm) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		styles.Title.Margin(1).Render("Create a new password"),
		f.inputs[0].View(),
		f.inputs[1].View(),
		f.inputs[2].View(),
	)
}
