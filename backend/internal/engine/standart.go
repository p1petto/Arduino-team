package engine

// import (
// 	"encoding/json"
// 	"fmt"
// )

// type Coords struct {
// 	X int
// 	Y int
// }

// type UserInput struct {
// 	Coords
// 	RGB [3]uint8 `json:"color"`
// }

// func NewStandartEngine(dx int, dy int) *StandartEngine {
// 	a := make([][][3]uint8, dy)
// 	for i := range a {
// 		a[i] = make([][3]uint8, dx)
// 	}
// 	return &StandartEngine{GameMatrix: a, inputChan: make(chan UserInput), dx: dx, dy: dy}
// }

// func (s *StandartEngine) Input(payload []byte) (Message, error) {
// 	var input UserInput
// 	// fmt.Printf("InputUser: %s", string(payload))
// 	err := json.Unmarshal(payload, &input)
// 	if err != nil {
// 		fmt.Printf("%+v", err)
// 		return Message{}, fmt.Errorf("%+w", err)
// 	}
// 	// fmt.Printf("%+v", input)
// 	fmt.Println("input send to the channel")
// 	s.inputChan <- input
// 	data, err := json.Marshal(s.GameMatrix)
// 	fmt.Println("message marshaled")
// 	return Message{Payload: data}, err

// }

// func (s *StandartEngine) Run() {
// 	go func() {
// 		for input := range s.inputChan {
// 			if input.Coords.Y >= 0 && input.Coords.Y < s.dy {
// 				if input.Coords.X >= 0 && input.Coords.X < s.dx {
// 					s.GameMatrix[input.Coords.Y][input.Coords.X] = input.RGB
// 					fmt.Println("input processed")
// 				}
// 			}

// 		}
// 		fmt.Println("input stopped")
// 	}()
// 	fmt.Println("goroutine started")
// }

// func (s *StandartEngine) CurStateMessage() (Message, error) {
// 	data, err := json.Marshal(s.GameMatrix)

// 	return Message{Payload: data}, err
// }
