package entities

type BankCardItem struct {
	StoredItem
	Number     string `json:"number"`
	CVCCode    string `json:"cvc_code"`
	Validity   string `json:"validity"`
	Cardholder string `json:"cardholder"`
}

func (i BankCardItem) FilterValue() string {
	return i.Number + i.Desc
}

func (i BankCardItem) Title() string {
	return i.Number
}

func (i BankCardItem) Description() string {
	return i.Desc
}
