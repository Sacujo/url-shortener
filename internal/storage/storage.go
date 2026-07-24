package storage

import "url-shortener/internal/model"

type Storage interface {
	Save(link model.Link) error
	FindByID(id string) (model.Link, error)
	IncrementClicks(id string) error
}
