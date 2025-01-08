package ormmodels

type NoteItem struct {
	StoredItem
	Text string `json:"text"`
}
