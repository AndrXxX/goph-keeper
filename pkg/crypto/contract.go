package crypto

type Decryptable interface {
	Encrypt(str string) (string, error)
	Decrypt(str string) (string, error)
}

type Verifiable interface {
	Encrypt(str string) (string, error)
	Verify(enc, decr string) (bool, error)
}
