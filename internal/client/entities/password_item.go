package entities

type PasswordItem struct {
	StoredItem
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (si PasswordItem) FilterValue() string {
	return si.Login + si.Desc
}

func (si PasswordItem) Title() string {
	return si.Login
}

func (si PasswordItem) Description() string {
	return si.Desc
}
