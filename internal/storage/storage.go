package storage

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	MaxFileSize = 100 * 1024 * 1024
	MinTTL      = 1 * time.Hour
	MaxTTL      = 28 * 24 * time.Hour
)

type FileMetadata struct {
	ID           string    `json:"id"`
	OriginalName string    `json:"original_name"`
	StoredName   string    `json:"stored_name"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

func (f *FileMetadata) ExpiresAt() time.Time {
	ttl := CalculateTTL(f.Size)
	return f.UploadedAt.Add(ttl)
}

func (f *FileMetadata) IsExpired() bool {
	return time.Now().After(f.ExpiresAt())
}

// TTL = MinTTL + (MaxTTL - MinTTL) * (1 - (size/maxSize)^0.5)
func CalculateTTL(size int64) time.Duration {
	if size <= 0 {
		return MaxTTL
	}
	if size >= MaxFileSize {
		return MinTTL
	}

	ratio := float64(size) / float64(MaxFileSize)
	factor := 1 - math.Sqrt(ratio)
	ttlRange := float64(MaxTTL - MinTTL)
	ttl := MinTTL + time.Duration(factor*ttlRange)

	return ttl
}

type DB struct {
	path string
	mu   sync.RWMutex
	data map[string]*FileMetadata
}

func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	db := &DB{
		path: path,
		data: make(map[string]*FileMetadata),
	}

	if _, err := os.Stat(path); err == nil {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		if err := json.NewDecoder(file).Decode(&db.data); err != nil {
			db.data = make(map[string]*FileMetadata)
		}
	}

	return db, nil
}

func (d *DB) Close() error {
	return nil
}

func (d *DB) save() error {
	file, err := os.Create(d.path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(d.data)
}

func (d *DB) SaveMetadata(meta *FileMetadata) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.data[meta.ID] = meta
	return d.save()
}

func (d *DB) GetMetadata(id string) (*FileMetadata, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	meta, ok := d.data[id]
	if !ok {
		return nil, nil
	}
	return meta, nil
}

func (d *DB) ListMetadata() ([]*FileMetadata, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	results := make([]*FileMetadata, 0, len(d.data))
	for _, meta := range d.data {
		results = append(results, meta)
	}
	return results, nil
}

func (d *DB) DeleteMetadata(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.data, id)
	return d.save()
}

func (d *DB) GetExpired() ([]*FileMetadata, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var expired []*FileMetadata
	for _, meta := range d.data {
		if meta.IsExpired() {
			expired = append(expired, meta)
		}
	}
	return expired, nil
}
