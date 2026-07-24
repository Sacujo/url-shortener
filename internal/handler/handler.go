package handler

import (
	"url-shortener/internal/storage"
)

type Handler struct {
	storage storage.Storage
}

func New(store storage.Storage) *Handler {
	return &Handler{
		storage: store,
	}
}
