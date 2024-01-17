package pubsub

import (
	"context"
	"errors"
	"github.com/yudhasubki/eventpool"
	"sync"
)

type EventPoolFactory struct {
	locker   *sync.RWMutex
	eventMap map[string][]*eventpool.Eventpool
}

func NewEventPoolFactory() *EventPoolFactory {
	return &EventPoolFactory{
		locker:   new(sync.RWMutex),
		eventMap: make(map[string][]*eventpool.Eventpool),
	}
}

func (f *EventPoolFactory) Subscribe(topic string, listener eventpool.EventpoolListener) error {
	event := eventpool.New()
	event.Submit(listener)
	event.Run()
	f.locker.Lock()
	events := f.eventMap[topic]
	if len(events) == 0 {
		f.eventMap[topic] = []*eventpool.Eventpool{event}
	} else {
		events = append(events, event)
		f.eventMap[topic] = events
	}

	f.locker.Unlock()
	return nil
}

func (f *EventPoolFactory) Publish(topic string, payload interface{}) error {
	events := f.eventMap[topic]
	if len(events) == 0 {
		return errors.New("not found event")
	}

	for _, event := range events {
		event.Publish(eventpool.SendJson(payload))
	}

	return nil
}

func (f *EventPoolFactory) Stop(ctx context.Context) error {
	var wg sync.WaitGroup
	for name, events := range f.eventMap {
		for _, q := range events {
			wg.Add(1)
			go func(name string, q *eventpool.Eventpool) {
				defer wg.Done()
				q.Close()
			}(name, q)
		}
	}
	wg.Wait()
	return nil
}
