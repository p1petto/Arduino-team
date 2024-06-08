package hub

import (
	"arduinoteam/internal/engine"
	"arduinoteam/internal/sl"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

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
				if errors.Is(err, engine.ErrNotValidInput) {
					h.RaiseWSError(err.Error(), client)
					continue
				}
				h.log.Error("Error getting engine responce", sl.Err(err))
				continue
			}
			go func() {
				if room.Status == "Connected" {
					for i := 0; i < 3; i++ {
						room.esp_chan <- fmt.Sprintf("%d|%d|%d|%d|%d|", input.Coords.X, input.Coords.Y, input.RGB[0], input.RGB[1], input.RGB[2])
					}
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
