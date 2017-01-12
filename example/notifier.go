package main

import (
	"fmt"
	"time"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/torukita/go-websocket/ws"
	"github.com/gorilla/websocket"
	"github.com/Sirupsen/logrus"
	"encoding/json"
)

var (
	logConfig = middleware.LoggerConfig{
		Format: "${time_rfc3339} ${status} ${method} ${uri}\n",
	}
	upgrader = websocket.Upgrader{}
	log = logrus.New()	
)

func Example(n *ws.Notifier) echo.HandlerFunc {
	return func(w echo.Context) error {
		c, err := upgrader.Upgrade(w.Response(), w.Request(), nil)
		if err != nil {
			log.Error(err)
			return err
		}
		defer c.Close()

		client := ws.NewClient(c)
		n.Register(client)
		defer n.UnRegister(client)
		client.Start()
		log.Info("Closed WebSocket Client")
		return nil
	}
}

type Data struct {
	Time    string
	Message string
}

var monitorExec ws.MonitorFunc = func() []byte {
	d := Data{
		Time: fmt.Sprint(time.Now()),
		Message: "This is a example message",
	}
	v, _ := json.Marshal(d)
	return v
}

func Run(addr string, debug bool) {
	broker := ws.NewNotifier()
	broker.SetTimer(500 * time.Millisecond)
	broker.SetHandler(monitorExec)
	broker.Run()
	
	e := echo.New()
	e.Debug = debug
	e.Use(middleware.LoggerWithConfig(logConfig))
	e.Use(middleware.Recover())
	e.GET("/ws/example", Example(broker))
	e.Logger.Fatal(e.Start(addr))
}

func main() {
	addr := "localhost:8088"
	debug := true
	Run(addr, debug)
}

