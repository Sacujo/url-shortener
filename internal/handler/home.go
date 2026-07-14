package handler

import "url-shortener/internal/web"

func (h *Handler) Home(req web.Request) web.Response {
	return web.Response{
		StatusCode: 200,
		Status:     "OK",
		Headers:    map[string]string{"Content-Type": "text/plain"},
		Body:       []byte("Welcome to the URL Shortener!"),
	}
}
