package emitter_test

import (
	"fmt"
	"testing"

	"github.com/janstuemmel/go-emitter"
)

func Example() {

	var e emitter.Emitter
	var bar emitter.Listener

	// init emitter
	e = emitter.NewEmitter()

	// some function
	bar = func(val interface{}) {
		fmt.Printf("hello %s!\n", val)
	}

	// register listener to event
	e.On("foo", bar)

	// fire event
	e.Emit("foo", "world")

	// unregister
	err := e.Off("foo", bar)

	fmt.Println(err)

	// Output:
	// hello world!
	// <nil>
}

func Example_on() {

	// init emitter
	e := emitter.NewEmitter()

	// register listener
	e.On("foo", func(val interface{}) {
		fmt.Printf("hello %s!", val)
	})

	// emit
	e.Emit("foo", "world")

	// Output: hello world!
}

func Example_once() {

	// init emitter
	e := emitter.NewEmitter()

	// register listener
	e.Once("foo", func(val interface{}) {
		fmt.Printf("hello %s!\n", val)
	})

	// emit
	e.Emit("foo", "world")

	err := e.Emit("foo", "huhu")

	fmt.Println(err.Error())

	// Output:
	// hello world!
	// event not registered
}

func Example_off() {

	// init emitter
	e := emitter.NewEmitter()
	foo := func(val interface{}) {
		fmt.Printf("hello %s!\n", val)
	}

	// register listener
	e.On("foo", foo)

	// unregister listener
	e.Off("foo", foo)

	// emit
	err := e.Emit("foo", "huhu")

	fmt.Println(err.Error())

	// Output:
	// event not registered
}

func TestEmitter(t *testing.T) {

	noop := func(val interface{}) {}

	t.Run("nil listeners", func(t *testing.T) {

		e := emitter.NewEmitter()

		e.On("foo", nil)
		e.Once("bar", nil)

		errFoo := e.Emit("foo", nil)
		errBar := e.Emit("bar", nil)

		assert(t, "event not registered", errFoo.Error())
		assert(t, "event not registered", errBar.Error())
	})

	t.Run("call", func(t *testing.T) {

		// given
		e := emitter.NewEmitter()
		done := false
		e.On("foo", func(val interface{}) {
			done = true
		})

		// when
		e.Emit("foo", nil)

		// then
		assert(t, true, done)
	})

	t.Run("call event twice", func(t *testing.T) {

		// given
		e := emitter.NewEmitter()
		calledWith := []interface{}{}
		fun1 := func(val interface{}) {
			calledWith = append(calledWith, val)
		}
		fun2 := func(val interface{}) {
			calledWith = append(calledWith, val)
		}
		e.On("foo", fun1)
		e.On("foo", fun2)

		// when
		e.Emit("foo", 1)

		// then
		assert(t, 2, len(calledWith))
		assert(t, 1, calledWith[0])
		assert(t, 1, calledWith[1])
	})

	t.Run("call event not registered", func(t *testing.T) {

		// given
		e := emitter.NewEmitter()

		// some random events before
		e.On("foo", noop)
		e.On("foo", noop)
		e.On("bar", noop)

		// when
		err := e.Emit("404", nil)

		// then
		assert(t, "event not registered", err.Error())
	})

	t.Run("unregister event", func(t *testing.T) {

		// given
		e := emitter.NewEmitter()
		e.On("bar", noop) // register another listener
		e.On("foo", noop)

		// when
		err := e.Emit("foo", nil)

		// then
		assert(t, nil, err)

		// given
		e.Off("foo", noop)

		// when
		err = e.Emit("foo", nil)

		// then
		assert(t, "event not registered", err.Error())
	})

	t.Run("unregister event not exists", func(t *testing.T) {

		// given
		e := emitter.NewEmitter()

		// when
		err := e.Off("foo", noop)

		// then
		assert(t, "event not registered", err.Error())
	})

	t.Run("register event once", func(t *testing.T) {

		// given
		e := emitter.NewEmitter()
		called := 0
		callback := func(val interface{}) {
			called = called + 1
		}
		e.Once("foo", callback)

		// when
		e.Emit("foo", nil)

		// then
		assert(t, 1, called)

		// when
		err := e.Emit("foo", callback)

		// then
		assert(t, "event not registered", err.Error())
		assert(t, 1, called)
	})

	t.Run("unregister second listener", func(t *testing.T) {

		// given
		e := emitter.NewEmitter()
		fun := func(val interface{}) {}
		e.On("foo", noop)
		e.On("foo", fun)

		// when
		err := e.Off("foo", fun)

		// then
		assert(t, nil, err)
	})

	t.Run("duplicate listeners not allowed", func(t *testing.T) {

		// given
		e := emitter.NewEmitter()
		called := 0
		callback := func(val interface{}) {
			called = called + 1
		}

		// register duplicate event listeners
		e.Once("foo", callback)
		e.Once("foo", callback)

		// when
		e.Emit("foo", nil)

		// then
		assert(t, 1, called)
	})
}

func assert(t *testing.T, want interface{}, have interface{}) {

	// mark as test helper function
	t.Helper()

	if want != have {
		t.Errorf("Assertion failed for %s\n\twant:\t%v\n\thave:\t%v", t.Name(), want, have)
	}
}
