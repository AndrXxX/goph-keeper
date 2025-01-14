package forms

import (
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/internal/client/views/styles"
)

var downloadFileFormKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Download},
	Full: [][]key.Binding{
		{kb.Back, kb.Save, kb.Quit},
		{kb.Up, kb.Down, kb.Enter},
	},
}

type downloadFileForm struct {
	help         help.Model
	title        string
	item         *entities.FileItem
	filePicker   filepicker.Model
	selectedPath string
	keys         *kb.KeyMap
	height       int
	inited       bool
}

func NewDownloadFileForm(item *entities.FileItem, height int) *downloadFileForm {
	m := downloadFileForm{
		help:       help.New(),
		title:      "Download file",
		item:       item,
		filePicker: filepicker.New(),
		height:     height,
	}
	m.filePicker.CurrentDirectory, _ = os.UserHomeDir()
	m.filePicker.AutoHeight = false
	m.filePicker.DirAllowed = true
	m.filePicker.FileAllowed = false
	m.keys = &downloadFileFormKeys
	return &m
}

func (f *downloadFileForm) Init() tea.Cmd {
	return tea.Batch(helpers.GenCmd(tea.WindowSizeMsg{Height: f.height}), f.filePicker.Init())
}

func (f *downloadFileForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdList []tea.Cmd
	if !f.inited {
		cmdList = append(cmdList, f.Init())
		f.inited = true
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.filePicker.Height = msg.Height / 2
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, kb.Keys.Back):
			return f, helpers.GenCmd(messages.ChangeView{Name: names.FileList})
		case key.Matches(msg, kb.Keys.Download):
			if f.selectedPath == "" {
				return f, helpers.GenCmd(messages.ShowError{Err: "Для скачивания нужно выбрать путь"})
			}
			return f, tea.Sequence(
				helpers.GenCmd(messages.DownloadFile{Path: f.selectedPath, Item: f.item}),
				helpers.GenCmd(messages.ShowMessage{Message: "Выполняется скачивание файла"}),
				helpers.GenCmd(messages.ChangeView{Name: names.FileList}),
				tea.WindowSize(),
			)
		}
	}
	fp, cmd := f.filePicker.Update(msg)
	f.filePicker = fp
	cmdList = append(cmdList, cmd)
	if didSelect, filePath := f.filePicker.DidSelectDisabledFile(msg); didSelect {
		f.selectedPath = ""
		return f, helpers.GenCmd(messages.ShowError{Err: filePath + " is not valid"})
	}
	if didSelect, filePath := f.filePicker.DidSelectFile(msg); didSelect {
		f.selectedPath = filePath
	}
	return f, tea.Sequence(cmdList...)
}

func (f *downloadFileForm) View() string {
	vList := []string{
		styles.Title.Margin(1).Render(f.title),
	}
	vList = append(vList, "File: "+f.filePicker.Styles.Selected.Render(f.item.Name)+"\n")
	if f.selectedPath != "" {
		vList = append(vList, "Selected path: "+f.filePicker.Styles.Selected.Render(f.selectedPath)+"\n")
	}
	vList = append(vList, f.filePicker.View())
	if f.keys != nil {
		vList = append(vList, f.help.View(*f.keys))
	}
	return lipgloss.JoinVertical(lipgloss.Left, vList...)
}
