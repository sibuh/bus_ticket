package scheduler

type Scheduler struct {
	Map map[string]chan string
}

func (s *Scheduler) Append(id string, ch chan string) {
	s.Map[id] = ch
}

func (s *Scheduler) Remove(id string) {
	delete(s.Map, id)
}
