# go-emitter  [![Build Status](https://travis-ci.org/janstuemmel/go-emitter.svg?branch=master)](https://travis-ci.org/janstuemmel/go-emitter)

A simple GO event emitter.

```go
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
```