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
		//  определение типа сообщения
		// fmt.Printf("InputUser: %s", string(payload))
		payload, payloadType, err := getMessageType(msg)
		if err != nil {
			h.log.Debug("Fail get messageType", sl.Err(err))
			h.RaiseWSError(err.Error(), client)
			continue
		}
		h.log.Debug("get payload", "struct", payload)

		switch payloadType {
		case "Input":
			h.InputGameLogic(payload, room, client)
		}

	}

}

func getMessageType(msg []byte) (map[string]interface{}, string, error) {

	var payload map[string]interface{}

	err := json.Unmarshal(msg, &payload)
	if err != nil {
		return nil, "", errors.New("fail to unmarshal struct")
	}
	payloadType, ok := payload["type"].(string)
	if !ok {
		return nil, "", errors.New("поле type не указано")
	}
	return payload, payloadType, nil
}

func (h *Hub) InputGameLogic(payload map[string]interface{}, room *Room, client *Client) {

	var input engine.UserInput
	err := mapstructure.Decode(payload, &input)
	if err != nil {
		h.log.Error("Error decode UserInput", sl.Err(err))
		return
	}
	response, err := room.engine.Input(input)
	if err != nil {
		if errors.Is(err, engine.ErrNotValidInput) {
			h.RaiseWSError(err.Error(), client)
			return
		}
		h.log.Error("Error getting engine responce", sl.Err(err))
		return
	}
	if !client.isActive() {
		h.log.Debug("User cannot input yet")
		h.RaiseWSError("Вы ещё не можете отправить новый ввод", client)
		return
	}
	client.setActive(false)
	client.ticker.Start()
	go func() {
		if room.Status == "Connected" {
			for i := 0; i < 3; i++ {
				room.esp_chan <- fmt.Sprintf("%d|%d|%d|%d|%d|", input.Coords.X, input.Coords.Y, input.RGB[0], input.RGB[1], input.RGB[2])
			}
		} else {
			h.RaiseWSError("ESP server is down", client)
		}
	}()
	data, err := json.Marshal(map[string]interface{}{"type": "Output", "message": response, "time": room.TickerDuration.Minutes()})
	if err != nil {
		h.log.Error("Error marshaling in ListenClient", sl.Err(err))
		return
	}

	h.Broadcast(Message{payload: data, room: room})
}

func (h *Hub) RaiseWSError(message string, client *Client) {
	data, err := json.Marshal(map[string]interface{}{"type": "Error", "message": message})
	if err != nil {
		h.log.Error("Failed to error raise", sl.Err(err))
		return
	}
	client.write(data)
}
