package event

import "sync"

type Event interface {
	Signal()
	Wait()
}

type event struct {
	cond *sync.Cond
}

func new() *event {
	return &event{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func New() Event {
	return new()
}

func (e *event) Signal() {
	e.cond.Signal()
}

func (e *event) Wait() {
	e.cond.L.Lock()
	e.cond.Wait()
	e.cond.L.Unlock()
}
