package forms

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiagomelo/go-clipboard/clipboard"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/locales"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	kb "github.com/AndrXxX/goph-keeper/internal/client/views/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/views/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
)

const (
	bcNumber = iota
	bcCVC
	bcValidity
	bcHolder
	bcDesc
)

var bankCardFormKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Save, kb.Copy},
	Full: [][]key.Binding{
		{kb.Back, kb.Save, kb.Copy, kb.Quit},
		{kb.Up, kb.Down, kb.Enter},
	},
}

type bankCardForm struct {
	item     *entities.BankCardItem
	creating bool
	*baseForm
}

func NewBankCardForm(item *entities.BankCardItem) *bankCardForm {
	m := bankCardForm{
		baseForm: NewBaseForm("Create/edit bank card", make([]textinput.Model, 5), form.FieldsUpdater{}),
		creating: item == nil,
		item:     item,
	}
	m.baseForm.keys = &bankCardFormKeys
	if m.creating {
		m.item = &entities.BankCardItem{}
	}

	m.baseForm.inputs[bcNumber].Prompt = locales.FINumber
	m.baseForm.inputs[bcNumber].SetValue(m.item.Number)

	m.baseForm.inputs[bcCVC].Prompt = locales.FICVCCode
	m.baseForm.inputs[bcCVC].SetValue(m.item.CVCCode)

	m.baseForm.inputs[bcValidity].Prompt = locales.FIValidity
	m.baseForm.inputs[bcValidity].SetValue(m.item.Validity)

	m.baseForm.inputs[bcHolder].Prompt = locales.FICardholder
	m.baseForm.inputs[bcHolder].SetValue(m.item.Cardholder)

	m.baseForm.inputs[bcDesc].Prompt = locales.FIDescription
	m.baseForm.inputs[bcDesc].SetValue(m.item.Desc)

	return &m
}

func (f *bankCardForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *bankCardForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Back):
			return f, helpers.GenCmd(messages.ChangeView{Name: names.BankCardList})
		case key.Matches(msg, kb.Keys.Save):
			item := f.getBankCardItem()
			if _, err := govalidator.ValidateStruct(item); err != nil {
				return f, helpers.GenCmd(messages.ValidityError{Error: err})
			}
			return f, tea.Batch(
				helpers.GenCmd(messages.UploadItemUpdates{Type: datatypes.BankCards, Items: []any{*item}}),
				helpers.GenCmd(messages.ChangeView{Name: names.BankCardList}),
				helpers.GenCmd(messages.AddBankCard{Item: item}),
				helpers.GenCmd(messages.ShowMessage{Message: "Выполняется синхронизация изменений"}),
			)
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

func (f *bankCardForm) getBankCardItem() *entities.BankCardItem {
	f.item.Number = f.baseForm.inputs[bcNumber].Value()
	f.item.CVCCode = f.baseForm.inputs[bcCVC].Value()
	f.item.Validity = f.baseForm.inputs[bcValidity].Value()
	f.item.Cardholder = f.baseForm.inputs[bcHolder].Value()
	f.item.Desc = f.baseForm.inputs[bcDesc].Value()
	now := time.Now()
	f.item.UpdatedAt = now
	return f.item
}

func (f *bankCardForm) View() string {
	return f.baseForm.View()
}
