package domain

import "strconv"

func (p *Pooling) journeyWorker(n int) {
	var worker = "Journey worker " + strconv.Itoa(n)
	p.logger.Printf("[%s] started", worker)
	defer p.logger.Printf("[%s] stopped", worker)

	for {
		p.DebugStates(worker)

		p.journeyEvent.Wait()

		var groups = p.waitingQueue.DequeueAll()

		for _, group := range groups {

			if ok := p.tryJourney(group); !ok {

				p.waitingQueue.Enqueue(&group)
			}
		}

		if p.context.Err() != nil {
			return
		}
	}
}

func (p *Pooling) tryJourney(group *Group) bool {
	p.workersMutex.Lock()
	car, freeSeats, ok := p.takeCarAvailable(group.People)
	if !ok {
		p.workersMutex.Unlock()
		return false
	}
	var journey = &Journey{group, car}
	p.journeyIndex.Insert(group.Id, &journey)
	freeSeats -= group.People
	p.freeSeatsIdx.Get(freeSeats).Insert(car.Id, &car)
	p.workersMutex.Unlock()
	return true
}

func (p *Pooling) takeCarAvailable(people int) (car *Car, freeSeats int, ok bool) {
	for i := people; i <= MaxSeats; i++ {
		if car, ok := p.freeSeatsIdx.Get(i).TakeOne(); ok {
			return car, i, true
		}
	}
	return nil, 0, false
}
