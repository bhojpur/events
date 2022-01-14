package syslogger

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

/*
Package syslogger uses the Events engine package to listen for any event that
implements the Syslogger interface. The listener calls the Syslog method on the
event, which should return a severity and a message.

The tag for messages sent to syslog will be the program name (os.Args[0]),
and the facility number will be 1 (user-level).

For example, to declare that your event type MyEvent should be sent to syslog,
implement the Syslog() method to define how the message should be formatted and
which severity it should have (see package "log/syslog" for details).

	import (
		"fmt"
		"log/syslog"
		"syslogger"
	)

	type MyEvent struct {
		field1, field2 string
	}

	func (ev *MyEvent) Syslog() (syslog.Priority, string) {
		return syslog.LOG_INFO, fmt.Sprintf("event: %v, %v", ev.field1, ev.field2)
	}
	var _ syslogger.Syslogger = (*MyEvent)(nil) // compile-time interface check

The compile-time interface check is optional but recommended because usually
there is no other static conversion in these cases.
*/

import (
	"fmt"
	"log/syslog"
	"os"

	"github.com/bhojpur/events/pkg/engine"
	"github.com/bhojpur/events/pkg/log"
)

// Syslogger is the interface that events should implement if they want to be
// dispatched to this package.
type Syslogger interface {
	// Syslog should return a severity (not a facility) and a message.
	Syslog() (syslog.Priority, string)
}

// syslogWriter is an interface that wraps syslog.Writer so it can be faked.
type syslogWriter interface {
	Alert(string) error
	Crit(string) error
	Debug(string) error
	Emerg(string) error
	Err(string) error
	Info(string) error
	Notice(string) error
	Warning(string) error
}

// writer holds a persistent connection to the syslog daemon
var writer syslogWriter

func listener(ev Syslogger) {
	// Ask the event to convert itself to a syslog message.
	sev, msg := ev.Syslog()

	// Call the corresponding Writer function.
	var err error
	switch sev {
	case syslog.LOG_EMERG:
		if writer != nil {
			err = writer.Emerg(msg)
		} else {
			log.Errorf(msg)
		}
	case syslog.LOG_ALERT:
		if writer != nil {
			err = writer.Alert(msg)
		} else {
			log.Errorf(msg)
		}
	case syslog.LOG_CRIT:
		if writer != nil {
			err = writer.Crit(msg)
		} else {
			log.Errorf(msg)
		}
	case syslog.LOG_ERR:
		if writer != nil {
			err = writer.Err(msg)
		} else {
			log.Errorf(msg)
		}
	case syslog.LOG_WARNING:
		if writer != nil {
			err = writer.Warning(msg)
		} else {
			log.Warningf(msg)
		}
	case syslog.LOG_NOTICE:
		if writer != nil {
			err = writer.Notice(msg)
		} else {
			log.Infof(msg)
		}
	case syslog.LOG_INFO:
		if writer != nil {
			err = writer.Info(msg)
		} else {
			log.Infof(msg)
		}
	case syslog.LOG_DEBUG:
		if writer != nil {
			err = writer.Debug(msg)
		} else {
			log.Infof(msg)
		}
	default:
		err = fmt.Errorf("invalid syslog severity: %v", sev)
	}
	if err != nil {
		log.Errorf("can't write syslog event: %v", err)
	}
}

func init() {
	var err error
	writer, err = syslog.New(syslog.LOG_INFO|syslog.LOG_USER, os.Args[0])
	if err != nil {
		log.Errorf("can't connect to syslog")
		writer = nil
	}

	engine.AddListener(listener)
}
