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
	"testing"
	"time"
)

// Only checking that the results are somewhat OK.
func TestTimeClock(t *testing.T) {
	flag.Set("time_default_clock_type", "time")
	flag.Set("time_time_clock_uncertainty", "20ms")

	clock := GetClock()
	i, err := clock.Now()
	if err != nil {
		t.Fatalf("Now failed: %v", err)
	}

	// Check the difference is exactly right: it needs to be twice
	// the uncertainty set above with the time_time_clock_uncertainty
	// flag.
	d := i.Latest().Sub(i.Earliest())
	if d != 40*time.Millisecond {
		t.Errorf("uncertainty not respected: %v", d)
	}

	// Check we're somewhat in range with time().
	now := time.Now()
	d = now.Sub(i.Earliest())
	if d.Seconds() > 1 || d.Seconds() < -1 {
		t.Errorf("now very late: %v %v %v", i, now, d)
	}
	d = now.Sub(i.Latest())
	if d.Seconds() > 1 || d.Seconds() < -1 {
		t.Errorf("now very early %v %v %v", now, i, d)
	}
}
