package store

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Event struct {
	Source     string         `json:"source"`
	EventType  string         `json:"event_type"`
	UserID     *int64         `json:"user_id,omitempty"`
	TenantID   string         `json:"tenant_id,omitempty"`
	Provider   string         `json:"provider,omitempty"`
	IP         string         `json:"ip,omitempty"`
	UserAgent  string         `json:"user_agent,omitempty"`
	Success    bool           `json:"success"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	ReceivedAt time.Time      `json:"received_at"`
}

type FileStore struct {
	path string
	mu   sync.Mutex
}

func NewFileStore(path string) (*FileStore, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	return &FileStore{path: path}, nil
}

func (s *FileStore) Append(ctx context.Context, event Event) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	return enc.Encode(event)
}
