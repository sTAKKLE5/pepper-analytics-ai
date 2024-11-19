package services

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

var (
	ErrPlantNotFound = errors.New("plant not found")
)

type FileService struct {
	uploadDir string
}

func NewFileService(uploadDir string) *FileService {
	return &FileService{uploadDir: uploadDir}
}

func (s *FileService) SaveFile(file io.Reader, destPath string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	dst, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	return err
}

func (s *FileService) DeleteFile(path string) error {
	return os.Remove(path)
}
