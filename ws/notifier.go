package ws

import (
	"errors"
	"fmt"
	"time"
	log "github.com/Sirupsen/logrus"
	"context"
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

func defaultMonitorHandler() MonitorFunc {
	return func() []byte {
		str := fmt.Sprintf("[%v] sending data ...", time.Now())
		log.Info(str)
		return []byte(str)
	}
}

func NewNotifier() *Notifier {
	return &Notifier {
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		message:    make(chan []byte),
		handler:    defaultMonitorHandler(),
		timer:      1 * time.Second,
	}
}

func (n *Notifier) SetHandler(t time.Duration, fn MonitorFunc) {
	n.timer = t
	n.handler = fn
}

func (n *Notifier) doExecHandler() []byte {
	return n.handler()
}

func (n *Notifier) doLoopExec(ctx context.Context) {
	for {
		time.Sleep(n.timer)
		if len(n.clients) > 0 {
			recvCh := make(chan []byte, 1)
			ctx, cancel := context.WithTimeout(ctx, n.timer /2)
			defer func() {
				log.Error("notifier: canceled")
				cancel()
			}()
			go func() {
				recvCh <- n.doExecHandler()
			}()
			select {
			case <-ctx.Done():
				log.Error("notifier:", ctx.Err())
			case recv := <- recvCh:
				n.Notify(recv)
			}
		}
	}
}

func (n *Notifier) Run(ctx context.Context) {
	go func() {
		n.doLoopExec(ctx)
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
			}
		}
	}
}

func (n *Notifier) Register(c *Client) error {
	if len(n.clients) > 2 {
		return errors.New("hoge")
    }
	n.register <- c
	return nil
}

func (n *Notifier) UnRegister(c *Client) {
	n.unregister <- c
}

