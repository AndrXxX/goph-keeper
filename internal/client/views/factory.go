package views

import (
	"github.com/AndrXxX/goph-keeper/internal/client/state"
)

type Factory struct {
	AppState   *state.AppState
	Loginer    loginer
	Registerer registerer
}

func (f *Factory) AuthMenu() *authMenu {
	m := newAuthMenu()
	m.f = f
	return m
}

func (f *Factory) LoginForm() *loginForm {
	lf := newLoginForm()
	lf.f = f
	lf.s = f.AppState
	lf.l = f.Loginer
	return lf
}

func (f *Factory) RegisterForm() *registerForm {
	return newRegisterForm()
}

func (f *Factory) MasterPassForm() *masterPassForm {
	mpf := newMasterPassForm()
	mpf.s = f.AppState
	return mpf
}
