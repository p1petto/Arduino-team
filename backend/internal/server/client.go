package server

import "github.com/gorilla/websocket"

type Client struct {
	login       string
	conn        *websocket.Conn
	apikey      string
	messageChan chan []byte
}

func NewClient(login string, apikey string) *Client {
	return &Client{login: login, apikey: apikey, messageChan: make(chan []byte)}
}

func (c *Client) listen() {
	go func() {
		for msg := range c.messageChan {
			c.conn.WriteMessage(websocket.TextMessage, msg)
		}
	}()
}

func (c *Client) write(message []byte) {
	c.messageChan <- message
}
