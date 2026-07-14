package handler

import (
	"math/rand/v2"
	"strings"
	"url-shortener/internal/model"
	"url-shortener/internal/storage"
	"url-shortener/internal/web"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const idLength = 8

func (h *Handler) CreateLink(req web.Request) web.Response {
	url := strings.TrimSpace(string(req.Body))
	if url == "" {
		return web.Response{
			StatusCode:  400,
			Status:      "Bad Request",
			ContentType: "text/plain",
			Body:        "400 Bad Request: URL is required",
		}
	}
	for {
		link := model.Link{
			ID:  generateShortID(), // Здесь нужно сгенерировать уникальный идентификатор для ссылки
			URL: url,               // Здесь нужно получить URL из тела запроса
		}
		err := h.storage.Save(link)
		if err == storage.ErrAlreadyExists {
			continue // Если идентификатор уже существует, генерируем новый

		}
		if err != nil {
			return web.Response{
				StatusCode:  500,
				Status:      "Internal Server Error",
				ContentType: "text/plain",
				Body:        "500 Internal Server Error",
			}
		}
		return web.Response{
			StatusCode:  201,
			Status:      "Created",
			ContentType: "text/plain",
			Body:        link.ID,
		}
	}
}

func generateShortID() string {

	var b strings.Builder
	b.Grow(idLength)

	for i := 0; i < idLength; i++ {
		b.WriteByte(chars[rand.IntN(len(chars))])
	}

	return b.String()
}
