package handler

import "url-shortener/internal/web"

func (h *Handler) NotFound(req web.Request) web.Response {
	// поиск ссылки через h.storage
	return web.Response{
		StatusCode:  404,
		Status:      "Not Found",
		ContentType: "text/plain",
		Body:        "404 Not Found",
	}
}
