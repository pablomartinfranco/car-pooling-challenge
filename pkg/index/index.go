package index

import "sync"

type Index[T any] struct {
	data  map[int]T
	mutex sync.Mutex
}

func New[T any]() *Index[T] {
	return &Index[T]{
		data:  make(map[int]T),
		mutex: sync.Mutex{},
	}
}

func (idx *Index[T]) Insert(key int, value *T) {
	idx.mutex.Lock()
	idx.data[key] = *value
	idx.mutex.Unlock()
}

func (idx *Index[T]) Update(key int, value *T) {
	idx.mutex.Lock()
	idx.data[key] = *value
	idx.mutex.Unlock()
}

func (idx *Index[T]) Delete(key int) {
	idx.mutex.Lock()
	delete(idx.data, key)
	idx.mutex.Unlock()
}

func (idx *Index[T]) Remove(key int) (value T, ok bool) {
	idx.mutex.Lock()
	value, ok = idx.data[key]
	delete(idx.data, key)
	idx.mutex.Unlock()
	return value, ok
}

func (idx *Index[T]) TakeOne() (value T, ok bool) {
	idx.mutex.Lock()
	for key, value := range idx.data {
		delete(idx.data, key)
		idx.mutex.Unlock()
		return value, true
	}
	idx.mutex.Unlock()
	return value, false
}

func (idx *Index[T]) TakeAll() (values []T) {
	idx.mutex.Lock()
	for key, value := range idx.data {
		delete(idx.data, key)
		values = append(values, value)
	}
	idx.mutex.Unlock()
	return values
}

func (idx *Index[T]) Lookup(key int) (value T, ok bool) {
	idx.mutex.Lock()
	value, ok = idx.data[key]
	idx.mutex.Unlock()
	return value, ok
}

func (idx *Index[T]) Size() (size int) {
	idx.mutex.Lock()
	size = len(idx.data)
	idx.mutex.Unlock()
	return size
}

func (idx *Index[T]) Clear() {
	idx.mutex.Lock()
	idx.data = make(map[int]T)
	idx.mutex.Unlock()
}

func (idx *Index[T]) HasKey(key int) bool {
	idx.mutex.Lock()
	_, ok := idx.data[key]
	idx.mutex.Unlock()
	return ok
}

func (idx *Index[T]) Keys() []int {
	idx.mutex.Lock()
	keys := make([]int, 0, len(idx.data))
	for key := range idx.data {
		keys = append(keys, key)
	}
	idx.mutex.Unlock()
	return keys
}

func (idx *Index[T]) LastKey() (key int) {
	idx.mutex.Lock()
	for k := range idx.data {
		if k > key {
			key = k
		}
	}
	idx.mutex.Unlock()
	return key
}

func (idx *Index[T]) Filter(fn func(T) bool) []T {
	idx.mutex.Lock()
	values := make([]T, 0, len(idx.data))
	for _, value := range idx.data {
		if fn(value) {
			values = append(values, value)
		}
	}
	idx.mutex.Unlock()
	return values
}

func (idx *Index[T]) Map(fn func(T) T) {
	idx.mutex.Lock()
	for key, value := range idx.data {
		idx.data[key] = fn(value)
	}
	idx.mutex.Unlock()
}

func (idx *Index[T]) Reduce(fn func(T, T) T) T {
	idx.mutex.Lock()
	var result T
	for _, value := range idx.data {
		result = fn(result, value)
	}
	idx.mutex.Unlock()
	return result
}

func (idx *Index[T]) ForEach(fn func(int, T)) {
	idx.mutex.Lock()
	for key, value := range idx.data {
		fn(key, value)
	}
	idx.mutex.Unlock()
}

func (idx *Index[T]) Any(fn func(int, T) bool) bool {
	idx.mutex.Lock()
	for key, value := range idx.data {
		if fn(key, value) {
			idx.mutex.Unlock()
			return true
		}
	}
	idx.mutex.Unlock()
	return false
}

func (idx *Index[T]) All(fn func(int, T) bool) bool {
	idx.mutex.Lock()
	for key, value := range idx.data {
		if !fn(key, value) {
			idx.mutex.Unlock()
			return false
		}
	}
	idx.mutex.Unlock()
	return true
}

func (idx *Index[T]) Count(fn func(int, T) bool) int {
	idx.mutex.Lock()
	count := 0
	for key, value := range idx.data {
		if fn(key, value) {
			count++
		}
	}
	idx.mutex.Unlock()
	return count
}
