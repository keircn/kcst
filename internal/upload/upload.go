package upload

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/keircn/kcst/internal/storage"
)

type Store struct {
	dir             string
	db              *storage.DB
	cleanupInterval time.Duration
}

func NewStore(dir string, db *storage.DB, cleanupInterval time.Duration) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Store{dir: dir, db: db, cleanupInterval: cleanupInterval}, nil
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

func (s *Store) Get(filename string) (*os.File, *storage.FileMetadata, error) {
	meta, err := s.db.GetMetadataByStoredName(filename)
	if err != nil {
		return nil, nil, err
	}
	if meta == nil {
		return nil, nil, os.ErrNotExist
	}
	if meta.IsExpired() {
		return nil, nil, os.ErrNotExist
	}

	file, err := os.Open(filepath.Join(s.dir, meta.StoredName))
	if err != nil {
		return nil, nil, err
	}

	return file, meta, nil
}

func (s *Store) Cleanup() error {
	expired, err := s.db.GetExpired()
	if err != nil {
		return err
	}

	for _, meta := range expired {
		filePath := filepath.Join(s.dir, meta.StoredName)
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			log.Printf("Failed to remove expired file %s: %v", filePath, err)
			continue
		}

		if err := s.db.DeleteMetadata(meta.ID); err != nil {
			log.Printf("Failed to delete metadata for %s: %v", meta.ID, err)
			continue
		}

		log.Printf("Cleaned up expired file: %s (size: %d, uploaded: %s)",
			meta.StoredName, meta.Size, meta.UploadedAt.Format(time.RFC3339))
	}

	return nil
}

func (s *Store) StartCleanupRoutine(stop <-chan struct{}) {
	ticker := time.NewTicker(s.cleanupInterval)
	go func() {
		if err := s.Cleanup(); err != nil {
			log.Printf("Cleanup error: %v", err)
		}

		for {
			select {
			case <-ticker.C:
				if err := s.Cleanup(); err != nil {
					log.Printf("Cleanup error: %v", err)
				}
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
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
