package ws

import (
	"errors"
	"time"
)

var (
	maxClients = 10
)

type notifier struct {
	name       string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	message    chan []byte
	cmd        MonitorCmd
}

type Notifier interface {
	Run()
	Notify([]byte)
	Regist(*Client) error
	UnRegist(*Client)
}

func NewNotifier(name string, fn MonitorCmd) Notifier {
	return &notifier {
		name:       name,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		message:    make(chan []byte),
		cmd: fn,
	}
}

func (n *notifier) doLoopExec() {
	t := time.NewTicker(n.cmd.GetInterval())
	for {
		select {
		case <- t.C:
			if len(n.clients) > 0 {
				result, err := n.cmd.Do()
				if err != nil {
					logger.Errorf("notifier(%s): %s", n.name, err)
					continue
				}
				n.Notify(result)
			}
		}
	}
	t.Stop()
}

func (n *notifier) Run() {
	go func() {
		n.doLoopExec()
	}()
	go n.run()
}

func (n *notifier) Notify(msg []byte) {
	n.message <- msg
}

func (n *notifier) run() {
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
			}
		}
	}
}

func (n *notifier) Regist(c *Client) error {
	if len(n.clients) >= maxClients {
		return errors.New("notifer: reached max registration")
    }
	n.register <- c
	return nil
}

func (n *notifier) UnRegist(c *Client) {
	n.unregister <- c
}

