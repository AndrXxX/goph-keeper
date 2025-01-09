package views

import (
	"github.com/charmbracelet/bubbles/help"

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

func (f *Factory) Container() *container {
	return &container{help: help.New(), views: NewMap(f), qr: f.QR, sm: f.SM}
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
