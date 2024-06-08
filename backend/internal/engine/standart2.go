package engine

import (
	"encoding/json"
	"sync"
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
	Coords `mapstructure:",squash"`
	RGB    [3]uint8 `json:"color" mapstructure:"color"`
}
type Coords struct {
	X int
	Y int
}

func (e *StandartEngine) Input(input UserInput) ([][][3]uint8, error) {
	if !(input.Coords.Y >= 0 && input.Coords.Y < e.dy) || !(input.Coords.X >= 0 && input.Coords.X < e.dx) {
		return [][][3]uint8{}, ErrNotValidInput
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.GameMatrix[input.Coords.Y][input.Coords.X] = input.RGB
	return copySlice(e.GameMatrix), nil
}

func (e *StandartEngine) CurStateMessage() (Message, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	data, err := json.Marshal(e.GameMatrix)

	return Message{Payload: data}, err
}

func copySlice(original [][][3]uint8) [][][3]uint8 {
	// Создаем новый слайс с такой же длиной, как у оригинального
	newSlice := make([][][3]uint8, len(original))

	// Копируем каждый элемент оригинального слайса в новый слайс
	for i, innerSlice := range original {
		// Создаем новый внутренний слайс с такой же длиной, как у оригинального
		newInnerSlice := make([][3]uint8, len(innerSlice))
		copy(newInnerSlice, innerSlice)
		newSlice[i] = newInnerSlice
	}

	return newSlice
}
