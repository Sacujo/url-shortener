package handler

import "url-shortener/internal/web"

func (h *Handler) NotFound(req web.Request) web.Response {
	// поиск ссылки через h.storage
	return web.Response{
		StatusCode: 404,
		Status:     "Not Found",
		Headers:    map[string]string{"Content-Type": "text/plain"},
		Body:       []byte("404 Not Found"),
	}
}
