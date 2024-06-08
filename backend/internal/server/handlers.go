package server

import (
	"arduinoteam/internal/hub"
	"arduinoteam/internal/sl"
	"arduinoteam/storage"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleRoomCreate(w http.ResponseWriter, r *http.Request) {
	s.setCORSPolicy(w)
	// name := chi.URLParam(r, "name")
	// if name == "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	r.ParseForm()
	var name string
	if name = r.Form.Get("name"); name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var esp_ip string
	if esp_ip = r.Form.Get("IP"); esp_ip == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	room, err := s.hub.CreateRoom(name, esp_ip)
	if err != nil {
		if errors.Is(err, storage.ErrRoomExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.log.Info("room created", "id", room.ID, "name", room.Name)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(room.ID))
}
func (s *Server) handleApiKeyCreate(w http.ResponseWriter, r *http.Request) {
	s.setCORSPolicy(w)
	login := chi.URLParam(r, "login")
	if login == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := GenerateRandomStringURLSafe(32)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	// NewClient(login, token)
	_, err = s.hub.CreateUser(login, token)
	if err != nil {
		if errors.Is(err, hub.ErrUserExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
	s.log.Info("client connected", "name", login)
	// s.log.Debug("clients slice", "struct", fmt.Sprintf("%+v", s.hub.freeUsers))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(token))
}

func (s *Server) handleRoomGet(w http.ResponseWriter, r *http.Request) {
	s.setCORSPolicy(w)
	roomID := chi.URLParam(r, "room")
	if roomID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room := s.hub.GetRoom(roomID)
	if room == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	s.log.Debug("room requested", "name", room.ID)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleRoomListGet(w http.ResponseWriter, r *http.Request) {
	s.setCORSPolicy(w)
	room := s.hub.GetRoomList()
	s.log.Debug("room list requested")

	encode(w, r, 200, room)
}

func (s *Server) handleOptions(w http.ResponseWriter, r *http.Request) {
	s.setCORSPolicy(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	op := "handlers.handleWS"
	roomID := chi.URLParam(r, "room")
	if roomID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room := s.hub.GetRoom(roomID)
	if room == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	reqToken := r.URL.Query().Get("token")
	if reqToken == "" {
		http.Error(w, "Incorrect token", http.StatusBadRequest)
		return
	}
	client, err := s.CheckToken(reqToken)
	if err != nil {
		s.log.Error(op+"bad token", sl.Err(err))
		http.Error(w, "Bad Token", http.StatusUnauthorized)
		return
	}
	s.log.Debug("new user connected", "op", op, "login", client.Login)
	w.Header().Set("Sec-Websocket-Protocol", reqToken)
	s.upgrader.Subprotocols = []string{reqToken}
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Error("Failed to upgrade connection", sl.Err(err))
		return
	}

	client.SetConnection(conn)
	client.StartTicker(room.TickerDuration)
	s.hub.Register(client, room)

	s.hub.ListenClient(client, room)
}
