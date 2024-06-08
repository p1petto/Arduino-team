package hub

import (
	"arduinoteam/internal/engine"
	"fmt"
	"net"
	"time"
)

type Engine interface {
	CurStateMessage() (engine.Message, error)
	Input(engine.UserInput) ([][][3]uint8, error)
	// Run()
}

type Room struct {
	ID             string    `json:"ID"`
	Name           string    `json:"name"`
	clients        []*Client `json:"-"`
	engine         Engine    `json:"-"`
	Ip             string    `json:"IP"`
	Status         string    `json:"status"`
	esp_chan       chan string
	TickerDuration time.Duration
}

func (r *Room) Run() {
	fmt.Println("runing: " + r.ID)
	conn, err := net.Dial("tcp", r.Ip+":80")
	if err != nil {
		r.Status = "Failed"
		return
	}
	r.Status = "Connected"
	go func() {
		for payload := range r.esp_chan {
			_, err = fmt.Fprintf(conn, "%s", payload)
			if err != nil {
				r.Status = "Failed to write"
				return
			}
			// fmt.Println("Successfuly send data to esp")
		}
	}()

}
