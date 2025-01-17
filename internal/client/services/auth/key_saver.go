package auth

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

type KeySaver struct {
	KeyPath string
}

func (s *KeySaver) Store(resp *http.Response) error {
	if s.KeyPath == "" {
		return nil
	}
	key, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Info("read server key", zap.Error(err))
		return nil
	}
	err = os.MkdirAll(filepath.Dir(s.KeyPath), 0755)
	if err != nil {
		logger.Log.Info("create server key", zap.Error(err))
		return nil
	}
	err = os.WriteFile(s.KeyPath, key, 0755)
	if err != nil {
		logger.Log.Info("create server key", zap.Error(err))
	}
	return nil
}
