package server

type Message struct {
	payload []byte
	room    *Room
}
