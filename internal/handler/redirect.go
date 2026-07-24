package handler

import (
	"log"
	"strings"
	"url-shortener/internal/web"
)

func (h *Handler) Redirect(req web.Request) web.Response {
	id := strings.TrimPrefix(req.Path, "/")
	link, err := h.storage.FindByID(id)
	if err != nil {
		return h.NotFound(req)
	}
	err = h.storage.IncrementClicks(id)
	if err != nil {
		log.Printf("Failed to increment clicks for ID %s: %v", id, err)
	}

	return web.Response{
		StatusCode: 302,
		Status:     "Found",
		Headers:    map[string]string{"Location": link.URL},
		Body:       nil,
	}
}
