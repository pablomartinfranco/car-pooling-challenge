package domain

import "strconv"

func (p *Pooling) dropoffWorker(n int) {
	var worker = "Dropoff worker " + strconv.Itoa(n)
	p.logger.Printf("[%s] started", worker)
	defer p.logger.Printf("[%s] stopped", worker)

	for {
		p.DebugStates(worker)

		p.dropoffEvent.Wait()

		var groups = p.dropoffQueue.DequeueAll()

		for _, group := range groups {

			if ok := p.tryDropoff(group); !ok {

				p.logger.Printf("Error: Group %d dropoff fail", group.Id)

				p.handleDropoffFail(group)
			}
		}

		p.journeyEvent.Signal()

		if p.context.Err() != nil {
			return
		}
	}
}

func (p *Pooling) tryDropoff(group *Group) bool {
	p.workersMutex.Lock()
	journey, ok := p.journeyIndex.Remove(group.Id)
	if !ok {
		p.workersMutex.Unlock()
		return false
	}
	car, freeSeats, ok := p.takeCarFromIndex(journey.Car.Id)
	if !ok {
		p.workersMutex.Unlock()
		return false
	}
	freeSeats += group.People
	if freeSeats > MaxSeats {
		// FIX: rare race condition here, out of range index in freeSeatsIdx
		p.workersMutex.Unlock()
		return false
	}
	p.freeSeatsIdx.Get(freeSeats).Insert(journey.Car.Id, &car)
	p.workersMutex.Unlock()
	return true
}

func (p *Pooling) takeCarFromIndex(carId int) (car *Car, freeSeats int, ok bool) {
	var hasCarId = func(_ int, c *Car) bool {
		return c.Id == carId
	}
	for i := 0; i <= MaxSeats; i++ {
		if ok := p.freeSeatsIdx.Get(i).Any(hasCarId); ok {
			if car, ok := p.freeSeatsIdx.Get(i).Remove(carId); ok {
				return car, i, true
			}
			return nil, 0, false
		}
	}
	return nil, 0, false
}

func (p *Pooling) handleDropoffFail(group *Group) {
	p.dropoffMutex.Lock()
	if ok := p.dropoffRetry.HasKey(group.Id); !ok {
		p.dropoffQueue.Enqueue(&group)
		p.dropoffRetry.Insert(group.Id, &p.config.DropoffRetryLimit)
	} else if retry, ok := p.dropoffRetry.Lookup(group.Id); ok && retry > 0 {
		retry--
		p.dropoffQueue.Enqueue(&group)
		p.dropoffRetry.Update(group.Id, &retry)
	}
	p.dropoffMutex.Unlock()
}
