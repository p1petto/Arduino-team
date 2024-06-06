package engine

import (
	"encoding/json"
	"errors"
	"sync"
)

var (
	ErrNotValidInput = errors.New("неправильный ввод от пользователя")
)

type StandartEngine struct {
	mu         sync.RWMutex
	GameMatrix [][][3]uint8 `json:"matrix"`
	dx         int
	dy         int
}

func NewStandartEngine(dx int, dy int) *StandartEngine {
	a := make([][][3]uint8, dy)
	for i := range a {
		a[i] = make([][3]uint8, dx)
	}
	return &StandartEngine{GameMatrix: a, dx: dx, dy: dy}
}

type UserInput struct {
	Coords
	RGB [3]uint8 `json:"color"`
}
type Coords struct {
	X int
	Y int
}

func (e *StandartEngine) Input(input UserInput) ([][][3]uint8, error) {
	// var input UserInput
	// // fmt.Printf("InputUser: %s", string(payload))
	// err := json.Unmarshal(payload, &input)
	// if err != nil {
	// 	fmt.Printf("%+v", err)
	// 	// return Message{}, fmt.Errorf("%+w", err)
	// }

	if !(input.Coords.Y >= 0 && input.Coords.Y < e.dy) {
		if !(input.Coords.X >= 0 && input.Coords.X < e.dx) {
			return [][][3]uint8{}, ErrNotValidInput
		}
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.GameMatrix[input.Coords.Y][input.Coords.X] = input.RGB
	// data, err := json.Marshal(e.GameMatrix)
	return e.GameMatrix, nil
}

func (e *StandartEngine) CurStateMessage() (Message, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	data, err := json.Marshal(e.GameMatrix)

	return Message{Payload: data}, err
}
