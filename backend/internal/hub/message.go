package hub

type Message struct {
	payload []byte
	room    *Room
}

type CastMessage struct {
	payload []byte
	room    *Room
	Client  *Client
}
