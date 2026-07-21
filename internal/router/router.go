package router

import (
	"url-shortener/internal/handler"
	"url-shortener/internal/web"
)

type Router struct {
	handler *handler.Handler
}

func New(h *handler.Handler) *Router {
	return &Router{
		handler: h,
	}
}

func (r *Router) Handle(req web.Request) web.Response {
	switch {
	case req.Method == "GET" && req.Path == "/":
		return r.handler.Home(req)

	case req.Method == "POST" && req.Path == "/links":
		return r.handler.CreateLink(req)

	case req.Method == "GET":
		return r.handler.Redirect(req)

	default:
		return r.handler.NotFound(req)
	}
}
