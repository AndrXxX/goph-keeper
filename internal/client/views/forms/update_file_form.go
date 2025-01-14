package forms

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
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
	"github.com/AndrXxX/goph-keeper/internal/client/views/styles"
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
)

const (
	ffDesc = iota
)

var updateFileFormKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Save},
	Full: [][]key.Binding{
		{kb.Back, kb.Save, kb.Quit},
		{kb.Up, kb.Down, kb.Enter},
	},
}

type updateFileForm struct {
	item     *entities.FileItem
	creating bool
	*baseForm
}

func NewUpdateFileForm(item *entities.FileItem) *updateFileForm {
	m := updateFileForm{
		baseForm: NewBaseForm("File info", make([]textinput.Model, 1), form.FieldsUpdater{}),
		item:     item,
	}
	m.baseForm.keys = &updateFileFormKeys
	m.baseForm.inputs[ffDesc].Prompt = locales.FIDescription
	m.baseForm.inputs[ffDesc].SetValue(m.item.Desc)
	return &m
}

func (f *updateFileForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *updateFileForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Back):
			if f.item.IsStored() {
				return f, helpers.GenCmd(messages.ChangeView{Name: names.FileList})
			}
			return f, helpers.GenCmd(messages.ChangeView{Name: names.UploadFileForm})
		case key.Matches(msg, kb.Keys.Save):
			item := f.getFileItem()
			if _, err := govalidator.ValidateStruct(item); err != nil {
				return f, helpers.GenCmd(messages.ShowError{Err: fmt.Sprintf("Ошибка при обновлении: %s", err)})
			}
			return f, tea.Batch(
				helpers.GenCmd(messages.UploadItemUpdates{Type: datatypes.Files, Items: []any{*item}}),
				helpers.GenCmd(messages.ChangeView{Name: names.FileList}),
				helpers.GenCmd(messages.AddFile{Item: item}),
				helpers.GenCmd(messages.ShowMessage{Message: "Выполняется синхронизация изменений"}),
			)
		}
	}
	_, cmd := f.baseForm.Update(msg)
	return f, cmd
}

func (f *updateFileForm) getFileItem() *entities.FileItem {
	f.item.Desc = f.baseForm.inputs[ffDesc].Value()
	f.item.UpdatedAt = time.Now()
	return f.item
}

func (f *updateFileForm) View() string {
	f.baseForm.afterTitle = []string{styles.Focused.Margin(1).Render(f.item.Name)}
	return f.baseForm.View()
}
