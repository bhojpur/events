package status

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
	"testing"

	"github.com/bhojpur/events/pkg/engine"
)

type testEvent struct {
	StatusUpdater
}

func TestUpdateInit(t *testing.T) {
	want := "status"
	ev := &testEvent{}
	ev.Update("status")

	if ev.Status != want {
		t.Errorf("ev.Status = %#v, want %#v", ev.Status, want)
	}
	if ev.EventID == 0 {
		t.Errorf("ev.EventID wasn't initialized")
	}
}

func TestUpdateEventID(t *testing.T) {
	want := int64(12345)
	ev := &testEvent{}
	ev.EventID = 12345

	ev.Update("status")

	if ev.EventID != want {
		t.Errorf("ev.EventID = %v, want %v", ev.EventID, want)
	}
}

func TestUpdateDispatch(t *testing.T) {
	triggered := false
	engine.AddListener(func(ev *testEvent) {
		triggered = true
	})

	want := "status"
	ev := &testEvent{}
	engine.DispatchUpdate(ev, "status")

	if ev.Status != want {
		t.Errorf("ev.Status = %#v, want %#v", ev.Status, want)
	}
	if !triggered {
		t.Errorf("listener wasn't triggered on Dispatch()")
	}
}
