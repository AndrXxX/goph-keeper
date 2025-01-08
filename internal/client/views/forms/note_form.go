package forms

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiagomelo/go-clipboard/clipboard"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
)

var noteFormKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Save, kb.Copy},
	Full: [][]key.Binding{
		{kb.Back, kb.Save, kb.Copy, kb.Quit},
		{kb.Up, kb.Down, kb.Enter},
	},
}

type noteForm struct {
	item     *entities.NoteItem
	creating bool
	fu       form.FieldsUpdater
	*baseForm
}

func NewNoteForm(item *entities.NoteItem) *noteForm {
	m := noteForm{
		baseForm: NewBaseForm("Create a new note", make([]textinput.Model, 2), form.FieldsUpdater{}),
		creating: item == nil,
		item:     item,
	}
	m.baseForm.keys = &noteFormKeys
	if m.creating {
		m.item = &entities.NoteItem{}
	}

	m.baseForm.inputs[0].Prompt = "Text: "
	m.baseForm.inputs[0].SetValue(m.item.Text)

	m.baseForm.inputs[1].Prompt = "Description: "
	m.baseForm.inputs[1].SetValue(m.item.Desc)

	return &m
}

func (f *noteForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *noteForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Back):
			return f, func() tea.Msg {
				return messages.ChangeView{
					Name: names.NotesList,
				}
			}
		case key.Matches(msg, kb.Keys.Save):
			return f, func() tea.Msg {
				return messages.ChangeView{
					Name: names.NotesList,
					Msg: messages.AddNote{
						Item: f.getNoteItem(),
					},
				}
			}
		case key.Matches(msg, kb.Keys.Copy):
			c := clipboard.New()
			err := c.CopyText(f.baseForm.inputs[f.baseForm.focusIndex].Value())
			if err != nil {
				return f, helpers.GenCmd(messages.ShowError{Err: fmt.Sprintf("failed to copy: %s", err.Error())})
			}
			return f, helpers.GenCmd(messages.ShowMessage{Message: "value copied to clipboard"})
		}
	}
	_, cmd := f.baseForm.Update(msg)
	return f, cmd
}

func (f *noteForm) getNoteItem() *entities.NoteItem {
	f.item.Text = f.baseForm.inputs[0].Value()
	f.item.Desc = f.baseForm.inputs[1].Value()
	return f.item
}

func (f *noteForm) View() string {
	return f.baseForm.View()
}
