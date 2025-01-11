package forms

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/locales"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
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

	m.baseForm.inputs[0].Prompt = locales.FIFile
	m.baseForm.inputs[0].SetValue(m.item.Data)

	m.baseForm.inputs[1].Prompt = locales.FIDescription
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
			return f, helpers.GenCmd(messages.ChangeView{
				Name: names.FileList,
			})
		case key.Matches(msg, kb.Keys.Save):
			var nMsg tea.Msg
			if f.creating {
				nMsg = messages.AddFile{Item: f.getFileItem()}
			}
			cmdList := []tea.Cmd{
				helpers.GenCmd(messages.ChangeView{Name: names.FileList, Msg: nMsg}),
				helpers.GenCmd(messages.ShowMessage{Message: "file saved"}),
			}
			return f, tea.Batch(cmdList...)
		}
	}
	_, cmd := f.baseForm.Update(msg)
	return f, cmd
}

func (f *fileForm) getFileItem() *entities.FileItem {
	f.item.Data = f.baseForm.inputs[0].Value()
	f.item.Desc = f.baseForm.inputs[1].Value()
	f.item.UpdatedAt = time.Now()
	return f.item
}

func (f *fileForm) View() string {
	return f.baseForm.View()
}
