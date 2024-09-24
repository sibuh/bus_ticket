package schedule

import (
	"time"
)

type Scheduler struct {
	smap map[string]chan string
}

func Init() *Scheduler {
	return &Scheduler{
		smap: make(map[string]chan string),
	}
}

func (s *Scheduler) append(id string, ch chan string) {
	s.smap[id] = ch
}

func (s *Scheduler) Get(id string) chan string {
	val, ok := s.smap[id]
	if !ok {
		return nil
	}
	return val
}

func (s *Scheduler) remove(id string) {
	delete(s.smap, id)
}

func (s *Scheduler) Schedule(id string, ch chan string, duration time.Duration, f func(id string) error) {
	s.append(id, ch)

	select {
	case <-ch:
		s.remove(id)
		return
	case <-time.After(duration * time.Second):
		s.remove(id)
		f(id)
	}
}
