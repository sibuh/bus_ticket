package scheduler

import "time"

type Scheduler struct {
	smap map[string]chan string
}

func Init(smap map[string]chan string) *Scheduler {
	return &Scheduler{smap}
}

func (s *Scheduler) Append(id string, ch chan string) {
	s.smap[id] = ch
}

func (s *Scheduler) Get(id string) chan string {
	val, ok := s.smap[id]
	if !ok {
		return nil
	}
	return val
}

func (s *Scheduler) Remove(id string) {
	delete(s.smap, id)
}

func (s *Scheduler) Scheduler(id string, ch chan string, duration time.Duration, queryFunc func() error) {
	select {
	case <-ch:
		s.Remove(id)
		return
	case <-time.After(duration * time.Second):
		go queryFunc()
	}
}
