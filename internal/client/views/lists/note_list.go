package lists

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/views/forms"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	kb "github.com/AndrXxX/goph-keeper/internal/client/views/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/views/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/internal/client/views/styles"
)

var noteListKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Edit, kb.New},
	Full: [][]key.Binding{
		{kb.Edit, kb.New, kb.Quit},
		{kb.Up, kb.Down, kb.Enter, kb.Back},
	},
}

type noteList struct {
	list list.Model
	help help.Model
	lr   refresher
}

func newNoteList() *noteList {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	defaultList.Title = "Notes"
	defaultList.Styles.Title = styles.Title
	return &noteList{list: defaultList, help: help.New()}
}

func (l *noteList) Init() tea.Cmd {
	l.lr.Refresh()
	return nil
}

func (l *noteList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(l.list.Items()) == 0 {
		l.lr.Refresh()
	}
	l.lr.RefreshIn(refreshListInterval)
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.list.SetSize(msg.Width, msg.Height/2)
	case messages.AddNote:
		l.lr.Refresh()
		return l, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Edit, kb.Keys.Enter):
			if len(l.list.VisibleItems()) != 0 {
				e := l.list.SelectedItem().(entities.NoteItem)
				f := forms.NewNoteForm(&e)
				return f, helpers.GenCmd(messages.ChangeView{Name: names.NoteForm, View: f})
			}
		case key.Matches(msg, kb.Keys.New):
			f := forms.NewNoteForm(nil)
			return f, helpers.GenCmd(messages.ChangeView{Name: names.NoteForm, View: f})
		case key.Matches(msg, kb.Keys.Back):
			return l, helpers.GenCmd(messages.ChangeView{Name: names.MainMenu})
		case key.Matches(msg, kb.Keys.Delete):
			// TODO: approve + action
			return l, l.DeleteCurrent()
		}
	}
	l.list, cmd = l.list.Update(msg)
	return l, cmd
}

func (l *noteList) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, l.list.View(), styles.Help.Render(l.help.View(noteListKeys)))
}

func (l *noteList) DeleteCurrent() tea.Cmd {
	if len(l.list.VisibleItems()) > 0 {
		l.list.RemoveItem(l.list.Index())
	}

	var cmd tea.Cmd
	l.list, cmd = l.list.Update(nil)
	return cmd
}
