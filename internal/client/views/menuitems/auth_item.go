package menuitems

type AuthItem struct {
	Name string
	Code string
	Desc string
}

func (i AuthItem) FilterValue() string {
	return i.Name + i.Desc
}

func (i AuthItem) Title() string {
	return i.Name
}

func (i AuthItem) Description() string {
	return i.Desc
}
