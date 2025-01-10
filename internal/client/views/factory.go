package views

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/AndrXxX/goph-keeper/internal/client/state"
	"github.com/AndrXxX/goph-keeper/internal/client/views/contract"
	"github.com/AndrXxX/goph-keeper/internal/client/views/forms"
	"github.com/AndrXxX/goph-keeper/internal/client/views/lists"
)

type Factory struct {
	AppState   *state.AppState
	Loginer    forms.Loginer
	Registerer forms.Registerer
	SM         contract.SyncManager
	S          *contract.Storages
	QR         contract.QueueRunner
}

func (f *Factory) Container(opts ...Option) *container {
	c := &container{
		help:  help.New(),
		views: NewMap(f),
		qr:    f.QR,
		sm:    f.SM,
		as:    f.AppState,
		uo:    make(map[tea.Msg]UpdateOption),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (f *Factory) FormsFactory() *forms.Factory {
	return &forms.Factory{
		AppState:   f.AppState,
		Loginer:    f.Loginer,
		Registerer: f.Registerer,
		SM:         f.SM,
	}
}

func (f *Factory) MenusFactory() *lists.Factory {
	return &lists.Factory{
		FF: f.FormsFactory(),
		SM: f.SM,
		S:  f.S,
	}
}
