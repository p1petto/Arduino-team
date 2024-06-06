package server

import "github.com/go-chi/chi/v5"

func (s *Server) configureRoutes() {
	s.mux.Post("/login/{login}", s.handleApiKeyCreate)
	s.mux.Group(func(r chi.Router) {
		r.Use(s.AuthMiddleware())
		r.Post("/rooms", s.handleRoomCreate)
		r.Get("/rooms/{room}", s.handleRoomGet)
		r.Get("/rooms", s.handleRoomListGet)
		r.Get("/ws/{room}", s.handleWS)
	})
}
