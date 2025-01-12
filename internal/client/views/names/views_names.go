package names

type ViewName string

const (
	PasswordList       = ViewName("passwordList")
	PasswordForm       = ViewName("passwordForm")
	NotesList          = ViewName("notesList")
	NoteForm           = ViewName("noteForm")
	BankCardList       = ViewName("bankCardList")
	BankCardForm       = ViewName("bankCardForm")
	FileList           = ViewName("fileList")
	UploadFileForm     = ViewName("uploadFileForm")
	MainMenu           = ViewName("mainMenu")
	AuthMenu           = ViewName("authMenu")
	LoginForm          = ViewName("loginForm")
	RegisterForm       = ViewName("registerForm")
	MasterPassRegForm  = ViewName("masterPassRegForm")
	MasterPassAuthForm = ViewName("masterPassAuthForm")
)
