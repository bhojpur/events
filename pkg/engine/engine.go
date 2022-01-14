package engine

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	listenersMutex sync.RWMutex // protects listeners and interfaces
	listeners      = make(map[reflect.Type][]interface{})
	interfaces     = make([]reflect.Type, 0)
)

// BadListenerError is raised via panic() when AddListener is called with an
// invalid listener function.
type BadListenerError string

func (why BadListenerError) Error() string {
	return fmt.Sprintf("bad listener func: %s", string(why))
}

// AddListener registers a listener function that will be called when a matching
// event is dispatched. The type of the function's first (and only) argument
// declares the event type (or interface) to listen for.
func AddListener(fn interface{}) {
	listenersMutex.Lock()
	defer listenersMutex.Unlock()

	fnType := reflect.TypeOf(fn)

	// check that the function type is what we think: # of inputs/outputs, etc.
	// panic if conditions not met (because it's a programming error to have that happen)
	switch {
	case fnType.Kind() != reflect.Func:
		panic(BadListenerError("listener must be a function"))
	case fnType.NumIn() != 1:
		panic(BadListenerError("listener must take exactly one input argument"))
	}

	// the first input parameter is the event
	evType := fnType.In(0)

	// keep a list of listeners for each event type
	listeners[evType] = append(listeners[evType], fn)

	// if eventType is an interface, store it in a separate list
	// so we can check non-interface objects against all interfaces
	if evType.Kind() == reflect.Interface {
		interfaces = append(interfaces, evType)
	}
}

// Dispatch sends an event to all registered listeners that were declared
// to accept values of the event's type, or interfaces that the value implements.
func Dispatch(ev interface{}) {
	listenersMutex.RLock()
	defer listenersMutex.RUnlock()

	evType := reflect.TypeOf(ev)
	vals := []reflect.Value{reflect.ValueOf(ev)}

	// call listeners for the actual static type
	callListeners(evType, vals)

	// also check if the type implements any of the registered interfaces
	for _, in := range interfaces {
		if evType.Implements(in) {
			callListeners(in, vals)
		}
	}
}

func callListeners(t reflect.Type, vals []reflect.Value) {
	for _, fn := range listeners[t] {
		reflect.ValueOf(fn).Call(vals)
	}
}

// Updater is an interface that events can implement to combine updating and
// dispatching into one call.
type Updater interface {
	// Update is called by DispatchUpdate() before the event is dispatched.
	Update(update interface{})
}

// DispatchUpdate calls Update() on the event and then dispatches it. This is a
// shortcut for combining updates and dispatches into a single call.
func DispatchUpdate(ev Updater, update interface{}) {
	ev.Update(update)
	Dispatch(ev)
}
