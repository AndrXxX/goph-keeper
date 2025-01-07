package forms

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/pkg/entities"
)

var fileFormKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Save},
	Full: [][]key.Binding{
		{kb.Back, kb.Save, kb.Quit},
		{kb.Up, kb.Down, kb.Enter},
	},
}

type fileForm struct {
	item     *entities.FileItem
	creating bool
	fu       form.FieldsUpdater
	*baseForm
}

func NewFileForm(item *entities.FileItem) *fileForm {
	m := fileForm{
		baseForm: NewBaseForm("Create/update file", make([]textinput.Model, 2), form.FieldsUpdater{}),
		creating: item == nil,
		item:     item,
	}
	m.baseForm.keys = &fileFormKeys
	if m.creating {
		m.item = &entities.FileItem{}
	}

	m.baseForm.inputs[0].Prompt = "File: "
	m.baseForm.inputs[0].SetValue(m.item.Data)

	m.baseForm.inputs[1].Prompt = "Description: "
	m.baseForm.inputs[1].SetValue(m.item.Desc)

	return &m
}

func (f *fileForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *fileForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Back):
			return f, func() tea.Msg {
				return messages.ChangeView{
					Name: names.FileList,
				}
			}
		case key.Matches(msg, kb.Keys.Save):
			// TODO: сделать уведомление
			var nMsg tea.Msg
			if f.creating {
				nMsg = messages.AddFile{
					Item: f.getFileItem(),
				}
			}
			return f, func() tea.Msg {
				return messages.ChangeView{
					Name: names.FileList,
					Msg:  nMsg,
				}
			}
		}
	}
	_, cmd := f.baseForm.Update(msg)
	return f, cmd
}

func (f *fileForm) getFileItem() *entities.FileItem {
	f.item.Data = f.baseForm.inputs[0].Value()
	f.item.Desc = f.baseForm.inputs[1].Value()
	return f.item
}

func (f *fileForm) View() string {
	return f.baseForm.View()
}
