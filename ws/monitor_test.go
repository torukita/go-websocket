package ws

import (
	"testing"
	"time"
	"fmt"
)

func TestDo(t *testing.T) {
	monitor := NewMonitorCmd(1 * time.Second, defaultMonitorFunc())
	result, err := monitor.Do()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(result))

	monitor = NewMonitorCmd(1 * time.Second, hogeMonitorFunc())
	result, err = monitor.Do()
	if err != nil {
		t.Log(err)
	} else {
		t.Error("Should be failed")
	}

	monitor = NewMonitorCmd(1 * time.Second, nil)
	result, err = monitor.Do()
	if err != nil {
		t.Log(err)
	}
	t.Log(string(result))	
	
}

func hogeMonitorFunc() MonitorFunc {
	return func () []byte {
		str := fmt.Sprintf("[%v] sending data ...", time.Now())
		time.Sleep(2 * time.Second)
		return []byte(str)
	}
}

