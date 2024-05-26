package server

import (
	"arduinoteam/internal/engine"
	"arduinoteam/storage/sqlite"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	ErrUserExists = errors.New("user already exists")
)

type Instruction struct {
	client *Client
	room   *Room
}

func NewInstruction(client *Client, room *Room) Instruction {
	return Instruction{client: client, room: room}
}

type Hub struct {
	registerChan   chan Instruction
	unregisterChan chan Instruction
	broadcastChan  chan Message
	log            *slog.Logger
	storage        *sqlite.Storage
	rooms          map[string]*Room
	roomMutex      sync.RWMutex
	freeUsers      map[string]*Client
	usersMutex     sync.RWMutex
}

func (h *Hub) Run() {
	go func() {
		for {
			select {
			case instruction := <-h.registerChan:
				h.handleRegister(instruction.client, instruction.room)
			case instruction := <-h.unregisterChan:
				h.handleUnregister(instruction.client, instruction.room)
			case message := <-h.broadcastChan:
				h.handleBroadcast(message)
			}
		}
	}()
}
func (h *Hub) Register(client *Client, room *Room) {
	h.registerChan <- NewInstruction(client, room)
}

func (h *Hub) Unregister(client *Client, room *Room) {
	h.unregisterChan <- NewInstruction(client, room)
}

func (h *Hub) Broadcast(message Message) {
	h.broadcastChan <- message
}
func (h *Hub) handleRegister(client *Client, room *Room) {
	room.clients = append(room.clients, client)
	h.log.Debug("new user registered", "struct", fmt.Sprintf("%+v", room.clients))
	client.listen()
}
func (h *Hub) handleUnregister(client *Client, room *Room) {
	for i, c := range room.clients {
		if c == client {
			room.clients = append(room.clients[:i], room.clients[i+1:]...)
			break
		}
	}

	// client.close()
}

func (h *Hub) handleBroadcast(message Message) {
	// encoded := message.Encode()

	for _, client := range message.room.clients {
		// if client != message.author {
		client.write(message.payload)
		// }
	}
}

func (h *Hub) ListenClient(client *Client, room *Room) {
	for {
		messageType, msg, err := client.conn.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			h.Unregister(client, room)
			return
		}
		output, err := room.engine.Input(msg)
		if err != nil {
			h.log.Error("error get engine output")
		}
		h.Broadcast(Message{payload: output.Payload, room: room})
	}

}

func NewHub(storage *sqlite.Storage, log *slog.Logger) *Hub {
	return &Hub{
		registerChan:   make(chan Instruction),
		unregisterChan: make(chan Instruction),
		broadcastChan:  make(chan Message),
		log:            log,
		storage:        storage,
		rooms:          make(map[string]*Room),
		freeUsers:      make(map[string]*Client),
	}
}

func (h *Hub) CreateRoom(name string, engine engine.Engine) (*Room, error) {
	op := "server.hub.CreateRoom"
	var room *Room
	id := generateID()
	_, err := h.storage.SaveRoom(name, id)
	if err != nil {
		h.log.Error("failed to save user", slErr(err))

		return room, fmt.Errorf("%s: %w", op, err)
	}
	room = &Room{ID: id, Name: name, engine: engine}
	h.roomMutex.Lock()
	defer h.roomMutex.Unlock()
	h.rooms[id] = room
	return room, nil
}
func (h *Hub) CreateUser(login string, token string) (*Client, error) {
	op := "server.hub.CreateUser"
	// _, err := h.storage.SaveRoom(login, id)
	// if err != nil {
	// 	h.log.Error("failed to save user", slErr(err))

	// 	return room, fmt.Errorf("%s: %w", op, err)
	// }
	var client *Client
	h.usersMutex.Lock()
	defer h.usersMutex.Unlock()
	if h.freeUsers[login] != nil {
		return client, fmt.Errorf("%s: %w", op, ErrUserExists)
	}
	client = NewClient(login, token)
	h.freeUsers[login] = client
	return client, nil
}

func (h *Hub) GetRoom(id string) *Room {
	h.roomMutex.RLock()
	defer h.roomMutex.RUnlock()
	return h.rooms[id]
}
func (h *Hub) GetRoomList() map[string]*Room {
	h.roomMutex.RLock()
	defer h.roomMutex.RUnlock()
	return h.rooms
}
func (h *Hub) GetUser(token string) *Client {
	h.usersMutex.RLock()
	defer h.usersMutex.RUnlock()
	return h.freeUsers[token]
}

func generateID() string {
	return uuid.New().String()
}
