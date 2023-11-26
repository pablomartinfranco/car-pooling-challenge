package domain

import (
	"car-pooling-challenge/pkg/array"
	"car-pooling-challenge/pkg/config"
	"car-pooling-challenge/pkg/event"
	"car-pooling-challenge/pkg/index"
	"car-pooling-challenge/pkg/logger"
	"car-pooling-challenge/pkg/queue"
	"context"
	"log"
	"sync"
)

const (
	MinSeats = 4
	MaxSeats = 6
	MinGroup = 1
	MaxGroup = 6
)

type Pooling struct {
	logger       *log.Logger
	config       *config.Config
	context      context.Context
	cancel       context.CancelFunc
	dropoffMutex *sync.Mutex
	journeyMutex *sync.Mutex
	workersMutex *sync.Mutex
	dropoffEvent event.Event
	journeyEvent event.Event
	dropoffRetry *index.Index[int]
	waitingQueue *queue.Queue[*Group]
	dropoffQueue *queue.Queue[*Group]
	journeyIndex *index.Index[*Journey]
	freeSeatsIdx *array.Array[*index.Index[*Car]]
}

func NewPooling(cfg *config.Config) *Pooling {

	context, cancel := context.WithCancel(context.Background())

	var service = &Pooling{
		logger:       logger.NewLogger("Pooling"),
		config:       cfg,
		context:      context,
		cancel:       cancel,
		dropoffMutex: &sync.Mutex{},
		journeyMutex: &sync.Mutex{},
		workersMutex: &sync.Mutex{},
		dropoffEvent: event.New(),
		journeyEvent: event.New(),
		dropoffRetry: index.New[int](),
		waitingQueue: queue.New[*Group](),
		dropoffQueue: queue.New[*Group](),
		journeyIndex: index.New[*Journey](),
		freeSeatsIdx: array.New[*index.Index[*Car]](MaxSeats + 1),
	}

	for i := 0; i <= MaxSeats; i++ {
		var cars = index.New[*Car]()
		service.freeSeatsIdx.TrySet(i, &cars)
	}

	return service
}

func (p *Pooling) Reset() {
	p.waitingQueue.Clear()
	p.dropoffQueue.Clear()
	p.journeyIndex.Clear()
	var clear = func(i int, v *index.Index[*Car]) { v.Clear() }
	p.freeSeatsIdx.ForEach(clear)
}

func (p *Pooling) Run() {
	// start the journey worker pool in goroutines
	for i := 0; i < p.config.JourneyWorkerPool; i++ {
		go p.journeyWorker(i)
	}

	// start the journey worker pool in goroutines
	for i := 0; i < p.config.DropoffWorkerPool; i++ {
		go p.dropoffWorker(i)
	}
}
