// Copyright 2025 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package runtime

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"
	"slices"
	"strings"
	"sync/atomic"
	"time"
)

// An Event sent by HAProxy.
type Event struct {
	Timestamp time.Time
	Message   string
}

func (e Event) String() string {
	return e.Timestamp.Format(time.StampMicro) + " " + e.Message
}

// An EventListener connects to HAProxy's master socket and listens
// for events on a specific sink using the "show events" command.
type EventListener struct {
	conn net.Conn
	// Buffered reader for conn.
	reader *bufio.Reader
	// Network to use with net.Dial (unix, udp...).
	network string
	// Address or path to HAProxy's master socket.
	address string
	// HAProxy's sink to read events from.
	sink string
	// Command used to listen for events.
	listenCmd string
	// Format used to parse events timestamps. Defaults to HAProxy's LOG_FORMAT_ISO.
	// Can be set before calling Listen.
	DateFormat string
	// Write timeout. Can be set before calling Listen.
	WriteTimeout time.Duration
	// Message delimiter. Either \n or 0 (zero).
	delim     byte
	events    chan Event
	lastError error
	closed    atomic.Bool
}

// NewEventListener connects to HAProxy's master socket and returns a new
// EventListener configured for the given sink and flags. The timeout parameter
// is used both as the default WriteTimeout and as the connect timeout.
//
// The caller must call Close to properly shutdown an EventListener.
func NewEventListener(network, address, sink string, timeout time.Duration, flags ...string) (*EventListener, error) {
	var err error

	l := &EventListener{
		network:      network,
		address:      address,
		sink:         sink,
		listenCmd:    fmt.Sprintf("show events %s %s", sink, strings.Join(flags, " ")),
		DateFormat:   "2006-01-02T15:04:05.000000-07:00", // LOG_FORMAT_TIMED, LOG_FORMAT_ISO
		WriteTimeout: timeout,
		delim:        '\n',
		events:       make(chan Event),
	}

	if slices.Contains(flags, "-0") {
		l.delim = byte(0)
	}

	l.conn, err = net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, l.errorf("%w", err)
	}

	l.reader = bufio.NewReader(l.conn)

	go l.listen()

	return l, nil
}

// Listen for for new events. Blocks until a new event is available,
// or until the ctx deadline is reached. On success, it returns the
// event's payload without its timestamp.
//
// Listen can return an error in the following cases: a network error
// (of type *net.OpError), EOF, an error returned by HAProxy or a parsing error.
//
// In case of a network error, the caller should Close the EventListener
// and create a new one to continue to receive events.
func (l *EventListener) Listen(ctx context.Context) (Event, error) {
	select {
	case event, ok := <-l.events:
		if !ok {
			return Event{}, l.lastError
		}
		return event, nil
	case <-ctx.Done():
		return Event{}, l.Close()
	}
}

// Close the EventListener cleanly.
func (l *EventListener) Close() error {
	if l.closed.CompareAndSwap(false, true) {
		defer close(l.events)
		if err := l.conn.Close(); err != nil {
			return l.errorf("%w", err)
		}
	}

	return nil
}

// Listen for events and push them on the events channel.
func (l *EventListener) listen() {
	if l.WriteTimeout > 0 {
		_ = l.conn.SetWriteDeadline(time.Now().Add(l.WriteTimeout))
	}

	_, err := fmt.Fprintf(l.conn, "@@1 set severity-output number;%s\n", l.listenCmd)
	if err != nil {
		l.lastError = l.errorf("%w", err)
		_ = l.Close()
		return
	}

	for {
		data, err := l.reader.ReadBytes(l.delim)
		if err != nil {
			l.lastError = l.errorf("%w", err)
			_ = l.Close()
			return
		}

		event, err := l.parseEvent(data)
		if err != nil {
			l.lastError = err
			_ = l.Close()
			return
		}

		l.events <- event
	}
}

func (l *EventListener) parseEvent(data []byte) (Event, error) {
	data = bytes.TrimPrefix(data, []byte{'\n'})
	data = bytes.TrimPrefix(data, []byte("<0>")) // syslog artefact
	if l.delim == 0 {
		data = bytes.TrimSuffix(data, []byte("\n\x00"))
	}
	data = bytes.TrimSpace(data)

	event := string(data)

	if len(event) > 4 && event[0] == '[' {
		switch event[0:4] {
		case "[3]:", "[2]:", "[1]:", "[0]:":
			return Event{}, l.errorf("[%c] %s [%s]", event[1], event[4:], l.listenCmd)
		}
	}

	timestamp, msg, found := strings.Cut(event, " ")
	if !found {
		return Event{}, l.errorf("parsing error: '%s'", event)
	}

	ts, err := time.Parse(l.DateFormat, timestamp)
	if err != nil {
		return Event{}, l.errorf("time parsing error: '%s'", timestamp)
	}

	return Event{Timestamp: ts, Message: msg}, nil
}

func (l *EventListener) errorf(format string, a ...any) error {
	format = fmt.Sprintf("EventListener(%s): %s", l.sink, format)
	return fmt.Errorf(format, a...)
}
