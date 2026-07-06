package storage

import "url-shortener/internal/model"

type MemoryStorage struct {
	links map[string]model.Link
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		links: make(map[string]model.Link),
	}
}

func (s *MemoryStorage) Save(link model.Link) {
	s.links[link.ID] = link
}
