package queue

import (
	"car-pooling-challenge/pkg/index"
	"sync"
)

type Queue[T any] struct {
	head     int
	tail     int
	data     *index.Index[T]
	mutex    *sync.Mutex
	notEmpty *sync.Cond
}

func New[T any]() *Queue[T] {
	mutex := sync.Mutex{}
	return &Queue[T]{
		head:     1,
		tail:     0,
		data:     index.New[T](),
		mutex:    &mutex,
		notEmpty: sync.NewCond(&mutex),
	}
}

func (q *Queue[T]) Enqueue(item *T) {
	q.mutex.Lock()
	q.tail += 1
	q.data.Insert(q.tail, item)
	q.mutex.Unlock()
	q.notEmpty.Signal()
}

func (q *Queue[T]) Dequeue() (item T) {
	q.mutex.Lock()
	for q.data.Size() == 0 {
		q.notEmpty.Wait()
	}
	item, _ = q.data.Remove(q.head)
	q.head += 1
	q.mutex.Unlock()
	return item
}

func (q *Queue[T]) DequeueAll() (items []T) {
	q.mutex.Lock()
	for q.data.Size() == 0 {
		q.notEmpty.Wait()
	}
	items = q.data.TakeAll()
	q.head = 1
	q.tail = 0
	q.mutex.Unlock()
	return items
}

func (q *Queue[T]) Any(fn func(T) bool) (ok bool) {
	q.mutex.Lock()
	ok = q.data.Any(func(_ int, value T) bool {
		return fn(value)
	})
	q.mutex.Unlock()
	return ok
}

func (q *Queue[T]) Size() int {
	return q.data.Size()
}

func (q *Queue[T]) Clear() {
	q.mutex.Lock()
	q.data.Clear()
	q.head = 1
	q.tail = 0
	q.mutex.Unlock()
}

func (q *Queue[T]) ForEach(fn func(T)) {
	q.mutex.Lock()
	q.data.ForEach(func(_ int, value T) {
		fn(value)
	})
	q.mutex.Unlock()
}
