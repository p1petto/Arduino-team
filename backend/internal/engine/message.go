package engine

import "errors"

type Message struct {
	Payload []byte
}

type EngineOutput struct{}
type EngineInput struct {
}

var (
	ErrNotValidInput = errors.New("неправильный ввод от пользователя")
)
