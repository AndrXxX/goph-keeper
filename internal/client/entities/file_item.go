package entities

type FileItem struct {
	StoredItem
	Data string `json:"data"`
}

func (i FileItem) FilterValue() string {
	return i.Desc
}

func (i FileItem) Title() string {
	return i.Desc
}

func (i FileItem) Description() string {
	return ""

}
