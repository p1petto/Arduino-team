package server

import (
	"arduinoteam/internal/hub"
	"arduinoteam/storage/sqlite"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Server struct {
	mux      *chi.Mux
	upgrader websocket.Upgrader
	hub      *hub.Hub
	log      *slog.Logger
	storage  *sqlite.Storage
}

func NewServer() *Server {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	mux := chi.NewRouter()
	storage := sqlite.NewStorage()

	hub := hub.NewHub(storage, log)
	hub.Run()

	server := &Server{upgrader: upgrader, hub: hub, log: log, storage: storage, mux: mux}
	server.configureRoutes()

	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Run(port string) {
	log.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, s))
}

func (s *Server) CheckToken(token string) (*hub.Client, error) {
	var v *hub.Client
	if token == "" {
		return v, fmt.Errorf("bad token")
	}
	for _, v = range s.hub.GetUserList() {
		if v.Apikey == token {
			return v, nil
		}
	}
	return v, fmt.Errorf("bad token")
}
