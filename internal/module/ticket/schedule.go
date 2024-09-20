package ticket

import (
	"time"
)

func Scheduler(id string, ch chan string, duration time.Duration, queryFunc func() error) {
	// var response interface{}, ch chan boolean
	select {
	case <-ch:
		// when the scheduler exites,
		return
	case <-time.After(duration * time.Second):
		go queryFunc()
	}
}
