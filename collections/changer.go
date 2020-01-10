// Tideland Go Data Structures and Algorithms - Collections
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
	"tideland.dev/go/trace/failure"
)

//--------------------
// CHANGER
//--------------------

// Changer defines the interface to perform changes on a tree
// node. It is returned by the addressing operations like At() and
// Create() of the Tree.
type Changer struct {
	node *node
	err  error
}

// Value returns the changer node value.
func (c *Changer) Value() (interface{}, error) {
	if c.err != nil {
		return nil, c.err
	}
	return c.node.content.value(), nil
}

// SetValue sets the changer node value. It also returns
// the previous value.
func (c *Changer) SetValue(v interface{}) (interface{}, error) {
	if c.err != nil {
		return nil, c.err
	}
	oldValue := c.node.content.value()
	newValue := justValue{v}
	if !c.node.isAllowed(newValue, true) {
		return nil, failure.New("setting duplicate value is not allowed")
	}
	c.node.content = newValue
	return oldValue, nil
}

// Add sets a child value.
func (c *Changer) Add(v interface{}) error {
	if c.err != nil {
		return c.err
	}
	_, err := c.node.addChild(justValue{v})
	return err
}

// Remove deletes this changer node.
func (c *Changer) Remove() error {
	if c.err != nil {
		return c.err
	}
	return c.node.remove()
}

// List returns the values of the children of the changer node.
func (c *Changer) List() ([]interface{}, error) {
	if c.err != nil {
		return nil, c.err
	}
	var list []interface{}
	err := c.node.doChildren(func(cn *node) error {
		list = append(list, cn.content.value())
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Error returns a potential error of the changer.
func (c *Changer) Error() error {
	return c.err
}

//--------------------
// STRING CHANGER
//--------------------

// StringChanger defines the interface to perform changes on a string
// tree node. It is returned by the addressing operations like
// At() and Create() of the StringTree.
type StringChanger struct {
	node *node
	err  error
}

// Value returns the changer node value.
func (c *StringChanger) Value() (string, error) {
	if c.err != nil {
		return "", c.err
	}
	if c.node.content.value() == nil {
		return "", nil
	}
	return c.node.content.value().(string), nil
}

// SetValue sets the changer node value. It also returns
// the previous value.
func (c *StringChanger) SetValue(v string) (string, error) {
	if c.err != nil {
		return "", c.err
	}
	oldValue := c.node.content.value().(string)
	newValue := justValue{v}
	if !c.node.isAllowed(newValue, true) {
		return "", failure.New("setting duplicate string value is not allowed")
	}
	c.node.content = newValue
	return oldValue, nil
}

// Add sets a child value. If it already exists it will be overwritten.
func (c *StringChanger) Add(v string) error {
	if c.err != nil {
		return c.err
	}
	_, err := c.node.addChild(justValue{v})
	return err
}

// Remove deletes this changer node.
func (c *StringChanger) Remove() error {
	if c.err != nil {
		return c.err
	}
	return c.node.remove()
}

// List returns the values of the children of the changer node.
func (c *StringChanger) List() ([]string, error) {
	if c.err != nil {
		return nil, c.err
	}
	var list []string
	err := c.node.doChildren(func(cn *node) error {
		list = append(list, cn.content.value().(string))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Error returns a potential error of the changer.
func (c *StringChanger) Error() error {
	return c.err
}

//--------------------
// KEY/VALUE CHANGER
//--------------------

// KeyValueChanger defines the interface to perform changes on a
// key/value tree node. It is returned by the addressing operations
// like At() and Create() of the KeyValueTree.
type KeyValueChanger struct {
	node *node
	err  error
}

// Key returns the changer node key.
func (c *KeyValueChanger) Key() (string, error) {
	if c.err != nil {
		return "", c.err
	}
	return c.node.content.key().(string), nil
}

// SetKey sets the changer node key. Its checks if duplicate
// keys are allowed and returns the previous key.
func (c *KeyValueChanger) SetKey(key string) (string, error) {
	if c.err != nil {
		return "", c.err
	}
	if !c.node.container.duplicates {
		if c.node.hasDuplicateSibling(key) {
			return "", failure.New("setting duplicate key is not allowed")
		}
	}
	current := c.node.content.key().(string)
	c.node.content = keyValue{key, c.node.content.value()}
	return current, nil
}

// Value returns the changer node value.
func (c *KeyValueChanger) Value() (interface{}, error) {
	if c.err != nil {
		return nil, c.err
	}
	return c.node.content.value(), nil
}

// SetValue sets the changer node value. It also returns
// the previous value.
func (c *KeyValueChanger) SetValue(value interface{}) (interface{}, error) {
	if c.err != nil {
		return nil, c.err
	}
	current := c.node.content.value()
	c.node.content = keyValue{c.node.content.key(), value}
	return current, nil
}

// Add sets a child key/value. If the key already exists the
// value will be overwritten.
func (c *KeyValueChanger) Add(k string, v interface{}) error {
	if c.err != nil {
		return c.err
	}
	_, err := c.node.addChild(keyValue{k, v})
	return err
}

// Remove deletes this changer node.
func (c *KeyValueChanger) Remove() error {
	if c.err != nil {
		return c.err
	}
	return c.node.remove()
}

// List returns the keys and values of the children of the changer node.
func (c *KeyValueChanger) List() ([]KeyValue, error) {
	if c.err != nil {
		return nil, c.err
	}
	var list []KeyValue
	err := c.node.doChildren(func(cn *node) error {
		list = append(list, KeyValue{cn.content.key().(string), cn.content.value()})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Error returns a potential error of the changer.
func (c *KeyValueChanger) Error() error {
	return c.err
}

//--------------------
// KEY/STRING VALUE CHANGER
//--------------------

// KeyStringValueChanger defines the interface to perform changes
// on a key/string value tree node. It is returned by the addressing
// operations like At() and Create() of the KeyStringValueTree.
type KeyStringValueChanger struct {
	node *node
	err  error
}

// Key returns the changer node key.
func (c *KeyStringValueChanger) Key() (string, error) {
	if c.err != nil {
		return "", c.err
	}
	return c.node.content.key().(string), nil
}

// SetKey sets the changer node key. Its checks if duplicate
// keys are allowed and returns the previous key.
func (c *KeyStringValueChanger) SetKey(key string) (string, error) {
	if c.err != nil {
		return "", c.err
	}
	if !c.node.container.duplicates {
		if c.node.hasDuplicateSibling(key) {
			return "", failure.New("setting duplicate key is not allowed")
		}
	}
	current := c.node.content.key().(string)
	c.node.content = keyValue{key, c.node.content.value()}
	return current, nil
}

// Value returns the changer node value.
func (c *KeyStringValueChanger) Value() (string, error) {
	if c.err != nil {
		return "", c.err
	}
	if c.node.content.value() == nil {
		return "", nil
	}
	return c.node.content.value().(string), nil
}

// SetValue sets the changer node value. It also returns
// the previous value.
func (c *KeyStringValueChanger) SetValue(value string) (string, error) {
	if c.err != nil {
		return "", c.err
	}
	current := c.node.content.value().(string)
	c.node.content = keyValue{c.node.content.key(), value}
	return current, nil
}

// Add sets a child key/value. If the key already exists the
// value will be overwritten.
func (c *KeyStringValueChanger) Add(k, v string) error {
	if c.err != nil {
		return c.err
	}
	_, err := c.node.addChild(keyValue{k, v})
	return err
}

// Remove deletes this changer node.
func (c *KeyStringValueChanger) Remove() error {
	if c.err != nil {
		return c.err
	}
	return c.node.remove()
}

// List returns the keys and values of the children of the changer node.
func (c *KeyStringValueChanger) List() ([]KeyStringValue, error) {
	if c.err != nil {
		return nil, c.err
	}
	var list []KeyStringValue
	err := c.node.doChildren(func(cn *node) error {
		list = append(list, KeyStringValue{cn.content.key().(string), cn.content.value().(string)})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Error returns a potential error of the changer.
func (c *KeyStringValueChanger) Error() error {
	return c.err
}

// EOF
