package lists

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/views/forms"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/internal/client/views/styles"
)

var fileListKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Edit, kb.Delete, kb.New, kb.Download},
	Full: [][]key.Binding{
		{kb.Edit, kb.Delete, kb.New, kb.Quit},
		{kb.Up, kb.Down, kb.Enter, kb.Back},
	},
}

type fileList struct {
	list list.Model
	help help.Model
	lr   refresher
}

func newFileList() *fileList {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	defaultList.Title = "Files"
	defaultList.Styles.Title = styles.Title
	return &fileList{list: defaultList, help: help.New()}
}

func (l *fileList) Init() tea.Cmd {
	l.lr.Refresh()
	return nil
}

func (l *fileList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(l.list.Items()) == 0 {
		l.lr.Refresh()
	}
	l.lr.RefreshIn(refreshListInterval)
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.list.SetSize(msg.Width/styles.InnerMargin, msg.Height/2)
	case messages.AddFile:
		l.lr.Refresh()
		return l, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Edit, kb.Keys.Enter):
			if len(l.list.VisibleItems()) != 0 {
				e := l.list.SelectedItem().(entities.FileItem)
				f := forms.NewUpdateFileForm(&e)
				return f, helpers.GenCmd(messages.ChangeView{Name: names.UpdateFileForm, View: f})
			}
		case key.Matches(msg, kb.Keys.New):
			f := forms.NewUploadFileForm(&entities.FileItem{}, l.list.Height()*2)
			return f, helpers.GenCmd(messages.ChangeView{Name: names.UploadFileForm, View: f})
		case key.Matches(msg, kb.Keys.Download):
			if len(l.list.VisibleItems()) != 0 {
				e := l.list.SelectedItem().(entities.FileItem)
				f := forms.NewDownloadFileForm(&e, l.list.Height()*2)
				return f, helpers.GenCmd(messages.ChangeView{Name: names.DownloadFileForm, View: f})
			}
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

func (l *fileList) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, l.list.View(), l.help.View(fileListKeys))
}

func (l *fileList) DeleteCurrent() tea.Cmd {
	if len(l.list.VisibleItems()) > 0 {
		l.list.RemoveItem(l.list.Index())
	}

	var cmd tea.Cmd
	l.list, cmd = l.list.Update(nil)
	return cmd
}
