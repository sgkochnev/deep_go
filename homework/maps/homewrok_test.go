package main

import (
	"cmp"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type node[K cmp.Ordered, V any] struct {
	value       V
	key         K
	left, right *node[K, V]
	isDel       bool
}

type OrderedMap[K cmp.Ordered, V any] struct {
	head *node[K, V]
	len  int
}

func NewOrderedMap[K cmp.Ordered, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	n := &node[K, V]{
		key:   key,
		value: value,
	}

	defer func() { m.len++ }()

	if m.len == 0 {
		m.head = n
		return
	}

	head := m.head
	for {
		switch {
		case head.key > key:
			if head.left == nil {
				head.left = n
				return
			}
			head = head.left

		case head.key < key:
			if head.right == nil {
				head.right = n
				return
			}
			head = head.right

		default:
			if head.isDel {
				head.value = n.value
			}
			return
		}
	}
}

func (m *OrderedMap[K, V]) Erase(key K) {
	head := m.head
	for {
		switch {
		case head != nil && head.key > key:
			head = head.left
		case head != nil && head.key < key:
			head = head.right
		default:
			if head == nil {
				return
			}

			if head.left == nil && head.right == nil {
				head = nil
			} else {
				head.isDel = true
			}

			m.len--

			return
		}
	}
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	head := m.head
	for {
		switch {
		case head != nil && head.key > key:
			head = head.left
		case head != nil && head.key < key:
			head = head.right
		default:
			if head == nil {
				return false
			}

			return !head.isDel && head.key == key
		}
	}
}

func (m *OrderedMap[K, V]) Size() int {
	return m.len
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	// if m == nil {
	// 	return
	// }

	// action(m.head.key, m.head.value)

	// m.head.left.ForEach(action)
	// m.head.right.ForEach(action)

}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
