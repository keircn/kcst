package models

import "time"

type PageData struct {
	Title   string
	Message string
	BaseURL string
}

type FilePreviewData struct {
	Title        string
	Description  string
	Filename     string
	OriginalName string
	Size         int64
	SizeHuman    string
	ContentType  string
	MediaType    string
	RawURL       string
	PreviewURL   string
	BaseURL      string
	UploadedAt   time.Time
	ExpiresAt    time.Time
}

type UploadResponse struct {
	Success      bool      `json:"success"`
	URL          string    `json:"url"`
	RawURL       string    `json:"raw_url"`
	PreviewURL   string    `json:"preview_url"`
	Filename     string    `json:"filename"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	SizeHuman    string    `json:"size_human"`
	ContentType  string    `json:"content_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	RetentionMS  int64     `json:"retention_ms"`
	ResponseMS   int64     `json:"response_ms"`
}
