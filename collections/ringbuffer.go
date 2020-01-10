// Tideland Go Data Structures and Algorithms - Collections - Ring Buffer
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package collections // import "tideland.dev/go/dsa/collections"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"strings"
)

//--------------------
// RING BUFFER
//--------------------

// valueLink is one ring buffer element containing one
// value and linking to the next element.
type valueLink struct {
	used  bool
	value interface{}
	next  *valueLink
}

// RingBuffer defines a buffer which is connected end-to-end. It
// grows if needed.
type RingBuffer struct {
	start *valueLink
	end   *valueLink
}

// NewRingBuffer creates a new ring buffer.
func NewRingBuffer(size int) *RingBuffer {
	rb := &RingBuffer{}
	rb.start = &valueLink{}
	rb.end = rb.start
	if size < 2 {
		size = 2
	}
	for i := 0; i < size-1; i++ {
		link := &valueLink{}
		rb.end.next = link
		rb.end = link
	}
	rb.end.next = rb.start
	return rb
}

// Push adds values to the end of the buffer.
func (rb *RingBuffer) Push(values ...interface{}) {
	for _, value := range values {
		if !rb.end.next.used {
			rb.end.next.used = true
			rb.end.next.value = value
			rb.end = rb.end.next
			continue
		}
		link := &valueLink{
			used:  true,
			value: value,
			next:  rb.start,
		}
		rb.end.next = link
		rb.end = rb.end.next
	}
}

// Peek returns the first value of the buffer. If the
// buffer is empty the second return value is false.
func (rb *RingBuffer) Peek() (interface{}, bool) {
	if !rb.start.used {
		return nil, false
	}
	return rb.start.value, true
}

// Pop removes and returns the first value of the buffer. If
// the buffer is empty the second return value is false.
func (rb *RingBuffer) Pop() (interface{}, bool) {
	if !rb.start.used {
		return nil, false
	}
	value := rb.start.value
	rb.start.used = false
	rb.start.value = nil
	rb.start = rb.start.next
	return value, true
}

// Len returns the number of values in the buffer.
func (rb *RingBuffer) Len() int {
	l := 0
	current := rb.start
	for current.used {
		l++
		current = current.next
		if current == rb.start {
			break
		}
	}
	return l
}

// Cap returns the capacity of the buffer.
func (rb *RingBuffer) Cap() int {
	c := 1
	current := rb.start
	for current.next != rb.start {
		c++
		current = current.next
	}
	return c
}

// String implements the fmt.Stringer interface.
func (rb *RingBuffer) String() string {
	vs := []string{}
	current := rb.start
	for current.used {
		vs = append(vs, fmt.Sprintf("[%v]", current.value))
		current = current.next
		if current == rb.start {
			break
		}
	}
	return strings.Join(vs, "->")
}

// EOF
