package values

type BankCardValue struct {
	Number     string `json:"number"`
	CVCCode    string `json:"cvc_code"`
	Validity   string `json:"validity"`
	Cardholder string `json:"cardholder"`
}
