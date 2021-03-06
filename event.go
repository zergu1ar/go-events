package go_events

import "context"

type (
	Event struct {
		name string
		ctx  context.Context
		in   chan interface{}
		out  []chan interface{}
	}
)

func NewEvent(name string, ctx context.Context) *Event {
	event := Event{
		name: name,
		ctx:  ctx,
		in:   make(chan interface{}, 100),
		out:  make([]chan interface{}, 0),
	}
	go event.iterate()
	return &event
}

func (e *Event) iterate() {
	for {
		select {
		case <-e.ctx.Done():
			e.close()
			return
		case data := <-e.in:
			if len(e.out) > 0 {
				for _, c := range e.out {
					c <- data
				}
			}
		}
	}
}

func (e *Event) Dispatch(data ...interface{}) {
	for _, d := range data {
		e.in <- d
	}
}

func (e *Event) Listen() <-chan interface{} {
	c := make(chan interface{}, 100)
	e.out = append(e.out, c)
	return c
}

func (e *Event) close() {
	for _, c := range e.out {
		close(c)
	}
}
