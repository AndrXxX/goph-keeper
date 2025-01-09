package hashgenerator

import (
	"crypto/sha256"
	"encoding/hex"
)

type sha256Generator struct {
	key string
}

// Generate выполняет генерацию хеша на основе ключа key и переданных данных data
func (h *sha256Generator) Generate(data []byte) string {
	t := sha256.New()
	t.Write(data)
	return hex.EncodeToString(t.Sum([]byte(h.key)))
}
