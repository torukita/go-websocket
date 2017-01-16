package ws

import (
	"github.com/gorilla/websocket"
	log "github.com/Sirupsen/logrus"
)

type Client struct {
	conn    *websocket.Conn
	message chan []byte
	notifier  *Notifier
}

func NewClient(conn *websocket.Conn, notifier *Notifier) *Client {
	return &Client{
		conn: conn,
		message: make(chan []byte),
		notifier: notifier,
	}
}

func (c *Client) Start() {
	if err := c.notifier.Register(c); err != nil {
		return
	}
	go func () {
	LOOP:
		for {
			select {
			case msg, ok := <- c.message:
				if !ok {
					c.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Error(err)
					break LOOP
				}
			}
		}
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Warnf("Client closed connection: %v", err)
			}
			break
		}
		// Never want to receive message from client
		log.Warnf("Recieved message from client: %s", string(msg))
		break
	}
}

func (c *Client) Send(msg []byte) {
	c.message <- msg
}

func (c *Client) Close() {
	c.notifier.UnRegister(c)
	close(c.message)
}

func (c *Client) Dump() {
	log.Infof("%v", c.conn.RemoteAddr())
}
