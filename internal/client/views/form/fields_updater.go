package form

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/AndrXxX/goph-keeper/internal/client/views/styles"
)

type FieldsUpdater struct {
}

func (u FieldsUpdater) Set(inputs []textinput.Model, focusIndex *int) []tea.Cmd {
	cmdList := make([]tea.Cmd, len(inputs))
	for i := 0; i <= len(inputs)-1; i++ {
		if i == *focusIndex {
			// Set focused state
			cmdList[i] = inputs[i].Focus()
			inputs[i].PromptStyle = styles.Focused
			inputs[i].TextStyle = styles.Focused
			continue
		}
		// Remove focused state
		inputs[i].Blur()
		inputs[i].PromptStyle = styles.Empty
		inputs[i].TextStyle = styles.Empty
	}
	return cmdList
}
