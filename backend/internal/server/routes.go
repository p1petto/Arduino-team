package server

import "github.com/go-chi/chi/v5"

func (s *Server) configureRoutes() {
	s.mux.Post("/login/{login}", s.handleApiKeyCreate)
	s.mux.Group(func(r chi.Router) {
		r.Use(s.AuthMiddleware())
		r.Post("/room/{name}", s.handleRoomCreate)
		r.Get("/room/{room}", s.handleRoomGet)
		r.Get("/room", s.handleRoomListGet)
		r.Get("/ws/{room}", s.handleWS)
	})
}
