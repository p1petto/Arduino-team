package hub

import "arduinoteam/internal/engine"

type Engine interface {
	CurStateMessage() (engine.Message, error)
	Input([]byte) ([]byte, error)
	// Run()
}

type Room struct {
	ID      string    `json:"ID"`
	Name    string    `json:"name"`
	clients []*Client `json:"-"`
	engine  Engine    `json:"-"`
}
