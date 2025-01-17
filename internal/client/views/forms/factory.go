package forms

type Factory struct {
	Loginer    Loginer
	Registerer Registerer
}

func (f *Factory) LoginForm() *loginForm {
	lf := newLoginForm()
	lf.f = f
	lf.l = f.Loginer
	return lf
}

func (f *Factory) RegisterForm() *registerForm {
	rf := newRegisterForm()
	rf.r = f.Registerer
	rf.f = f
	return rf
}

func (f *Factory) MasterPassRegForm() *masterPassRegForm {
	mpf := newMasterPassRegForm()
	return mpf
}

func (f *Factory) MasterPassAuthForm() *masterPassAuthForm {
	mpf := newMasterPassAuthForm()
	return mpf
}
