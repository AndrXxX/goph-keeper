package views

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/forms"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/internal/client/views/styles"
	"github.com/AndrXxX/goph-keeper/pkg/entities"
)

var bankCardListKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Edit, kb.Delete, kb.New},
	Full: [][]key.Binding{
		{kb.Edit, kb.Delete, kb.New, kb.Quit},
		{kb.Up, kb.Down, kb.Enter, kb.Back},
	},
}

type bankCardList struct {
	list   list.Model
	help   help.Model
	height int
	width  int
}

func NewBankCardList() *bankCardList {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	defaultList.Title = "Bank cards"
	defaultList.Styles.Title = styles.Title
	return &bankCardList{list: defaultList, help: help.New()}
}

func (pl *bankCardList) Init() tea.Cmd {
	return nil
}

func (pl *bankCardList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		pl.setSize(msg.Width, msg.Height)
		pl.list.SetSize(msg.Width/margin, msg.Height/2)
	case messages.AddBankCard:
		pl.list.InsertItem(-1, msg.Item)
		pl.View()
		return pl, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Edit, kb.Keys.Enter):
			if len(pl.list.VisibleItems()) != 0 {
				e := pl.list.SelectedItem().(*entities.BankCardItem)
				f := forms.NewBankCardForm(e)
				return f, func() tea.Msg {
					return messages.ChangeView{
						Name: names.BankCardForm,
						View: f,
					}
				}
			}
		case key.Matches(msg, kb.Keys.New):
			f := forms.NewBankCardForm(nil)
			return f, func() tea.Msg {
				return messages.ChangeView{
					Name: names.BankCardForm,
					View: f,
				}
			}
		case key.Matches(msg, kb.Keys.Back):
			return pl, func() tea.Msg {
				return messages.ChangeView{
					Name: names.MainMenu,
				}
			}
		case key.Matches(msg, kb.Keys.Delete):
			return pl, pl.DeleteCurrent()
		}
	}
	pl.list, cmd = pl.list.Update(msg)
	return pl, cmd
}

func (pl *bankCardList) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, pl.list.View(), pl.help.View(bankCardListKeys))
}

func (pl *bankCardList) DeleteCurrent() tea.Cmd {
	if len(pl.list.VisibleItems()) > 0 {
		pl.list.RemoveItem(pl.list.Index())
	}

	var cmd tea.Cmd
	pl.list, cmd = pl.list.Update(nil)
	return cmd
}

func (pl *bankCardList) setSize(width, height int) {
	pl.width = width / margin
	pl.height = height
}
