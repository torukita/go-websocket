package ws

import (
	"github.com/gorilla/websocket"
	log "github.com/Sirupsen/logrus"
)

type Client struct {
	conn    *websocket.Conn
	message chan []byte
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		message: make(chan []byte),
	}
}

func (c *Client) Start() {
	go func () {
	LOOP:
		for {
			select {
			case msg, ok := <- c.message:
				if !ok {
					log.Warn("Server closed connection")
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
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Warnf("Client closed connection: %v", err)
			}
			break
		}
		// log.Infof("Recieved message: %s", string(message))
	}
}

func (c *Client) Send(msg []byte) {
	c.message <- msg
}

func (c *Client) Close() {
	close(c.message)
}

func (c *Client) Dump() {
	log.Infof("%v", c.conn.RemoteAddr())
}
