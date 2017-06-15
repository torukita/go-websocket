package ws

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	logger = logrus.New()
)

type MonitorFunc func() []byte

func defaultMonitorFunc() MonitorFunc {
	return func() []byte {
		str := fmt.Sprintf("[%v] sending data ...", time.Now())
		return []byte(str)
	}
}

type MonitorCmd interface {
	GetInterval() (time.Duration)
	Do() ([]byte, error)
}

func NewMonitorCmd(interval time.Duration, fn MonitorFunc) MonitorCmd {
	if fn == nil {
		fn = defaultMonitorFunc()
	}
	cmd := &monitor{
		interval: interval,
		handler:  fn,
	}
	return cmd
}

type monitor struct {
	interval time.Duration
	handler  MonitorFunc
}

func (m *monitor) GetInterval() time.Duration {
	return m.interval
}

func (m *monitor) Do() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.interval)
	defer cancel()
	recvCh := make(chan []byte, 1)
	go func() {
		recvCh <- m.handler()
	}()
	select {
	case <-ctx.Done():
		logger.Debug(ctx.Err())
		return nil, ctx.Err()
	case result := <-recvCh:
		return result, nil
	}
}
