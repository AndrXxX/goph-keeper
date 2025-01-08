package forms

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiagomelo/go-clipboard/clipboard"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/contract"
	"github.com/AndrXxX/goph-keeper/internal/client/views/form"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
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
	fu       form.FieldsUpdater
	*baseForm
	sm contract.SyncManager
}

func NewBankCardForm(item *entities.BankCardItem, sm contract.SyncManager) *bankCardForm {
	m := bankCardForm{
		baseForm: NewBaseForm("Create/edit bank card", make([]textinput.Model, 5), form.FieldsUpdater{}),
		creating: item == nil,
		item:     item,
		sm:       sm,
	}
	m.baseForm.keys = &bankCardFormKeys
	if m.creating {
		m.item = &entities.BankCardItem{}
	}

	m.baseForm.inputs[bcNumber].Prompt = "Number: "
	m.baseForm.inputs[bcNumber].SetValue(m.item.Number)

	m.baseForm.inputs[bcCVC].Prompt = "CVCCode: "
	m.baseForm.inputs[bcCVC].SetValue(m.item.CVCCode)

	m.baseForm.inputs[bcValidity].Prompt = "Validity: "
	m.baseForm.inputs[bcValidity].SetValue(m.item.Validity)

	m.baseForm.inputs[bcHolder].Prompt = "Cardholder: "
	m.baseForm.inputs[bcHolder].SetValue(m.item.Cardholder)

	m.baseForm.inputs[bcDesc].Prompt = "Description: "
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
			err := f.sm.Sync(datatypes.BankCards, []any{*f.getBankCardItem()})
			if err != nil {
				return f, helpers.GenCmd(messages.ShowError{Err: fmt.Sprintf("Ошибка при обновлении: %s", err)})
			}
			return f, tea.Batch(
				helpers.GenCmd(messages.ChangeView{Name: names.BankCardList}),
				helpers.GenCmd(messages.AddBankCard{Item: f.getBankCardItem()}),
				helpers.GenCmd(messages.ShowMessage{Message: "Изменения сохранены"}),
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
	*f.item.UpdatedAt = time.Now()
	return f.item
}

func (f *bankCardForm) View() string {
	return f.baseForm.View()
}
