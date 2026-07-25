package storage

import (
	"errors"
	"sync"
	"url-shortener/internal/model"
)

var ErrNotFound = errors.New("link not found")
var ErrAlreadyExists = errors.New("link already exists")

type MemoryStorage struct {
	mu    sync.RWMutex
	links map[string]model.Link
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		links: make(map[string]model.Link),
	}
}

func (s *MemoryStorage) Save(link model.Link) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.links[link.ID]; exists {
		return ErrAlreadyExists
	}
	s.links[link.ID] = link
	return nil
}

func (s *MemoryStorage) FindByID(id string) (model.Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	link, exists := s.links[id]
	if !exists {
		return model.Link{}, ErrNotFound
	}
	return link, nil
}

func (s *MemoryStorage) IncrementClicks(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	link, exists := s.links[id]
	if !exists {
		return ErrNotFound
	}
	link.Clicks++
	s.links[id] = link
	return nil
}
