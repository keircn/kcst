package upload

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/keircn/kcst/internal/storage"
)

type Store struct {
	dir string
	db  *storage.DB
}

func NewStore(dir string, db *storage.DB) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Store{dir: dir, db: db}, nil
}

func (s *Store) Save(file multipart.File, header *multipart.FileHeader) (string, error) {
	randName, err := generateRandomName()
	if err != nil {
		return "", err
	}

	ext := getExtension(header.Filename)
	filename := randName + ext

	dst, err := os.Create(filepath.Join(s.dir, filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	meta := &storage.FileMetadata{
		ID:           randName,
		OriginalName: header.Filename,
		StoredName:   filename,
		Size:         size,
		ContentType:  header.Header.Get("Content-Type"),
		UploadedAt:   time.Now(),
	}
	if err := s.db.SaveMetadata(meta); err != nil {
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
