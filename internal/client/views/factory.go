package views

import (
	"github.com/AndrXxX/goph-keeper/internal/client/state"
	"github.com/AndrXxX/goph-keeper/internal/client/views/forms"
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

func (f *Factory) AuthMenu() *authMenu {
	m := newAuthMenu()
	m.f = f.FormsFactory()
	return m
}
