package filestorage

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/google/uuid"
)

type storage struct {
	path string
	hg   hashGenerator
}

func New(path string, hg hashGenerator) (*storage, error) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("make dir %w", err)
	}
	return &storage{path, hg}, nil
}

func (s *storage) Store(src io.Reader, id uuid.UUID) error {
	fullPath := path.Join(s.path, s.FileId(id))
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	_, err = io.Copy(file, src)
	return err
}

func (s *storage) Get(id uuid.UUID) (file io.ReadCloser, err error) {
	fullPath := path.Join(s.path, s.FileId(id))
	return os.OpenFile(fullPath, os.O_RDONLY, os.ModePerm)
}

func (s *storage) IsExist(id uuid.UUID) bool {
	fullPath := path.Join(s.path, s.FileId(id))
	_, err := os.Stat(fullPath)
	return err == nil || !os.IsNotExist(err)
}

func (s *storage) Delete(id uuid.UUID) (err error) {
	fullPath := path.Join(s.path, s.FileId(id))
	return os.Remove(fullPath)
}

func (s *storage) FileId(id uuid.UUID) string {
	return s.hg.Generate([]byte(id.String()))
}
