package names

type ViewName string

const (
	PasswordList  = ViewName("passwordList")
	PasswordForm  = ViewName("passwordForm")
	NotesList     = ViewName("notesList")
	NoteForm      = ViewName("noteForm")
	BankCardsList = ViewName("bankCardsList")
	BankCardsForm = ViewName("bankCardsForm")
	FilesList     = ViewName("filesList")
	FilesForm     = ViewName("filesForm")
	TypesList     = ViewName("typesList")
	AuthMenu      = ViewName("authMenu")
	AuthForm      = ViewName("authForm")
)
