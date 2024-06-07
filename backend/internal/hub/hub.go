package hub

import (
	"arduinoteam/internal/engine"
	"arduinoteam/internal/sl"
	"arduinoteam/storage/sqlite"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
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
	castChan       chan CastMessage
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
			case message := <-h.castChan:
				h.handleCast(message)
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
func (h *Hub) Cast(message CastMessage) {
	h.castChan <- message
}
func (h *Hub) handleRegister(client *Client, room *Room) {
	room.clients = append(room.clients, client)
	h.log.Debug("new user registered", "op", "handleRegister", "struct", fmt.Sprintf("%+v", room.clients))
	client.listen()
}
func (h *Hub) handleUnregister(client *Client, room *Room) {
	for i, c := range room.clients {
		if c == client {
			room.clients = append(room.clients[:i], room.clients[i+1:]...)
			break
		}
	}
	client.close()
}

func (h *Hub) handleBroadcast(message Message) {
	// encoded := message.Encode()

	for _, client := range message.room.clients {
		// if client != message.author {
		client.write(message.payload)
		// }
	}
}
func (h *Hub) handleCast(message CastMessage) {
	// encoded := message.Encode()
	for _, client := range message.room.clients {
		if client == message.Client {
			client.write(message.payload)
		}
	}
}

func (h *Hub) ListenClient(client *Client, room *Room) {
	for {
		messageType, msg, err := client.conn.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			h.Unregister(client, room)
			return
		}

		var payload map[string]interface{}
		// fmt.Printf("InputUser: %s", string(payload))
		err = json.Unmarshal(msg, &payload)
		if err != nil {
			fmt.Printf("%+v", err)
		}
		payloadType, ok := payload["type"]
		if !ok {
			h.log.Error("Error getting payload type", sl.Err(err))
			continue
		}
		h.log.Debug("get payload", "struct", payload)
		switch payloadType {
		case "Input":
			var input engine.UserInput
			err := mapstructure.Decode(payload, &input)
			if err != nil {
				h.log.Error("Error decode UserInput", sl.Err(err))
				continue
			}
			response, err := room.engine.Input(input)
			if err != nil {
				h.log.Error("Error getting engine responce", sl.Err(err))
				continue
			}
			go func() {
				if room.Status == "Connected" {
					room.esp_chan <- fmt.Sprintf("%d|%d|%d|%d|%d|", input.Coords.X, input.Coords.Y, input.RGB[0], input.RGB[1], input.RGB[2])
				} else {
					h.RaiseWSError("ESP server is down", client)
				}
			}()
			data, err := json.Marshal(map[string]interface{}{"type": "Output", "message": response})
			if err != nil {
				h.log.Error("Error marshaling in ListenClient", sl.Err(err))
				continue
			}
			h.Broadcast(Message{payload: data, room: room})
		}

	}

}

func (h *Hub) RaiseWSError(message string, client *Client) {
	data, err := json.Marshal(map[string]interface{}{"type": "Error", "message": message})
	if err != nil {
		h.log.Error("Failed to error raise", sl.Err(err))
		return
	}
	client.write(data)
}

func NewHub(storage *sqlite.Storage, log *slog.Logger) *Hub {
	return &Hub{
		registerChan:   make(chan Instruction),
		unregisterChan: make(chan Instruction),
		broadcastChan:  make(chan Message),
		castChan:       make(chan CastMessage),
		log:            log,
		storage:        storage,
		rooms:          make(map[string]*Room),
		freeUsers:      make(map[string]*Client),
	}
}

func (h *Hub) CreateRoom(name string, esp_ip string) (*Room, error) {
	op := "server.hub.CreateRoom"
	var room *Room
	id := generateID()
	_, err := h.storage.SaveRoom(name, id)
	if err != nil {
		h.log.Error("failed to save user", sl.Err(err))

		return room, fmt.Errorf("%s: %w", op, err)
	}
	standartEngn := engine.NewStandartEngine(24, 12)
	room = &Room{ID: id, Name: name, Ip: esp_ip, Status: "Pending", engine: standartEngn, esp_chan: make(chan string)}
	go room.Run()
	// room.engine.Run()

	h.roomMutex.Lock()
	defer h.roomMutex.Unlock()
	h.rooms[id] = room
	return room, nil
}
func (h *Hub) CreateUser(login string, token string) (*Client, error) {
	op := "server.hub.CreateUser"
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
	copy := make(map[string]*Room)
	for key, value := range h.rooms {
		copy[key] = value
	}
	return copy
}
func (h *Hub) GetUserList() map[string]*Client {
	h.usersMutex.RLock()
	defer h.usersMutex.RUnlock()
	copy := make(map[string]*Client)
	for key, value := range h.freeUsers {
		copy[key] = value
	}
	return copy
}
func (h *Hub) GetUser(token string) *Client {
	h.usersMutex.RLock()
	defer h.usersMutex.RUnlock()
	return h.freeUsers[token]
}

func generateID() string {
	return uuid.New().String()
}
