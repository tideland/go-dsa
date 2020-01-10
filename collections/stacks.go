// Tideland Go Data Structures and Algorithms - Collections - Stacks
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

	"tideland.dev/go/trace/failure"
)

//--------------------
// STACK
//--------------------

// Stack defines a stack containing any kind of values.
type Stack struct {
	values []interface{}
}

// NewStack creates a stack with the passed values
// as initial content.
func NewStack(vs ...interface{}) *Stack {
	return &Stack{
		values: vs,
	}
}

// Push adds values to the top of the stack.
func (s *Stack) Push(vs ...interface{}) {
	s.values = append(s.values, vs...)
}

// Pop removes and returns the top value of the stack.
func (s *Stack) Pop() (interface{}, error) {
	lv := len(s.values)
	if lv == 0 {
		return nil, failure.New("stack is empty")
	}
	v := s.values[lv-1]
	s.values = s.values[:lv-1]
	return v, nil
}

// Peek returns the top value of the stack.
func (s *Stack) Peek() (interface{}, error) {
	lv := len(s.values)
	if lv == 0 {
		return nil, failure.New("stack is empty")
	}
	v := s.values[lv-1]
	return v, nil
}

// All returns all values bottom-up.
func (s *Stack) All() []interface{} {
	sl := len(s.values)
	all := make([]interface{}, sl)
	copy(all, s.values)
	return all
}

// AllReverse returns all values top-down.
func (s *Stack) AllReverse() []interface{} {
	sl := len(s.values)
	all := make([]interface{}, sl)
	for i, value := range s.values {
		all[sl-1-i] = value
	}
	return all
}

// Len returns the number of entries in the stack.
func (s *Stack) Len() int {
	return len(s.values)
}

// Deflate cleans the stack.
func (s *Stack) Deflate() {
	s.values = []interface{}{}
}

// String implements the fmt.Stringer interface.
func (s *Stack) String() string {
	return fmt.Sprintf("%v", s.values)
}

//--------------------
// STRING STACK
//--------------------

// StringStack defines a stack containing string values.
type StringStack struct {
	values []string
}

// NewStringStack creates a string stack with the passed strings
// as initial content.
func NewStringStack(vs ...string) *StringStack {
	return &StringStack{
		values: vs,
	}
}

// Push adds strings to the top of the stack.
func (s *StringStack) Push(vs ...string) {
	s.values = append(s.values, vs...)
}

// Pop removes and returns the top string of the stack.
func (s *StringStack) Pop() (string, error) {
	lv := len(s.values)
	if lv == 0 {
		return "", failure.New("string stack is empty")
	}
	v := s.values[lv-1]
	s.values = s.values[:lv-1]
	return v, nil
}

// Peek returns the top string of the stack.
func (s *StringStack) Peek() (string, error) {
	lv := len(s.values)
	if lv == 0 {
		return "", failure.New("string stack is empty")
	}
	v := s.values[lv-1]
	return v, nil
}

// All returns all strings bottom-up.
func (s *StringStack) All() []string {
	sl := len(s.values)
	all := make([]string, sl)
	copy(all, s.values)
	return all
}

// AllReverse returns all strings top-down.
func (s *StringStack) AllReverse() []string {
	sl := len(s.values)
	all := make([]string, sl)
	for i, value := range s.values {
		all[sl-1-i] = value
	}
	return all
}

// Len returns the number of entries in the stack.
func (s *StringStack) Len() int {
	return len(s.values)
}

// Deflate cleans the stack.
func (s *StringStack) Deflate() {
	s.values = []string{}
}

// String implements the fmt.Stringer interface.
func (s *StringStack) String() string {
	return fmt.Sprintf("%v", s.values)
}

// EOF
