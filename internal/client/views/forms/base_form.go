package forms

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	kb "github.com/AndrXxX/goph-keeper/internal/client/views/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/views/styles"
)

type baseForm struct {
	help       help.Model
	title      string
	focusIndex int
	inputs     []textinput.Model
	fu         form.FieldsUpdater
	keys       *kb.KeyMap
	afterTitle []string
}

func NewBaseForm(title string, inputs []textinput.Model, fu form.FieldsUpdater) *baseForm {
	f := baseForm{
		help:   help.New(),
		title:  title,
		inputs: inputs,
		fu:     fu,
	}

	for i := range f.inputs {
		t := textinput.New()
		t.Cursor.Style = styles.Cursor
		t.CharLimit = 150
		f.inputs[i] = t
	}
	f.fu.Set(f.inputs, &f.focusIndex)
	return &f
}

func (f *baseForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *baseForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
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

func (f *baseForm) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(f.inputs))
	for i := range f.inputs {
		f.inputs[i], cmds[i] = f.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (f *baseForm) View() string {
	vList := []string{
		styles.Title.Margin(1).Render(f.title),
	}
	vList = append(vList, f.afterTitle...)
	for i := range f.inputs {
		vList = append(vList, f.inputs[i].View())
	}
	if f.keys != nil {
		vList = append(vList, styles.Help.Render(f.help.View(*f.keys)))
	}
	return lipgloss.JoinVertical(lipgloss.Left, vList...)
}
