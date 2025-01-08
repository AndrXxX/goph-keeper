package ormmodels

type PasswordItem struct {
	StoredItem
	Login    string
	Password string
}
