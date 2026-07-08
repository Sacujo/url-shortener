package handler

import (
	"url-shortener/internal/storage"
)

type Handler struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}
