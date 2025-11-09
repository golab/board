/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// the core package provides basic functionality to all the major components of the code
package message

import (
	"sync"
	"time"
)

type Message struct {
	Text      string
	ExpiresAt *time.Time
	notified  map[string]bool
	mu        sync.Mutex
}

func NewMessage(text string, ttl int) *Message {
	// calculate the expiration time using TTL
	now := time.Now()
	expiresAt := now.Add(time.Duration(ttl) * time.Second)

	return &Message{
		Text:      text,
		ExpiresAt: &expiresAt,
		notified:  make(map[string]bool),
		mu:        sync.Mutex{},
	}
}

func (m *Message) MarkNotified(id string) {
	m.mu.Lock()
	m.notified[id] = true
	m.mu.Unlock()
}

func (m *Message) IsNotified(id string) bool {
	var n bool
	m.mu.Lock()
	_, n = m.notified[id]
	m.mu.Unlock()
	return n
}
