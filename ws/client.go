package ws

import (
	"github.com/gorilla/websocket"
)


type Client struct {
	conn    *websocket.Conn
	message chan []byte
	broker  Notifier
}

func NewClient(conn *websocket.Conn, broker Notifier) *Client {
	return &Client{
		conn: conn,
		message: make(chan []byte),
		broker: broker,
	}
}


func (c *Client) Start() {
	if err := c.broker.Regist(c); err != nil {
		logger.Error(err)
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
					logger.Error(err)
					break LOOP
				}
			}
		}
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logger.Warnf("Client closed connection: %v", err)
			}
			break
		}
		// Never want to receive message from client
		logger.Warnf("Recieved message from client: %s", string(msg))
		break
	}
}

func (c *Client) Send(msg []byte) {
	c.message <- msg
}

func (c *Client) Close() {
	c.broker.UnRegist(c)
	close(c.message)
}

func (c *Client) Dump() {
	logger.Infof("%v", c.conn.RemoteAddr())
}
