package ticket

import (
	"time"
)

func Scheduler(id string, ch chan string, duration time.Duration, queryFunc func() error) {
	// var response interface{}, ch chan boolean
	select {
	case <-ch:
		return
	case <-time.After(duration * time.Second):
		// If query fails retry or log
		go queryFunc()

		// res := <- ch

		// if(res.status !== 200) go queryFunc(response, )
	}
}
