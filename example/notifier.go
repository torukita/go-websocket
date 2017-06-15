package main

import (
	"fmt"
	"time"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/torukita/go-websocket/ws"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"encoding/json"
)

var (
	logConfig = middleware.LoggerConfig{
		Format: "${time_rfc3339} ${status} ${method} ${uri} ${remote_ip} ${latency_human}\n",
	}
	upgrader = websocket.Upgrader{}
	log = logrus.New()
	interval = time.Duration(1 * time.Second)
)

func Example(n ws.Notifier) echo.HandlerFunc {
	return func(w echo.Context) error {
		c, err := upgrader.Upgrade(w.Response(), w.Request(), nil)
		if err != nil {
			return err
		}
		defer c.Close()

		client := ws.NewClient(c, n)
		client.Start()
		defer client.Close()
		return nil
	}
}

type Data struct {
	Time    string
	Message string
}

func monitorExec() ws.MonitorFunc {
	return func() []byte {
//		time.Sleep(2 * time.Second) 
		d := Data{
			Time: fmt.Sprint(time.Now()),
			Message: "This is a example message",
		}
		v, _ := json.Marshal(d)
		return v
	}
}

func Run(addr string, debug bool) {
	cmd := ws.NewMonitorCmd(1 * time.Second, monitorExec())
	broker := ws.NewNotifier("hoge", cmd)
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

