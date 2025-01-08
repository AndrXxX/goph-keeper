package entities

type User struct {
	Login    string `json:"login" valid:"minstringlength(3),required"`
	Password string `json:"password" valid:"minstringlength(5),required"`
}
