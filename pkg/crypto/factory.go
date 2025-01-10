package crypto

import "github.com/dromara/dongle"

type factory struct {
	key string
}

func Factory(key string) *factory {
	return &factory{key: key}
}

func (f *factory) NewDecryptableEncryptor() Decryptable {
	cipher := dongle.NewCipher()
	cipher.SetMode(dongle.CBC)      // CBC、CFB、OFB、CTR、ECB
	cipher.SetPadding(dongle.PKCS7) // No、Empty、Zero、PKCS5、PKCS7、AnsiX923、ISO97971
	cipher.SetKey(f.key)            // key must be 16, 24 or 32 bytes
	return &decryptableEncryptor{cipher: cipher}
}

func (f *factory) NewVerifiableEncryptor() Verifiable {
	return &verifiableEncryptor{key: f.key}
}
