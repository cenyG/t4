package collections

import (
	"github.com/pkg/errors"
	"sync"
)

type ConcurrentArray[T any] interface {
	Append(value T)
	Get(index int) (T, error)
	Set(index int, values T) error
	SetAll(value []T)
	Length() int
	Remove(index int) bool
	Values() []T
}

// concurrentArray представляет потокобезопасный массив с использованием обобщений.
type concurrentArray[T any] struct {
	mu  sync.RWMutex
	arr []T
}

// NewConcurrentArray создает новый потокобезопасный массив.
func NewConcurrentArray[T any](size int) ConcurrentArray[T] {
	return &concurrentArray[T]{
		arr: make([]T, size),
	}
}

// Append добавляет элемент в конец массива.
func (c *concurrentArray[T]) Append(value T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.arr = append(c.arr, value)
}

// Get - get element by index
func (c *concurrentArray[T]) Get(index int) (T, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if index < 0 || index >= len(c.arr) {
		var zeroValue T
		return zeroValue, errors.New("index out of range")
	}
	return c.arr[index], nil
}

// Set - set element by index
func (c *concurrentArray[T]) Set(index int, value T) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if index < 0 || index >= len(c.arr) {
		return errors.New("index out of range")
	}
	c.arr[index] = value

	return nil
}

// SetAll - rewrite array with new values
func (c *concurrentArray[T]) SetAll(values []T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.arr = values
}

// Length возвращает длину массива.
func (c *concurrentArray[T]) Length() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.arr)
}

// Remove удаляет элемент по индексу.
func (c *concurrentArray[T]) Remove(index int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if index < 0 || index >= len(c.arr) {
		return false
	}
	c.arr = append(c.arr[:index], c.arr[index+1:]...)
	return true
}

// Values возвращает копию всех элементов массива.
func (c *concurrentArray[T]) Values() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	copyArr := make([]T, len(c.arr))
	copy(copyArr, c.arr)
	return copyArr
}
