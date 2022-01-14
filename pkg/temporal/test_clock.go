package temporal

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
	"sync"
	"time"
)

var (
	// testClock is the global test instance of the clock for tests.
	testClock *TestClock
)

// TestClock is an implementation of Clock for tests, where it is
// possible to set the current time and the uncertainty at any time.
//
// To use it:
// time.UseTestClock()
// time.SetTestClockTime(now)
// time.SetTestClockUncertainty(dur)
type TestClock struct {
	// mu protects the following fields
	mu          sync.Mutex
	now         time.Time
	uncertainty time.Duration
}

// Now is part of the Clock interface.
func (t *TestClock) Now() (Interval, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return NewInterval(t.now.Add(-(t.uncertainty)), t.now.Add(t.uncertainty))
}

// Set let the user set the time
func (t *TestClock) Set(now time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.now = now
}

// SetUncertainty lets the user set the uncertainty
func (t *TestClock) SetUncertainty(uncertainty time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.uncertainty = uncertainty
}

// SetTestClockTime sets the 'test' implementation time to the provided value.
func SetTestClockTime(now time.Time) {
	testClock.Set(now)
}

// SetTestClockUncertainty sets the 'test' implementation uncertainty
// to the provided value.
func SetTestClockUncertainty(uncertainty time.Duration) {
	testClock.SetUncertainty(uncertainty)
}

// UseTestClock is meant to be used in tests to start using the test clock.
func UseTestClock() {
	*defaultClockType = "test"
}

func init() {
	testClock = &TestClock{
		now:         time.Now(),
		uncertainty: 10 * time.Millisecond,
	}
	clockTypes["test"] = testClock
}
