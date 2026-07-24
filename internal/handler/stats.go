package handler

import (
	"encoding/json"
	"log"
	"strings"
	"url-shortener/internal/web"
)

func (h *Handler) Stats(req web.Request) web.Response {
	id := strings.TrimPrefix(req.Path, "/stats/")
	link, err := h.storage.FindByID(id)
	if err != nil {
		return h.NotFound(req)
	}
	body, err := json.Marshal(link)
	if err != nil {
		log.Printf("Failed to marshal link data for ID %s: %v", id, err)
		return h.InternalServerError(req)
	}
	return web.Response{
		StatusCode: 200,
		Status:     "OK",
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       body,
	}
}
