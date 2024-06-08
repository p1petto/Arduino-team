package hub

import (
	"sync"
	"time"

	"github.com/tjgq/ticker"

	"github.com/gorilla/websocket"
)

type Client struct {
	Login       string
	Apikey      string
	conn        *websocket.Conn
	messageChan chan []byte
	ticker      *ticker.Ticker
	done        chan bool
	Active      bool
	muActive    sync.Mutex
}

func (c *Client) StartTicker(t time.Duration) {
	c.ticker = ticker.New(t)
	c.ticker.Start()
	go func() {
		for {
			select {
			case <-c.done:
				return
			case <-c.ticker.C:
				c.setActive(true)
			}
		}
	}()
}

func (c *Client) setActive(b bool) {
	c.muActive.Lock()
	defer c.muActive.Unlock()
	c.Active = b
}

func (c *Client) isActive() bool {
	c.muActive.Lock()
	defer c.muActive.Unlock()
	return c.Active
}

func NewClient(login string, apikey string) *Client {
	return &Client{Login: login, Apikey: apikey, messageChan: make(chan []byte), done: make(chan bool), Active: true}
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
	c.done <- true
}

func (c *Client) write(message []byte) {
	c.messageChan <- message
}

func (c *Client) SetConnection(conn *websocket.Conn) {
	c.conn = conn
}
