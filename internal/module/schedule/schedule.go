package schedule

import (
	"time"

	"github.com/google/uuid"
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

func (s *Scheduler) Schedule(id uuid.UUID, ch chan string, duration time.Duration, f func(id uuid.UUID) error) {
	s.append(id.String(), ch)

	select {
	case <-ch:
		s.remove(id.String())
		return
	case <-time.After(duration * time.Second):
		s.remove(id.String())
		f(id)
	}
}
