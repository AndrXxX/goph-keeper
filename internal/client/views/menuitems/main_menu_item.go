package menuitems

type MainMenuItem struct {
	Name string
	Code string
	Desc string
}

func (i MainMenuItem) FilterValue() string {
	return i.Name + i.Desc
}

func (i MainMenuItem) Title() string {
	return i.Name
}

func (i MainMenuItem) Description() string {
	return i.Desc
}
