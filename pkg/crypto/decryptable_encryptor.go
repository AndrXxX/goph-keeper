package crypto

import "github.com/dromara/dongle"

type decryptableEncryptor struct {
	cipher *dongle.Cipher
}

func (m *decryptableEncryptor) Encrypt(str string) (string, error) {
	e := dongle.Encrypt.FromString(str).ByAes(m.cipher)
	if e.Error != nil {
		return "", e.Error
	}
	return e.ToRawString(), nil
}

func (m *decryptableEncryptor) Decrypt(str string) (string, error) {
	d := dongle.Decrypt.FromRawString(str).ByAes(m.cipher)
	if d.Error != nil {
		return "", d.Error
	}
	return d.ToString(), nil
}
