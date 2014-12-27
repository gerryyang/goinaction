// Copyright © 2011-12 Qtrac Ltd.
// 
// This program or package and any associated files are licensed under the
// Apache License, Version 2.0 (the "License"); you may not use these files
// except in compliance with the License. You can get a copy of the License
// at: http://www.apache.org/licenses/LICENSE-2.0.
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The data structure is a left-leaning red-black tree based on Robert
// Sedgewick's implementations as described in
// http://www.cs.princeton.edu/~rs/talks/LLRB/LLRB.pdf and
// http://www.cs.princeton.edu/~rs/talks/LLRB/RedBlack.pdf, with some of
// the code based on Lee Stanza's C++ code at
// http://www.teachsolaisgames.com/articles/balanced_left_leaning.html
// Thanks also to Russ Cox for many useful suggestions.

// Package omap implements an efficient key-ordered map.
//
// Keys and values may be of any type, but all keys must be comparable
// using the less than function that is passed in to the omap.New()
// function, or the less than function provided by the omap.New*()
// construction functions.
package omap

import "strings"

// NewStringKeyed returns an empty Map that accepts case-sensitive
// string keys.
func NewStringKeyed() *Map {
    return &Map{less: func(a, b interface{}) bool {
        return a.(string) < b.(string)
    }}
}

// NewCaseFoldedKeyed returns an empty Map that accepts case-insensitive
// string keys.
func NewCaseFoldedKeyed() *Map {
    return &Map{less: func(a, b interface{}) bool {
        return strings.ToLower(a.(string)) < strings.ToLower(b.(string))
    }}
}

// NewIntKeyed returns an empty Map that accepts int keys.
func NewIntKeyed() *Map {
    return &Map{less: func(a, b interface{}) bool {
        return a.(int) < b.(int)
    }}
}

// NewFloat64Keyed returns an empty Map that accepts float64 keys.
func NewFloat64Keyed() *Map {
    return &Map{less: func(a, b interface{}) bool {
        return a.(float64) < b.(float64)
    }}
}

// New returns an empty Map that uses the given less than function to
// compare keys. For example:
//      type Point { X, Y int }
//      pointMap := omap.New(func(a, b interface{}) bool {
//              α, β := a.(Point), b.(Point)
//              if α.X != β.X {
//                  return α.X < β.X
//              }
//              return α.Y < β.Y
//          })
func New(less func(interface{}, interface{}) bool) *Map {
    return &Map{less: less}
}

// Map is a key-ordered map.
// The zero value is an invalid map! Use one of the construction functions
// (e.g., New()), to create a map for a specific key type.
type Map struct {
    root   *node
    less   func(interface{}, interface{}) bool
    length int
}

type node struct {
    key, value  interface{}
    red         bool
    left, right *node
}

// Insert inserts a new key-value into the Map and returns true; or
// replaces an existing key-value pair's value if the keys are equal and
// returns false. For example:
//      inserted := myMap.Insert(key, value).
func (m *Map) Insert(key, value interface{}) (inserted bool) {
    m.root, inserted = m.insert(m.root, key, value)
    m.root.red = false
    if inserted {
        m.length++
    }
    return inserted
}

// Find returns the value and true if the key is in the Map or nil and
// false otherwise. For example:
//      value, found := myMap.Find(key).
func (m *Map) Find(key interface{}) (value interface{}, found bool) {
    root := m.root
    for root != nil {
        if m.less(key, root.key) {
            root = root.left
        } else if m.less(root.key, key) {
            root = root.right
        } else {
            return root.value, true
        }
    }
    return nil, false
}

// Delete deletes the key-value with the given key from the Map and returns
// true, or does nothing and returns false if there is no key-value with
// the given key. For example:
//      deleted := myMap.Delete(key).
func (m *Map) Delete(key interface{}) (deleted bool) {
    if m.root != nil {
        if m.root, deleted = m.remove(m.root, key); m.root != nil {
            m.root.red = false
        }
    }
    if deleted {
        m.length--
    }
    return deleted
}

// Do calls the given function on every key-value in the Map in order.
func (m *Map) Do(function func(interface{}, interface{})) {
    do(m.root, function)
}

// Len returns the number of key-value pairs in the map.
func (m *Map) Len() int {
    return m.length
}

func (m *Map) insert(root *node, key, value interface{}) (*node, bool) {
    inserted := false
    if root == nil { // If the key was in the tree it would belong here
        return &node{key: key, value: value, red: true}, true
    }
    if isRed(root.left) && isRed(root.right) {
        colorFlip(root)
    }
    if m.less(key, root.key) {
        root.left, inserted = m.insert(root.left, key, value)
    } else if m.less(root.key, key) {
        root.right, inserted = m.insert(root.right, key, value)
    } else { // The key is already in the tree so just replace its value
        root.value = value
    }
    if isRed(root.right) && !isRed(root.left) {
        root = rotateLeft(root)
    }
    if isRed(root.left) && isRed(root.left.left) {
        root = rotateRight(root)
    }
    return root, inserted
}

func isRed(root *node) bool { return root != nil && root.red }

func colorFlip(root *node) {
    root.red = !root.red
    if root.left != nil {
        root.left.red = !root.left.red
    }
    if root.right != nil {
        root.right.red = !root.right.red
    }
}

func rotateLeft(root *node) *node {
    x := root.right
    root.right = x.left
    x.left = root
    x.red = root.red
    root.red = true
    return x
}

func rotateRight(root *node) *node {
    x := root.left
    root.left = x.right
    x.right = root
    x.red = root.red
    root.red = true
    return x
}

func do(root *node, function func(interface{}, interface{})) {
    if root != nil {
        do(root.left, function)
        function(root.key, root.value)
        do(root.right, function)
    }
}

// We do not provide an exported First() method because this is an
// implementation detail.
func first(root *node) *node {
    for root.left != nil {
        root = root.left
    }
    return root
}

func (m *Map) remove(root *node, key interface{}) (*node, bool) {
    deleted := false
    if m.less(key, root.key) {
        if root.left != nil {
            if !isRed(root.left) && !isRed(root.left.left) {
                root = moveRedLeft(root)
            }
            root.left, deleted = m.remove(root.left, key)
        }
    } else {
        if isRed(root.left) {
            root = rotateRight(root)
        }
        if !m.less(key, root.key) && !m.less(root.key, key) &&
            root.right == nil {
            return nil, true
        }
        if root.right != nil {
            if !isRed(root.right) && !isRed(root.right.left) {
                root = moveRedRight(root)
            }
            if !m.less(key, root.key) && !m.less(root.key, key) {
                smallest := first(root.right)
                root.key = smallest.key
                root.value = smallest.value
                root.right = deleteMinimum(root.right)
                deleted = true
            } else {
                root.right, deleted = m.remove(root.right, key)
            }
        }
    }
    return fixUp(root), deleted
}

func moveRedLeft(root *node) *node {
    colorFlip(root)
    if root.right != nil && isRed(root.right.left) {
        root.right = rotateRight(root.right)
        root = rotateLeft(root)
        colorFlip(root)
    }
    return root
}

func moveRedRight(root *node) *node {
    colorFlip(root)
    if root.left != nil && isRed(root.left.left) {
        root = rotateRight(root)
        colorFlip(root)
    }
    return root
}

func deleteMinimum(root *node) *node {
    if root.left == nil {
        return nil
    }
    if !isRed(root.left) && !isRed(root.left.left) {
        root = moveRedLeft(root)
    }
    root.left = deleteMinimum(root.left)
    return fixUp(root)
}

func fixUp(root *node) *node {
    if isRed(root.right) {
        root = rotateLeft(root)
    }
    if isRed(root.left) && isRed(root.left.left) {
        root = rotateRight(root)
    }
    if isRed(root.left) && isRed(root.right) {
        colorFlip(root)
    }
    return root
}
