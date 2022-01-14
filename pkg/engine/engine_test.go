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
	"reflect"
	"testing"
	"time"
)

type testInterface1 interface {
	TestFunc1()
}

type testInterface2 interface {
	TestFunc2()
}

type testEvent1 struct {
}

type testEvent2 struct {
	triggered bool
}

func (testEvent1) TestFunc1() {}

func (*testEvent2) TestFunc2() {}

func clearListeners() {
	listenersMutex.Lock()
	defer listenersMutex.Unlock()

	listeners = make(map[reflect.Type][]interface{})
	interfaces = make([]reflect.Type, 0)
}

func TestStaticListener(t *testing.T) {
	clearListeners()

	triggered := false
	AddListener(func(testEvent1) { triggered = true })
	AddListener(func(testEvent2) { t.Errorf("wrong listener type triggered") })
	Dispatch(testEvent1{})

	if !triggered {
		t.Errorf("static listener failed to trigger")
	}
}

func TestPointerListener(t *testing.T) {
	clearListeners()

	testEvent := new(testEvent2)
	AddListener(func(ev *testEvent2) { ev.triggered = true })
	AddListener(func(testEvent2) { t.Errorf("non-pointer listener triggered on pointer type") })
	Dispatch(testEvent)

	if !testEvent.triggered {
		t.Errorf("pointer listener failed to trigger")
	}
}

func TestInterfaceListener(t *testing.T) {
	clearListeners()

	triggered := false
	AddListener(func(testInterface1) { triggered = true })
	AddListener(func(testInterface2) { t.Errorf("interface listener triggered on non-matching type") })
	Dispatch(testEvent1{})

	if !triggered {
		t.Errorf("interface listener failed to trigger")
	}
}

func TestEmptyInterfaceListener(t *testing.T) {
	clearListeners()

	triggered := false
	AddListener(func(interface{}) { triggered = true })
	Dispatch("this should match interface{}")

	if !triggered {
		t.Errorf("interface{} listener failed to trigger")
	}
}

func TestMultipleListeners(t *testing.T) {
	clearListeners()

	triggered1, triggered2 := false, false
	AddListener(func(testEvent1) { triggered1 = true })
	AddListener(func(testEvent1) { triggered2 = true })
	Dispatch(testEvent1{})

	if !triggered1 || !triggered2 {
		t.Errorf("not all matching listeners triggered")
	}
}

func TestBadListenerWrongInputs(t *testing.T) {
	clearListeners()

	defer func() {
		err := recover()

		if err == nil {
			t.Errorf("bad listener func (wrong # of inputs) failed to trigger panic")
			return
		}

		blErr, ok := err.(BadListenerError)
		if !ok {
			panic(err) // this is not the error we were looking for; re-panic
		}

		want := "bad listener func: listener must take exactly one input argument"
		if got := blErr.Error(); got != want {
			t.Errorf(`BadListenerError.Error() = "%s", want "%s"`, got, want)
		}
	}()

	AddListener(func() {})
	Dispatch(testEvent1{})
}

func TestBadListenerWrongType(t *testing.T) {
	clearListeners()

	defer func() {
		err := recover()

		if err == nil {
			t.Errorf("bad listener type (not a func) failed to trigger panic")
		}

		blErr, ok := err.(BadListenerError)
		if !ok {
			panic(err) // this is not the error we were looking for; re-panic
		}

		want := "bad listener func: listener must be a function"
		if got := blErr.Error(); got != want {
			t.Errorf(`BadListenerError.Error() = "%s", want "%s"`, got, want)
		}
	}()

	AddListener("this is not a function")
	Dispatch(testEvent1{})
}

func TestAsynchronousDispatch(t *testing.T) {
	clearListeners()

	triggered := make(chan bool)
	AddListener(func(testEvent1) { triggered <- true })
	go Dispatch(testEvent1{})

	select {
	case <-triggered:
	case <-time.After(time.Second):
		t.Errorf("asynchronous dispatch failed to trigger listener")
	}
}

func TestDispatchPointerToValueInterfaceListener(t *testing.T) {
	clearListeners()

	triggered := false
	AddListener(func(ev testInterface1) {
		triggered = true
	})
	Dispatch(&testEvent1{})

	if !triggered {
		t.Errorf("Dispatch by pointer failed to trigger interface listener")
	}
}

func TestDispatchValueToValueInterfaceListener(t *testing.T) {
	clearListeners()

	triggered := false
	AddListener(func(ev testInterface1) {
		triggered = true
	})
	Dispatch(testEvent1{})

	if !triggered {
		t.Errorf("Dispatch by value failed to trigger interface listener")
	}
}

func TestDispatchPointerToPointerInterfaceListener(t *testing.T) {
	clearListeners()

	triggered := false
	AddListener(func(testInterface2) { triggered = true })
	Dispatch(&testEvent2{})

	if !triggered {
		t.Errorf("interface listener failed to trigger for pointer")
	}
}

func TestDispatchValueToPointerInterfaceListener(t *testing.T) {
	clearListeners()

	AddListener(func(testInterface2) {
		t.Errorf("interface listener triggered for value dispatch")
	})
	Dispatch(testEvent2{})
}

type testUpdateEvent struct {
	update interface{}
}

func (ev *testUpdateEvent) Update(update interface{}) {
	ev.update = update
}

func TestDispatchUpdate(t *testing.T) {
	clearListeners()

	triggered := false
	AddListener(func(*testUpdateEvent) {
		triggered = true
	})

	ev := &testUpdateEvent{}
	DispatchUpdate(ev, "hello")

	if !triggered {
		t.Errorf("listener failed to trigger on DispatchUpdate()")
	}
	want := "hello"
	if got := ev.update.(string); got != want {
		t.Errorf("ev.update = %#v, want %#v", got, want)
	}
}
