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

	"github.com/bhojpur/events/pkg/log"
)

var (
	// clockTypes maps implementation name to Clock object.
	// Should only be written to at init() time.
	clockTypes = make(map[string]Clock)

	// defaultClockType is the flag used to define the runtime clock type.
	defaultClockType = flag.String("time_default_clock_type", "time", "The type of clock to be used by default time library.")
)

// Clock returns the current time.
type Clock interface {
	// Now returns the current time as Interval.
	// This method should be thread safe (i.e. multiple go routines can
	// safely call this at the same time).
	// The returned interval is guaranteed to have earliest <= latest,
	// and all implementations enforce it.
	Now() (Interval, error)
}

// GetClock returns the global Clock object.
// Since it depends on flags, be sure to call this after they have been parsed
// (i.e. *not* in init() functions), otherwise this will panic.
func GetClock() Clock {
	if !flag.Parsed() {
		panic("GetClock() called before flags are parsed")
	}

	c, ok := clockTypes[*defaultClockType]
	if !ok {
		log.Fatalf("No Clock type named %v", *defaultClockType)
	}
	return c
}
