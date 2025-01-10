package forms

import (
	"github.com/AndrXxX/goph-keeper/internal/client/state"
	"github.com/AndrXxX/goph-keeper/internal/client/views/contract"
)

type Factory struct {
	AppState   *state.AppState
	Loginer    Loginer
	Registerer Registerer
	SM         contract.SyncManager
}

func (f *Factory) LoginForm() *loginForm {
	lf := newLoginForm()
	lf.f = f
	lf.s = f.AppState
	lf.l = f.Loginer
	return lf
}

func (f *Factory) RegisterForm() *registerForm {
	rf := newRegisterForm()
	rf.r = f.Registerer
	rf.s = f.AppState
	rf.f = f
	return rf
}

func (f *Factory) MasterPassRegForm() *masterPassRegForm {
	mpf := newMasterPassRegForm()
	mpf.s = f.AppState
	return mpf
}

func (f *Factory) MasterPassAuthForm() *masterPassAuthForm {
	mpf := newMasterPassAuthForm()
	return mpf
}
