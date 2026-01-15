package memtable

import (
	"math/rand"
	"time"
)

const (
	maxLevel    = 12
	pProbability = 0.5
)

// Node represents a single element in the skip list.
type Node struct {
	key     []byte
	value   []byte
	forward []*Node
}

func newNode(key, value []byte, level int) *Node {
	return &Node{
		key:     key,
		value:   value,
		forward: make([]*Node, level),
	}
}

// SkipList is a probabilistic data structure that provides O(log n) search and insertion.
type SkipList struct {
	head  *Node
	level int // current maximum level of the skip list
	rnd   *rand.Rand
}

// NewSkipList creates a new SkipList.
func NewSkipList() *SkipList {
	return &SkipList{
		head:  newNode(nil, nil, maxLevel),
		level: 1,
		rnd:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// randomLevel generates a random level for a new node.
func (sl *SkipList) randomLevel() int {
	level := 1
	for sl.rnd.Float64() < pProbability && level < maxLevel {
		level++
	}
	return level
}

// Get finds the value associated with a key.
func (sl *SkipList) Get(key []byte) ([]byte, bool) {
	curr := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil && compare(curr.forward[i].key, key) < 0 {
			curr = curr.forward[i]
		}
	}

	curr = curr.forward[0]
	if curr != nil && compare(curr.key, key) == 0 {
		return curr.value, true
	}

	return nil, false
}

// Put inserts or updates a key-value pair.
func (sl *SkipList) Put(key, value []byte) {
	update := make([]*Node, maxLevel)
	curr := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil && compare(curr.forward[i].key, key) < 0 {
			curr = curr.forward[i]
		}
		update[i] = curr
	}

	curr = curr.forward[0]

	if curr != nil && compare(curr.key, key) == 0 {
		// Update existing node
		curr.value = value
		return
	}

	// Insert new node
	newLevel := sl.randomLevel()
	if newLevel > sl.level {
		for i := sl.level; i < newLevel; i++ {
			update[i] = sl.head
		}
		sl.level = newLevel
	}

	node := newNode(key, value, newLevel)
	for i := 0; i < newLevel; i++ {
		node.forward[i] = update[i].forward[i]
		update[i].forward[i] = node
	}
}

// Iterator provides a way to traverse the skip list in order.
type Iterator struct {
	sl   *SkipList
	curr *Node
}

// NewIterator creates a new iterator for the skip list.
func (sl *SkipList) NewIterator() *Iterator {
	return &Iterator{
		sl:   sl,
		curr: nil,
	}
}

// Next moves the iterator to the next element.
func (it *Iterator) Next() {
	if it.curr == nil {
		it.curr = it.sl.head.forward[0]
	} else {
		it.curr = it.curr.forward[0]
	}
}

// Valid returns true if the iterator is positioned at a valid node.
func (it *Iterator) Valid() bool {
	return it.curr != nil
}

// Key returns the key of the node at the current position.
func (it *Iterator) Key() []byte {
	return it.curr.key
}

// Value returns the value of the node at the current position.
func (it *Iterator) Value() []byte {
	return it.curr.value
}

// Seek moves the iterator to the first node whose key is >= target.
func (it *Iterator) Seek(target []byte) {
	curr := it.sl.head
	for i := it.sl.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil && compare(curr.forward[i].key, target) < 0 {
			curr = curr.forward[i]
		}
	}
	it.curr = curr.forward[0]
}

// SeekToFirst moves the iterator to the first node.
func (it *Iterator) SeekToFirst() {
	it.curr = it.sl.head.forward[0]
}

// compare is a helper to compare byte slices.
// Returns -1 if a < b, 1 if a > b, 0 if a == b.
func compare(a, b []byte) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	if len(a) < len(b) {
		return -1
	} else if len(a) > len(b) {
		return 1
	}
	return 0
}
