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
	"fmt"

	"tideland.dev/go/trace/failure"
)

//--------------------
// EXCHANGE TYPES
//--------------------

// KeyValue wraps a key and a value for the key/value iterator.
type KeyValue struct {
	Keys  string
	Value interface{}
}

// KeyStringValue carries a combination of key and string value.
type KeyStringValue struct {
	Key   string
	Value string
}

//--------------------
// NODE VALUE
//--------------------

// nodeContent is the base interface for the different value
// containing types.
type nodeContent interface {
	// key returns the key for finding ands
	// check for duplicates.
	key() interface{}

	// value returns the value itself.
	value() interface{}

	// deepCopy creates a copy of the node content.
	deepCopy() nodeContent
}

// justValue has same key and value.
type justValue struct {
	v interface{}
}

// key implements the nodeContent interface.
func (v justValue) key() interface{} {
	return v.v
}

// value implements the nodeContent interface.
func (v justValue) value() interface{} {
	return v.v
}

// deepCopy implements the nodeContent interface.
func (v justValue) deepCopy() nodeContent {
	return justValue{v.v}
}

// String implements the Stringer interface.
func (v justValue) String() string {
	return fmt.Sprintf("%v", v.v)
}

// keyValue has different key and value.
type keyValue struct {
	k interface{}
	v interface{}
}

// key implements the nodeContent interface.
func (v keyValue) key() interface{} {
	return v.k
}

// value implements the nodeContent interface.
func (v keyValue) value() interface{} {
	return v.v
}

// deepCopy implements the nodeContent interface.
func (v keyValue) deepCopy() nodeContent {
	return keyValue{v.k, v.v}
}

// String implements the Stringer interface.
func (v keyValue) String() string {
	return fmt.Sprintf("%v = '%v'", v.k, v.v)
}

//--------------------
// NODE CONTAINER
//--------------------

// nodeContainer is the top element of all nodes and provides
// configuration.
type nodeContainer struct {
	root       *node
	duplicates bool
}

// newNodeContainer creates a new node container.
func newNodeContainer(c nodeContent, duplicates bool) *nodeContainer {
	nc := &nodeContainer{
		root: &node{
			content: c,
		},
		duplicates: duplicates,
	}
	nc.root.container = nc
	return nc
}

// deepCopy creates a copy of the container.
func (nc *nodeContainer) deepCopy() *nodeContainer {
	cnc := &nodeContainer{
		duplicates: nc.duplicates,
	}
	cnc.root = nc.root.deepCopy(cnc, nil)
	return cnc
}

//--------------------
// NODE
//--------------------

// node contains the value and structural information of a node.
type node struct {
	container *nodeContainer
	parent    *node
	content   nodeContent
	children  []*node
}

// isAllowed returns true, if adding the content or setting
// it is allowed depending on allowed duplicates.
func (n *node) isAllowed(c nodeContent, here bool) bool {
	if n.container.duplicates {
		return true
	}
	checkNode := n
	if here {
		checkNode = n.parent
	}
	for _, child := range checkNode.children {
		if child.content.key() == c.key() {
			return false
		}
	}
	return true
}

// hasDuplicateSibling checks if the node has a sibling with the same key.
func (n *node) hasDuplicateSibling(key interface{}) bool {
	if n.parent == nil {
		return false
	}
	for _, sibling := range n.parent.children {
		if sibling == n {
			continue
		}
		if sibling.content.key() == key {
			return true
		}
	}
	return false
}

// addChild adds a child node depending on allowed duplicates.
func (n *node) addChild(c nodeContent) (*node, error) {
	if !n.isAllowed(c, false) {
		return nil, failure.New("adding duplicate node is not allowed")
	}
	child := &node{
		container: n.container,
		parent:    n,
		content:   c,
	}
	n.children = append(n.children, child)
	return child, nil
}

// remove deletes this node from its parent.
func (n *node) remove() error {
	if n.parent == nil {
		return failure.New("cannot remove root node")
	}
	for i, child := range n.parent.children {
		if child == n {
			n.parent.children = append(n.parent.children[:i], n.parent.children[i+1:]...)
			return nil
		}
	}
	panic("cannot find node to remove at parent")
}

// at finds a node by its path.
func (n *node) at(path ...nodeContent) (*node, error) {
	if len(path) == 0 || path[0].key() != n.content.key() {
		return nil, failure.New("cannot find node")
	}
	if len(path) == 1 {
		return n, nil
	}
	// Check children for rest of the path.
	for _, child := range n.children {
		found, err := child.at(path[1:]...)
		if err != nil && !failure.Contains(err, "cannot find") {
			return nil, failure.Annotate(err, "invalid path")
		}
		if found != nil {
			return found, nil
		}
	}
	return nil, failure.New("cannot find node")
}

// create acts like at but if nodes don't exist they will be created.
func (n *node) create(path ...nodeContent) (*node, error) {
	if len(path) == 0 || path[0].key() != n.content.key() {
		return nil, failure.New("cannot find parent node for creation")
	}
	if len(path) == 1 {
		return n, nil
	}
	// Check children for the next path element.
	var found *node
	for _, child := range n.children {
		if path[1].key() == child.content.key() {
			found = child
			break
		}
	}
	if found == nil {
		child, err := n.addChild(path[1])
		if err != nil {
			return nil, failure.Annotate(err, "cannot add child node")
		}
		return child.create(path[1:]...)
	}
	return found.create(path[1:]...)
}

// findFirst returns the first node for which the passed function
// returns true.
func (n *node) findFirst(f func(fn *node) (bool, error)) (*node, error) {
	hasFound, err := f(n)
	if err != nil {
		return nil, failure.Annotate(err, "cannot find first matching node")
	}
	if hasFound {
		return n, nil
	}
	for _, child := range n.children {
		found, err := child.findFirst(f)
		if err != nil && !failure.Contains(err, "cannot find") {
			return nil, failure.Annotate(err, "cannot find first matching node")
		}
		if found != nil {
			return found, nil
		}
	}
	return nil, failure.New("cannot find first matching node")
}

// findAll returns all nodes for which the passed function
// returns true.
func (n *node) findAll(f func(fn *node) (bool, error)) ([]*node, error) {
	var allFound []*node
	hasFound, err := f(n)
	if err != nil {
		return nil, failure.Annotate(err, "cannot find all matching nodes")
	}
	if hasFound {
		allFound = append(allFound, n)
	}
	for _, child := range n.children {
		found, err := child.findAll(f)
		if err != nil {
			return nil, failure.Annotate(err, "cannot find all matching nodes")
		}
		if found != nil {
			allFound = append(allFound, found...)
		}
	}
	return allFound, nil
}

// doAll performs the passed function for the node
// and all its children deep to the leafs.
func (n *node) doAll(f func(dn *node) error) error {
	if err := f(n); err != nil {
		return failure.Annotate(err, "cannot perform function on all nodes")
	}
	for _, child := range n.children {
		if err := child.doAll(f); err != nil {
			return failure.Annotate(err, "cannot perform function on all nodes")
		}
	}
	return nil
}

// doChildren performs the passed function for all children.
func (n *node) doChildren(f func(cn *node) error) error {
	for _, child := range n.children {
		if err := f(child); err != nil {
			return failure.Annotate(err, "cannot perform function on all children")
		}
	}
	return nil
}

// size recursively calculates the size of the nodeContainer.
func (n *node) size() int {
	l := 1
	for _, child := range n.children {
		l += child.size()
	}
	return l
}

// deepCopy creates a copy of the node.
func (n *node) deepCopy(c *nodeContainer, p *node) *node {
	cn := &node{
		container: c,
		parent:    p,
		content:   n.content.deepCopy(),
		children:  make([]*node, len(n.children)),
	}
	for i, child := range n.children {
		cn.children[i] = child.deepCopy(c, cn)
	}
	return cn
}

// String implements the Stringer interface.
func (n *node) String() string {
	out := fmt.Sprintf("[%v", n.content)
	if len(n.children) > 0 {
		out += " "
		for _, child := range n.children {
			out += child.String()
		}
	}
	out += "]"
	return out
}

//--------------------
// TREE
//--------------------

// Tree defines the interface for a tree able to store any type
// of values.
type Tree struct {
	container *nodeContainer
}

// NewTree creates a new tree with or without duplicate
// values for children.
func NewTree(v interface{}, duplicates bool) *Tree {
	return &Tree{
		container: newNodeContainer(justValue{v}, duplicates),
	}
}

// At returns the changer of the path defined by the given
// values. If it does not exist it will not be created. Use
// Create() here. So to set a child at a given node path do
//
// err := tree.At("path", 1, "to", "use").Set(12345)
func (t *Tree) At(values ...interface{}) *Changer {
	var path []nodeContent
	for _, value := range values {
		path = append(path, justValue{value})
	}
	n, err := t.container.root.at(path...)
	return &Changer{n, err}
}

// Root returns the top level changer.
func (t *Tree) Root() *Changer {
	return &Changer{t.container.root, nil}
}

// Create returns the changer of the path defined by the
// given keys. If it does not exist it will be created,
// but at least the root key has to be correct.
func (t *Tree) Create(values ...interface{}) *Changer {
	var path []nodeContent
	for _, value := range values {
		path = append(path, justValue{value})
	}
	n, err := t.container.root.create(path...)
	return &Changer{n, err}
}

// FindFirst returns the changer for the first node found
// by the passed function.
func (t *Tree) FindFirst(f func(v interface{}) (bool, error)) *Changer {
	n, err := t.container.root.findFirst(func(fn *node) (bool, error) {
		return f(fn.content.value())
	})
	return &Changer{n, err}
}

// FindAll returns all changers for the nodes found
// by the passed function.
func (t *Tree) FindAll(f func(v interface{}) (bool, error)) []*Changer {
	ns, err := t.container.root.findAll(func(fn *node) (bool, error) {
		return f(fn.content.value())
	})
	if err != nil {
		return []*Changer{{nil, err}}
	}
	var cs []*Changer
	for _, n := range ns {
		cs = append(cs, &Changer{n, nil})
	}
	return cs
}

// DoAll executes the passed function on all nodes.
func (t *Tree) DoAll(f func(v interface{}) error) error {
	return t.container.root.doAll(func(dn *node) error {
		return f(dn.content.value())
	})
}

// DoAllDeep executes the passed function on all nodes
// passing a deep list of values ordered top-down.
func (t *Tree) DoAllDeep(f func(vs []interface{}) error) error {
	return t.container.root.doAll(func(dn *node) error {
		values := []interface{}{}
		cn := dn
		for cn != nil {
			values = append([]interface{}{cn.content.value()}, values...)
			cn = cn.parent
		}
		return f(values)
	})
}

// Len returns the number of nodes of the tree.
func (t *Tree) Len() int {
	return t.container.root.size()
}

// Copy creates a copy of the tree.
func (t *Tree) Copy() *Tree {
	return &Tree{
		container: t.container.deepCopy(),
	}
}

// Deflate cleans the tree with a new root value.
func (t *Tree) Deflate(v interface{}) {
	t.container.root = &node{
		content: justValue{v},
	}
}

// String implements the fmt.Stringer interface.
func (t *Tree) String() string {
	return t.container.root.String()
}

//--------------------
// STRING TREE
//--------------------

// StringTree defines the interface for a tree able to store strings.
type StringTree struct {
	container *nodeContainer
}

// NewStringTree creates a new string tree with or without
// duplicate values for children.
func NewStringTree(v string, duplicates bool) *StringTree {
	return &StringTree{
		container: newNodeContainer(justValue{v}, duplicates),
	}
}

// At returns the changer of the path defined by the given
// values. If it does not exist it will not be created. Use
// Create() here. So to set a child at a given node path do
//
// err := tree.At("path", "one", "to", "use").Set("12345")
func (t *StringTree) At(values ...string) *StringChanger {
	var path []nodeContent
	for _, value := range values {
		path = append(path, justValue{value})
	}
	n, err := t.container.root.at(path...)
	return &StringChanger{n, err}
}

// Root returns the top level changer.
func (t *StringTree) Root() *StringChanger {
	return &StringChanger{t.container.root, nil}
}

// Create returns the changer of the path defined by the
// given keys. If it does not exist it will be created,
// but at least the root key has to be correct.
func (t *StringTree) Create(values ...string) *StringChanger {
	var path []nodeContent
	for _, value := range values {
		path = append(path, justValue{value})
	}
	n, err := t.container.root.create(path...)
	return &StringChanger{n, err}
}

// FindFirst returns the changer for the first node found
// by the passed function.
func (t *StringTree) FindFirst(f func(v string) (bool, error)) *StringChanger {
	n, err := t.container.root.findFirst(func(fn *node) (bool, error) {
		return f(fn.content.value().(string))
	})
	return &StringChanger{n, err}
}

// FindAll returns all changers for the nodes found
// by the passed function.
func (t *StringTree) FindAll(f func(v string) (bool, error)) []*StringChanger {
	ns, err := t.container.root.findAll(func(fn *node) (bool, error) {
		return f(fn.content.value().(string))
	})
	if err != nil {
		return []*StringChanger{{nil, err}}
	}
	var cs []*StringChanger
	for _, n := range ns {
		cs = append(cs, &StringChanger{n, nil})
	}
	return cs
}

// DoAll executes the passed function on all nodes.
func (t *StringTree) DoAll(f func(v string) error) error {
	return t.container.root.doAll(func(dn *node) error {
		return f(dn.content.value().(string))
	})
}

// DoAllDeep executes the passed function on all nodes
// passing a deep list of values ordered top-down.
func (t *StringTree) DoAllDeep(f func(vs []string) error) error {
	return t.container.root.doAll(func(dn *node) error {
		values := []string{}
		cn := dn
		for cn != nil {
			values = append([]string{cn.content.value().(string)}, values...)
			cn = cn.parent
		}
		return f(values)
	})
}

// Len returns the number of nodes of the tree.
func (t *StringTree) Len() int {
	return t.container.root.size()
}

// Copy creates a copy of the tree.
func (t *StringTree) Copy() *StringTree {
	return &StringTree{
		container: t.container.deepCopy(),
	}
}

// Deflate cleans the tree with a new root value.
func (t *StringTree) Deflate(v string) {
	t.container.root = &node{
		content: justValue{v},
	}
}

// String implements the fmt.Stringer interface.
func (t *StringTree) String() string {
	return t.container.root.String()
}

//--------------------
// KEY/VALUE TREE
//--------------------

// KeyValueTree defines the interface for a tree able to store key/value pairs.
type KeyValueTree struct {
	container *nodeContainer
}

// NewKeyValueTree creates a new key/value tree with or without
// duplicate values for children.
func NewKeyValueTree(k string, v interface{}, duplicates bool) *KeyValueTree {
	return &KeyValueTree{
		container: newNodeContainer(keyValue{k, v}, duplicates),
	}
}

// At returns the changer of the path defined by the given
// values. If it does not exist it will not be created. Use
// Create() here. So to set a child at a given node path do
//
// err := tree.At("path", "one", "to", "use").Set(12345)
func (t *KeyValueTree) At(keys ...string) *KeyValueChanger {
	var path []nodeContent
	for _, key := range keys {
		path = append(path, keyValue{key, nil})
	}
	n, err := t.container.root.at(path...)
	return &KeyValueChanger{n, err}
}

// Root returns the top level changer.
func (t *KeyValueTree) Root() *KeyValueChanger {
	return &KeyValueChanger{t.container.root, nil}
}

// Create returns the changer of the path defined by the
// given keys. If it does not exist it will be created,
// but at least the root key has to be correct.
func (t *KeyValueTree) Create(keys ...string) *KeyValueChanger {
	var path []nodeContent
	for _, key := range keys {
		path = append(path, keyValue{key, nil})
	}
	n, err := t.container.root.create(path...)
	return &KeyValueChanger{n, err}
}

// FindFirst returns the changer for the first node found
// by the passed function.
func (t *KeyValueTree) FindFirst(f func(k string, v interface{}) (bool, error)) *KeyValueChanger {
	n, err := t.container.root.findFirst(func(fn *node) (bool, error) {
		return f(fn.content.key().(string), fn.content.value())
	})
	return &KeyValueChanger{n, err}
}

// FindAll returns all changers for the nodes found
// by the passed function.
func (t *KeyValueTree) FindAll(f func(k string, v interface{}) (bool, error)) []*KeyValueChanger {
	ns, err := t.container.root.findAll(func(fn *node) (bool, error) {
		return f(fn.content.key().(string), fn.content.value())
	})
	if err != nil {
		return []*KeyValueChanger{{nil, err}}
	}
	var cs []*KeyValueChanger
	for _, n := range ns {
		cs = append(cs, &KeyValueChanger{n, nil})
	}
	return cs
}

// DoAll executes the passed function on all nodes.
func (t *KeyValueTree) DoAll(f func(k string, v interface{}) error) error {
	return t.container.root.doAll(func(dn *node) error {
		return f(dn.content.key().(string), dn.content.value())
	})
}

// DoAllDeep executes the passed function on all nodes
// passing a deep list of keys ordered top-down.
func (t *KeyValueTree) DoAllDeep(f func(ks []string, v interface{}) error) error {
	return t.container.root.doAll(func(dn *node) error {
		keys := []string{}
		cn := dn
		for cn != nil {
			keys = append([]string{cn.content.key().(string)}, keys...)
			cn = cn.parent
		}
		return f(keys, dn.content.value())
	})
}

// Len returns the number of nodes of the tree.
func (t *KeyValueTree) Len() int {
	return t.container.root.size()
}

// Copy creates a copy of the tree.
func (t *KeyValueTree) Copy() *KeyValueTree {
	return &KeyValueTree{
		container: t.container.deepCopy(),
	}
}

// CopyAt creates a copy of a subtree.
func (t *KeyValueTree) CopyAt(keys ...string) (*KeyValueTree, error) {
	var path []nodeContent
	for _, key := range keys {
		path = append(path, keyValue{key, ""})
	}
	n, err := t.container.root.at(path...)
	if err != nil {
		return nil, err
	}
	nc := &nodeContainer{
		duplicates: t.container.duplicates,
	}
	nc.root = n.deepCopy(nc, nil)
	return &KeyValueTree{nc}, nil
}

// Deflate cleans the tree with a new root value.
func (t *KeyValueTree) Deflate(k string, v interface{}) {
	t.container.root = &node{
		content: keyValue{k, v},
	}
}

// String implements the fmt.Stringer interface.
func (t *KeyValueTree) String() string {
	return t.container.root.String()
}

//--------------------
// KEY/STRING VALUE TREE
//--------------------

// KeyStringValueTree defines the interface for a tree able to store
// key/string value pairs.
type KeyStringValueTree struct {
	container *nodeContainer
}

// NewKeyStringValueTree creates a new key/value tree with or without
// duplicate values for children and strings as values.
func NewKeyStringValueTree(k, v string, duplicates bool) *KeyStringValueTree {
	return &KeyStringValueTree{
		container: newNodeContainer(keyValue{k, v}, duplicates),
	}
}

// At returns the changer of the path defined by the given
// values. If it does not exist it will not be created. Use
// Create() here. So to set a child at a given node path do
//
// err := tree.At("path", "one", "to", "use").Set(12345)
func (t *KeyStringValueTree) At(keys ...string) *KeyStringValueChanger {
	var path []nodeContent
	for _, key := range keys {
		path = append(path, keyValue{key, ""})
	}
	n, err := t.container.root.at(path...)
	return &KeyStringValueChanger{n, err}
}

// Root returns the top level changer.
func (t *KeyStringValueTree) Root() *KeyStringValueChanger {
	return &KeyStringValueChanger{t.container.root, nil}
}

// Create returns the changer of the path defined by the
// given keys. If it does not exist it will be created,
// but at least the root key has to be correct.
func (t *KeyStringValueTree) Create(keys ...string) *KeyStringValueChanger {
	var path []nodeContent
	for _, key := range keys {
		path = append(path, keyValue{key, ""})
	}
	n, err := t.container.root.create(path...)
	return &KeyStringValueChanger{n, err}
}

// FindFirst returns the changer for the first node found
// by the passed function.
func (t *KeyStringValueTree) FindFirst(f func(k, v string) (bool, error)) *KeyStringValueChanger {
	n, err := t.container.root.findFirst(func(fn *node) (bool, error) {
		return f(fn.content.key().(string), fn.content.value().(string))
	})
	return &KeyStringValueChanger{n, err}
}

// FindAll returns all changers for the nodes found
// by the passed function.
func (t *KeyStringValueTree) FindAll(f func(k, v string) (bool, error)) []*KeyStringValueChanger {
	ns, err := t.container.root.findAll(func(fn *node) (bool, error) {
		return f(fn.content.key().(string), fn.content.value().(string))
	})
	if err != nil {
		return []*KeyStringValueChanger{{nil, err}}
	}
	var cs []*KeyStringValueChanger
	for _, n := range ns {
		cs = append(cs, &KeyStringValueChanger{n, nil})
	}
	return cs
}

// DoAll executes the passed function on all nodes.
func (t *KeyStringValueTree) DoAll(f func(k, v string) error) error {
	return t.container.root.doAll(func(dn *node) error {
		return f(dn.content.key().(string), dn.content.value().(string))
	})
}

// DoAllDeep executes the passed function on all nodes
// passing a deep list of keys ordered top-down.
func (t *KeyStringValueTree) DoAllDeep(f func(ks []string, v string) error) error {
	return t.container.root.doAll(func(dn *node) error {
		keys := []string{}
		cn := dn
		for cn != nil {
			keys = append([]string{cn.content.key().(string)}, keys...)
			cn = cn.parent
		}
		return f(keys, dn.content.value().(string))
	})
}

// Len returns the number of nodes of the tree.
func (t *KeyStringValueTree) Len() int {
	return t.container.root.size()
}

// Copy creates a copy of the tree.
func (t *KeyStringValueTree) Copy() *KeyStringValueTree {
	return &KeyStringValueTree{
		container: t.container.deepCopy(),
	}
}

// CopyAt creates a copy of a subtree.
func (t *KeyStringValueTree) CopyAt(keys ...string) (*KeyStringValueTree, error) {
	var path []nodeContent
	for _, key := range keys {
		path = append(path, keyValue{key, ""})
	}
	n, err := t.container.root.at(path...)
	if err != nil {
		return nil, err
	}
	nc := &nodeContainer{
		duplicates: t.container.duplicates,
	}
	nc.root = n.deepCopy(nc, nil)
	return &KeyStringValueTree{nc}, nil
}

// Deflate cleans the tree with a new root value.
func (t *KeyStringValueTree) Deflate(k, v string) {
	t.container.root = &node{
		content: keyValue{k, v},
	}
}

// String implements the fmt.Stringer interface.
func (t *KeyStringValueTree) String() string {
	return t.container.root.String()
}

// EOF
