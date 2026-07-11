package handler

import (
	"url-shortener/internal/model"
	"url-shortener/internal/web"
)

func (h *Handler) CreateLink(req web.Request) web.Response {
	link := model.Link{
		ID:  req.Path,         // Здесь нужно сгенерировать уникальный идентификатор для ссылки
		URL: string(req.Body), // Здесь нужно получить URL из тела запроса
	}
	err := h.storage.Save(link)
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

// тут будет работа с h.storage
