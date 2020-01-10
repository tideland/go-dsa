// Tideland Go Data Structures and Algorithms - Collections - Sets
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
// SET
//--------------------

// Set defines a set containing any kind of values.
type Set struct {
	values map[interface{}]struct{}
}

// NewSet creates a set with the passed values
// as initial content.
func NewSet(vs ...interface{}) *Set {
	s := &Set{make(map[interface{}]struct{})}
	s.Add(vs...)
	return s
}

// Add adds values to the set.
func (s *Set) Add(vs ...interface{}) {
	for _, v := range vs {
		s.values[v] = struct{}{}
	}
}

// Remove removes values out if the set. It doesn't
// matter if the set does not contain them.
func (s *Set) Remove(vs ...interface{}) {
	for _, v := range vs {
		delete(s.values, v)
	}
}

// Contains checks if the set contains a given value.
func (s *Set) Contains(v interface{}) bool {
	_, ok := s.values[v]
	return ok
}

// All returns all values.
func (s *Set) All() []interface{} {
	all := []interface{}{}
	for v := range s.values {
		all = append(all, v)
	}
	return all
}

// FindAll returns all values found by the passed function.
func (s *Set) FindAll(f func(v interface{}) (bool, error)) ([]interface{}, error) {
	found := []interface{}{}
	for v := range s.values {
		ok, err := f(v)
		if err != nil {
			return nil, failure.Annotate(err, "cannot find all matching values")
		}
		if ok {
			found = append(found, v)
		}
	}
	return found, nil
}

// DoAll executes the passed function on all values.
func (s *Set) DoAll(f func(v interface{}) error) error {
	for v := range s.values {
		if err := f(v); err != nil {
			return failure.Annotate(err, "cannot process all entries")
		}
	}
	return nil
}

// Len returns the number of entries in the set.
func (s *Set) Len() int {
	return len(s.values)
}

// Deflate cleans the stack.
func (s *Set) Deflate() {
	s.values = make(map[interface{}]struct{})
}

// String implements the fmt.Stringer interface.
func (s *Set) String() string {
	all := s.All()
	return fmt.Sprintf("%v", all)
}

//--------------------
// STRING SET
//--------------------

// StringSet defines a set containing string values.
type StringSet struct {
	values map[string]struct{}
}

// NewStringSet creates a string set with the passed strings
// as initial content.
func NewStringSet(vs ...string) *StringSet {
	s := &StringSet{make(map[string]struct{})}
	s.Add(vs...)
	return s
}

// Add adds strings to the set.
func (s *StringSet) Add(vs ...string) {
	for _, v := range vs {
		s.values[v] = struct{}{}
	}
}

// Remove removes strings out if the set. It doesn't
// matter if the set does not contain them.
func (s *StringSet) Remove(vs ...string) {
	for _, v := range vs {
		delete(s.values, v)
	}
}

// Contains checks if the set contains a given string.
func (s *StringSet) Contains(v string) bool {
	_, ok := s.values[v]
	return ok
}

// All returns all strings.
func (s *StringSet) All() []string {
	all := []string{}
	for v := range s.values {
		all = append(all, v)
	}
	return all
}

// FindAll returns all strings found by the passed function.
func (s *StringSet) FindAll(f func(v string) (bool, error)) ([]string, error) {
	found := []string{}
	for v := range s.values {
		ok, err := f(v)
		if err != nil {
			return nil, failure.Annotate(err, "cannot find all matching string values")
		}
		if ok {
			found = append(found, v)
		}
	}
	return found, nil
}

// DoAll executes the passed function on all strings.
func (s *StringSet) DoAll(f func(v string) error) error {
	for v := range s.values {
		if err := f(v); err != nil {
			return failure.Annotate(err, "cannot process all string entries")
		}
	}
	return nil
}

// Len returns the number of entries in the set.
func (s *StringSet) Len() int {
	return len(s.values)
}

// Deflate cleans the stack.
func (s *StringSet) Deflate() {
	s.values = make(map[string]struct{})
}

// String implements the fmt.Stringer interface.
func (s *StringSet) String() string {
	all := s.All()
	return fmt.Sprintf("%v", all)
}

// EOF
