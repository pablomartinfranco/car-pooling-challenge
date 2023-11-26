package array

import "sync"

type Array[T any] struct {
	slice []T
	mutex sync.Mutex
}

func New[T any](size int) *Array[T] {
	return &Array[T]{
		slice: make([]T, size),
		mutex: sync.Mutex{},
	}
}

func (a *Array[T]) Get(index int) (value T) {
	a.mutex.Lock()
	value = a.slice[index]
	a.mutex.Unlock()
	return value
}

func (a *Array[T]) TryGet(index int) (value T, ok bool) {
	a.mutex.Lock()
	if index >= len(a.slice) {
		a.mutex.Unlock()
		return value, false
	}
	value = a.slice[index]
	a.mutex.Unlock()
	return value, true
}

func (a *Array[T]) Set(index int, value *T) {
	a.mutex.Lock()
	a.slice[index] = *value
	a.mutex.Unlock()
}

func (a *Array[T]) TrySet(index int, value *T) bool {
	a.mutex.Lock()
	if index >= len(a.slice) {
		a.mutex.Unlock()
		return false
	}
	a.slice[index] = *value
	a.mutex.Unlock()
	return true
}

func (a *Array[T]) Remove(index int) (value T, ok bool) {
	a.mutex.Lock()
	if index >= len(a.slice) {
		a.mutex.Unlock()
		return value, false
	}
	value = a.slice[index]
	a.slice = append(a.slice[:index], a.slice[index+1:]...)
	a.mutex.Unlock()
	return value, true
}

func (a *Array[T]) Delete(index int) bool {
	a.mutex.Lock()
	if index >= len(a.slice) {
		a.mutex.Unlock()
		return false
	}
	a.slice = append(a.slice[:index], a.slice[index+1:]...)
	a.mutex.Unlock()
	return true
}

func (a *Array[T]) Insert(value *T) {
	a.mutex.Lock()
	a.slice = append(a.slice, *value)
	a.mutex.Unlock()
}

func (a *Array[T]) Append(values []T) {
	a.mutex.Lock()
	a.slice = append(a.slice, values...)
	a.mutex.Unlock()
}

func (a *Array[T]) Prepend(values []T) {
	a.mutex.Lock()
	a.slice = append(values, a.slice...)
	a.mutex.Unlock()
}

func (a *Array[T]) Take() (value T, ok bool) {
	a.mutex.Lock()
	if len(a.slice) == 0 {
		a.mutex.Unlock()
		return value, false
	}
	value = a.slice[0]
	a.slice = a.slice[1:]
	a.mutex.Unlock()
	return value, true
}

func (a *Array[T]) Pop() (value T, ok bool) {
	a.mutex.Lock()
	if len(a.slice) == 0 {
		a.mutex.Unlock()
		return value, false
	}
	value = a.slice[len(a.slice)-1]
	a.slice = a.slice[:len(a.slice)-1]
	a.mutex.Unlock()
	return value, true
}

func (a *Array[T]) Any(predicate func(int, T) bool) bool {
	a.mutex.Lock()
	for index, value := range a.slice {
		if predicate(index, value) {
			a.mutex.Unlock()
			return true
		}
	}
	a.mutex.Unlock()
	return false
}

func (a *Array[T]) Len() int {
	a.mutex.Lock()
	length := len(a.slice)
	a.mutex.Unlock()
	return length
}

func (a *Array[T]) Clear() {
	a.mutex.Lock()
	a.slice = make([]T, 0)
	a.mutex.Unlock()
}

func (a *Array[T]) ForEach(fn func(int, T)) {
	a.mutex.Lock()
	for index, value := range a.slice {
		fn(index, value)
	}
	a.mutex.Unlock()
}
