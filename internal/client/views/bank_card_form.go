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

	m.baseForm.inputs[0].Prompt = "Number: "
	m.baseForm.inputs[0].SetValue(m.item.Number)

	m.baseForm.inputs[1].Prompt = "CVCCode: "
	m.baseForm.inputs[1].SetValue(m.item.CVCCode)

	m.baseForm.inputs[2].Prompt = "Validity: "
	m.baseForm.inputs[2].SetValue(m.item.Validity)

	m.baseForm.inputs[3].Prompt = "Cardholder: "
	m.baseForm.inputs[3].SetValue(m.item.Cardholder)

	m.baseForm.inputs[4].Prompt = "Description: "
	m.baseForm.inputs[4].SetValue(m.item.Desc)

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
			return f, func() tea.Msg {
				return messages.ChangeView{
					Name: names.BankCardList,
				}
			}
		case key.Matches(msg, kb.Keys.Save):
			// TODO: сделать уведомление
			var nMsg tea.Msg
			if f.creating {
				nMsg = messages.AddBankCard{
					Item: f.getBankCardItem(),
				}
			}
			return f, func() tea.Msg {
				return messages.ChangeView{
					Name: names.BankCardList,
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

func (f *bankCardForm) getBankCardItem() *entities.BankCardItem {
	f.item.Number = f.baseForm.inputs[0].Value()
	f.item.CVCCode = f.baseForm.inputs[1].Value()
	f.item.Validity = f.baseForm.inputs[2].Value()
	f.item.Cardholder = f.baseForm.inputs[3].Value()
	f.item.Desc = f.baseForm.inputs[4].Value()
	return f.item
}

func (f *bankCardForm) View() string {
	return f.baseForm.View()
}
