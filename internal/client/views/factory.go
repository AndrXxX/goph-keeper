package views

import (
	"github.com/AndrXxX/goph-keeper/internal/client/state"
	"github.com/AndrXxX/goph-keeper/internal/client/views/forms"
	"github.com/AndrXxX/goph-keeper/internal/client/views/lists"
)

type Factory struct {
	AppState   *state.AppState
	Loginer    forms.Loginer
	Registerer forms.Registerer
}

func (f *Factory) FormsFactory() *forms.Factory {
	return &forms.Factory{
		AppState:   f.AppState,
		Loginer:    f.Loginer,
		Registerer: f.Registerer,
	}
}

func (f *Factory) MenusFactory() *lists.Factory {
	return &lists.Factory{
		FF: f.FormsFactory(),
	}
}
