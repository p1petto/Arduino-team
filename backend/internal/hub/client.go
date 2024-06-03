package hub

import "github.com/gorilla/websocket"

type Client struct {
	Login       string
	Apikey      string
	conn        *websocket.Conn
	messageChan chan []byte
}

func NewClient(login string, apikey string) *Client {
	return &Client{Login: login, Apikey: apikey, messageChan: make(chan []byte)}
}

func (c *Client) listen() {
	go func() {
		for msg := range c.messageChan {
			c.conn.WriteMessage(websocket.TextMessage, msg)
		}
	}()
}

func (c *Client) close() {
	close(c.messageChan)
	c.conn.Close()
}

func (c *Client) write(message []byte) {
	c.messageChan <- message
}

func (c *Client) SetConnection(conn *websocket.Conn) {
	c.conn = conn
}
