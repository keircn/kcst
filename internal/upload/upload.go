package upload

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type Store struct {
	dir string
}

func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}

func (s *Store) Save(file multipart.File, originalName string) (string, error) {
	randName, err := generateRandomName()
	if err != nil {
		return "", err
	}

	ext := getExtension(originalName)
	filename := randName + ext

	dst, err := os.Create(filepath.Join(s.dir, filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return filename, nil
}

func generateRandomName() (string, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func getExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ".bin"
	}
	return strings.ToLower(ext)
}
