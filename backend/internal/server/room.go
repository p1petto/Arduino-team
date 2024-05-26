package server

import "arduinoteam/internal/engine"

type Room struct {
	ID      string        `json:"ID"`
	Name    string        `json:"name"`
	clients []*Client     `json:"-"`
	engine  engine.Engine `json:"-"`
}
