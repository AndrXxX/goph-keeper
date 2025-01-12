package forms

import (
	"os"
	"path"
	"path/filepath"
	"time"

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

var uploadFileFormKeys = kb.KeyMap{
	Short: []key.Binding{kb.Back, kb.Save},
	Full: [][]key.Binding{
		{kb.Back, kb.Save, kb.Quit},
		{kb.Up, kb.Down, kb.Enter},
	},
}

type uploadFileForm struct {
	help         help.Model
	title        string
	item         *entities.FileItem
	creating     bool
	filePicker   filepicker.Model
	selectedFile string
	keys         *kb.KeyMap
	height       int
	inited       bool
}

func NewUploadFileForm(item *entities.FileItem, height int) *uploadFileForm {
	m := uploadFileForm{
		help:       help.New(),
		title:      "Upload file",
		creating:   item == nil,
		item:       item,
		filePicker: filepicker.New(),
		height:     height,
	}
	m.filePicker.CurrentDirectory, _ = os.UserHomeDir()
	if !m.creating {
		if item.Path != "" {
			m.filePicker.CurrentDirectory = path.Dir(item.Path)
			m.selectedFile = item.Path
		}
	}
	m.filePicker.AutoHeight = false
	m.keys = &uploadFileFormKeys
	if m.creating {
		m.item = &entities.FileItem{}
	}
	return &m
}

func (f *uploadFileForm) Init() tea.Cmd {
	return tea.Batch(helpers.GenCmd(tea.WindowSizeMsg{Height: f.height}), f.filePicker.Init())
}

func (f *uploadFileForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case key.Matches(msg, kb.Keys.Save):
			if f.selectedFile == "" {
				return f, helpers.GenCmd(messages.ShowError{Err: "Для сохранения нужно выбрать файл"})
			}
			return f, helpers.GenCmd(messages.ChangeView{Name: names.UpdateFileForm, View: NewUpdateFileForm(f.getFileItem())})
		}
	}
	fp, cmd := f.filePicker.Update(msg)
	f.filePicker = fp
	cmdList = append(cmdList, cmd)
	if didSelect, filePath := f.filePicker.DidSelectDisabledFile(msg); didSelect {
		f.selectedFile = ""
		return f, helpers.GenCmd(messages.ShowError{Err: filePath + " is not valid"})
	}
	if didSelect, filePath := f.filePicker.DidSelectFile(msg); didSelect {
		f.selectedFile = filePath
	}
	return f, tea.Batch(cmdList...)
}

func (f *uploadFileForm) getFileItem() *entities.FileItem {
	f.item.Name = filepath.Base(f.selectedFile)
	f.item.Path = f.selectedFile
	f.item.UpdatedAt = time.Now()
	return f.item
}

func (f *uploadFileForm) View() string {
	vList := []string{
		styles.Title.Margin(1).Render(f.title),
	}
	if f.selectedFile != "" {
		vList = append(vList, "Selected file: "+f.filePicker.Styles.Selected.Render(f.selectedFile))
	}
	vList = append(vList, f.filePicker.View())
	if f.keys != nil {
		vList = append(vList, f.help.View(*f.keys))
	}
	return lipgloss.JoinVertical(lipgloss.Left, vList...)
}
