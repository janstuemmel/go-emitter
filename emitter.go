package emitter

import (
	"errors"
	"reflect"
)

// ErrEventNotRegistered is thrown if a event not exists
var ErrEventNotRegistered = errors.New("event not registered")

type (

	// Emitter interface
	Emitter interface {
		On(string, Listener)
		Once(string, Listener)
		Off(string, Listener) error
		Emit(string, interface{}) error
	}

	// Listener is a function with an empty interface value
	Listener func(interface{})

	// Emitter struct
	emitter struct {
		listeners map[string][]Listener
	}
)

// NewEmitter returns a new emitter struct
// as Emitter interface type
func NewEmitter() Emitter {

	// initialze new events map
	listeners := make(map[string][]Listener)

	return &emitter{
		listeners: listeners,
	}
}

// On registers a listener to a event
func (e emitter) On(event string, listener Listener) {
	e.listeners[event] = append(e.listeners[event], listener)
}

// Once registers a listener to a event
// but execute only once
func (e emitter) Once(event string, listener Listener) {

	var wrapper Listener

	wrapper = func(val interface{}) {

		// call function once
		listener(val)

		// unregister after first call
		e.Off(event, wrapper)
	}

	e.On(event, wrapper)
}

// Off unregisters a listener by event
func (e emitter) Off(event string, listener Listener) error {

	for key, events := range e.listeners {

		if key != event {
			continue
		}

		// iterate listerners for event
		for i, fn := range events {

			// compare pointers of listeners
			if reflect.ValueOf(fn).Pointer() != reflect.ValueOf(listener).Pointer() {
				continue
			}

			// remove listener
			events = append(events[:i], events[i+1:]...)

			// if no more listeners, delete event
			if len(events) < 1 {
				delete(e.listeners, key)
			}

			// success
			return nil
		}
	}

	return ErrEventNotRegistered
}

// Emit emits executes all listeners registered to a event
func (e emitter) Emit(event string, val interface{}) error {

	for key, events := range e.listeners {

		if key != event {
			continue
		}

		// execute each callback function
		for _, listener := range events {

			// execute listener
			listener(val)
		}

		// success
		return nil

	}

	return ErrEventNotRegistered
}
