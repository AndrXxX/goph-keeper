package views

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/AndrXxX/goph-keeper/internal/client/views/contract"
	"github.com/AndrXxX/goph-keeper/internal/client/views/forms"
	"github.com/AndrXxX/goph-keeper/internal/client/views/lists"
)

type Factory struct {
	Loginer    forms.Loginer
	Registerer forms.Registerer
	S          *contract.Storages
}

func (f *Factory) Container(opts ...Option) *container {
	c := &container{
		help: help.New(),
		uo:   make(map[tea.Msg]UpdateOption),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (f *Factory) FormsFactory() *forms.Factory {
	return &forms.Factory{
		Loginer:    f.Loginer,
		Registerer: f.Registerer,
	}
}

func (f *Factory) MenusFactory() *lists.Factory {
	return &lists.Factory{
		FF: f.FormsFactory(),
		S:  f.S,
	}
}
