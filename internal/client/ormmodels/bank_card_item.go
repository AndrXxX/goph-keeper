package ormmodels

type BankCardItem struct {
	StoredItem
	Number     string
	CVCCode    string
	Validity   string
	Cardholder string
}
