package views

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiagomelo/go-clipboard/clipboard"

	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/pkg/entities"
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
		baseForm: NewBaseForm("Create a new note", make([]textinput.Model, 1), form.FieldsUpdater{}),
		creating: item == nil,
		item:     item,
	}
	m.baseForm.keys = &noteFormKeys
	if m.creating {
		m.item = &entities.NoteItem{}
	}

	m.baseForm.inputs[0].Prompt = "Text: "
	m.baseForm.inputs[0].SetValue(m.item.Text)

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
			// TODO: сделать уведомление
			var nMsg tea.Msg
			if f.creating {
				nMsg = messages.AddNote{
					Item: f.getNoteItem(),
				}
			}
			return f, func() tea.Msg {
				return messages.ChangeView{
					Name: names.NotesList,
					Msg:  nMsg,
				}
			}
		case key.Matches(msg, kb.Keys.Copy):
			c := clipboard.New()
			err := c.CopyText(f.baseForm.inputs[f.baseForm.focusIndex].Value())
			if err != nil {
				println(err.Error())
			}
			// TODO: process error
			return f, nil
		}
	}
	_, cmd := f.baseForm.Update(msg)
	return f, cmd
}

func (f *noteForm) getNoteItem() *entities.NoteItem {
	f.item.Text = f.baseForm.inputs[0].Value()
	return f.item
}

func (f *noteForm) View() string {
	return f.baseForm.View()
}
