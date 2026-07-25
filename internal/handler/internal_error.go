package handler

import "url-shortener/internal/web"

func (h *Handler) InternalServerError(req web.Request) web.Response {
	return web.Response{
		StatusCode: 500,
		Status:     "Internal Server Error",
		Headers:    map[string]string{"Content-Type": "text/plain"},
		Body:       []byte("500 Internal Server Error"),
	}
}
