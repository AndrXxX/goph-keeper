package lists

import (
	"github.com/AndrXxX/goph-keeper/internal/client/views/forms"
)

type Factory struct {
	FF *forms.Factory
}

func (f *Factory) AuthMenu() *authMenu {
	m := newAuthMenu()
	m.f = f.FF
	return m
}
