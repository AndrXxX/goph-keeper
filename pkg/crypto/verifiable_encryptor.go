package crypto

import "github.com/dromara/dongle"

type verifiableEncryptor struct {
	key string
}

func (ve *verifiableEncryptor) Encrypt(str string) (string, error) {
	e := dongle.Encrypt.FromString(str).ByHmacSha256([]byte(ve.key))
	if e.Error != nil {
		return "", e.Error
	}
	return e.ToBase64String(), nil
}

func (ve *verifiableEncryptor) Verify(enc, decr string) (bool, error) {
	newEnc, err := ve.Encrypt(decr)
	if err != nil {
		return false, err
	}
	return newEnc == enc, nil
}
