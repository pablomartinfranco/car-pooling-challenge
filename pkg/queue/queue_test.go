package queue_test

import (
	"car-pooling-challenge/pkg/queue"
	"testing"
	"time"
)

func Test_EnqueueDequeue(t *testing.T) {
	type args struct {
		item int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "When Queue item Then Dequeue same item",
			args: args{
				item: 7,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := queue.New[int]()
			q.Enqueue(&tt.args.item)

			if q.Dequeue() != tt.args.item {
				t.Errorf("Queue() = %v, want %v", q.Dequeue(), tt.args.item)
			}
		})
	}
}

func Test_EnqueueDequeueThreaded(t *testing.T) {
	type args struct {
		item int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "When Queue item in principal thread Then Dequeue same item in worker thread",
			args: args{
				item: 7,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var aQueue = queue.New[int]()
			var number = new(int)
			var callback = func(n int) {
				*number = n
			}
			var worker = func(q *queue.Queue[int], callback func(int)) {
				item := q.Dequeue()
				callback(item)
			}

			go worker(aQueue, callback)

			time.Sleep(100 * time.Microsecond)

			aQueue.Enqueue(&tt.args.item)

			time.Sleep(100 * time.Microsecond)

			if *number != tt.args.item {
				t.Errorf("Queue() = %v, want %v", *number, tt.args.item)
			}
		})
	}
}
