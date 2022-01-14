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
	"flag"
	"time"
)

var (
	uncertainty = flag.Duration("time_time_clock_uncertainty", 10*time.Millisecond, "An assumed time uncertainty on the local machine that will be used by time-based implementation of Clock.")
)

// TimeClock is an implementation of Clock that uses time.Now() and a
// flag-configured uncertainty.
type TimeClock struct{}

// Now is part of the Clock interface.
func (t TimeClock) Now() (Interval, error) {
	now := time.Now()
	return NewInterval(now.Add(-(*uncertainty)), now.Add(*uncertainty))
}

func init() {
	clockTypes["time"] = TimeClock{}
}
