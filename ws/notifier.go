package ws

import (
	"fmt"
	"time"
	log "github.com/Sirupsen/logrus"
)

type MonitorFunc func() []byte

type Notifier struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	message    chan []byte
	handler    MonitorFunc
	timer      time.Duration
}

var defaultMonitorHandler MonitorFunc = func() []byte {
	str := fmt.Sprintf("[%v] sending data ...", time.Now())
	log.Info(str)
	return []byte(str)	
}

func NewNotifier() *Notifier {
	return &Notifier {
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		message:    make(chan []byte),
		handler:    defaultMonitorHandler,
		timer:      1 * time.Second,
	}
}

func (n *Notifier) SetHandler(fn MonitorFunc) {
	n.handler = fn
}

func (n *Notifier) SetTimer(t time.Duration) {
	n.timer = t
}

func (n *Notifier) loopExec(fn MonitorFunc) {
	for {
		time.Sleep(n.timer)
		if len(n.clients) > 0 {
			n.Notify(fn())
		}
	}
}

func (n *Notifier) Run() {
	go func() {
		n.loopExec(n.handler)
	}()
	go n.run()
}

func (n *Notifier) Notify(msg []byte) {
	n.message <- msg
}

func (n *Notifier) run() {
	for {
		select {
		case message := <- n.message:
			for c, _ := range n.clients {
				c.Send(message)
			}
		case client := <- n.register:
			n.clients[client] = true
		case client := <- n.unregister:
			if _, ok := n.clients[client]; ok {
				delete(n.clients, client)
				client.Close()
			}
		}
	}
}

func (n *Notifier) Register(c *Client) {
	n.register <- c
}

func (n *Notifier) UnRegister(c *Client) {
	n.unregister <- c
}

