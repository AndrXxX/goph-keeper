package messages

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/AndrXxX/goph-keeper/internal/client/views/list"
)

type ChangeView struct {
	Name list.ViewName
	View tea.Model
	Msg  tea.Msg
}
