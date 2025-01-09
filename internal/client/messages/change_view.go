package messages

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
)

type ChangeView struct {
	Name names.ViewName
	View tea.Model
	Msg  tea.Msg
}
