package ticket

import (
	"time"
)

func Scheduler(id string, ch chan int, duration time.Duration, queryFunc func() error) {
	time.Sleep(duration * time.Second)
	if err := queryFunc(); err != nil {
		return
	}
}
