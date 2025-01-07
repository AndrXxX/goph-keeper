package names

type ViewName string

const (
	PasswordList = ViewName("passwordList")
	PasswordForm = ViewName("passwordForm")
	NotesList    = ViewName("notesList")
	NoteForm     = ViewName("noteForm")
	BankCardList = ViewName("bankCardList")
	BankCardForm = ViewName("bankCardForm")
	FileList     = ViewName("fileList")
	FileForm     = ViewName("fileForm")
	MainMenu     = ViewName("mainMenu")
	AuthMenu     = ViewName("authMenu")
	LoginForm    = ViewName("loginForm")
)
