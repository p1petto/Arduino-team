package engine

import (
	"encoding/json"
	"fmt"
)

type Engine interface {
	CurStateMessage() (Message, error)
	Input([]byte) (Message, error)
}

type Color struct {
	RGBL [4]int
}
type Coords struct {
	X, Y int
}

// func (i *Input) UnmarshalJSON(b []byte) error {

// }

// func (i Input) MarshalJSON() ([]byte, error) {
// 	return []byte(p), nil
// }

type StandartEngine struct {
	GameMatrix [][]Color `json:"matrix"`
}

type Input struct {
	Coords Coords `json:"coords"`
	Color  Color  `json:"color"`
}

func NewStandartEngine(dx int, dy int) StandartEngine {
	a := make([][]Color, dy)
	for i := range a {
		a[i] = make([]Color, dx)
	}
	return StandartEngine{GameMatrix: a}
}

func (s *StandartEngine) Input(payload []byte) (Message, error) {
	var input Input
	err := json.Unmarshal(payload, &input)
	if err != nil {
		fmt.Printf("%+v", err)
	}
	fmt.Printf("%+v", input)
	s.GameMatrix[input.Coords.Y][input.Coords.X] = input.Color

	data, err := json.Marshal(s.GameMatrix)

	return Message{Payload: data}, err

}

func (s StandartEngine) CurStateMessage() (Message, error) {
	data, err := json.Marshal(s.GameMatrix)

	return Message{Payload: data}, err
}
