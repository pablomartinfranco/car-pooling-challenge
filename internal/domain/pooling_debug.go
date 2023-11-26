package domain

func (p *Pooling) DebugStates(caller string) {
	p.logger.Printf("[%s] dropoff groups %d", caller, p.dropoffQueue.Size())
	p.logger.Printf("[%s] waiting groups %d", caller, p.waitingQueue.Size())
	p.logger.Printf("[%s] journey groups %d", caller, p.journeyIndex.Size())
}

func (p *Pooling) InspectStates(caller string) {
	p.inspectFreeSeatsIdx(caller)
	p.inspectWaitingQueue(caller)
	p.inspectJourneyIndex(caller)
	p.inspectDropoffRetry(caller)
	p.inspectDropoffQueue(caller)
}

func (p *Pooling) inspectDropoffQueue(caller string) {
	p.dropoffQueue.ForEach(
		func(group *Group) {
			p.logger.Printf("[%s] dropoffQueue: Group %d, People %d", caller, group.Id, group.People)
		},
	)
}

func (p *Pooling) inspectWaitingQueue(caller string) {
	p.waitingQueue.ForEach(
		func(group *Group) {
			p.logger.Printf("[%s] waitingQueue: Group %d, People %d", caller, group.Id, group.People)
		},
	)
}

func (p *Pooling) inspectFreeSeatsIdx(caller string) {
	for i := 0; i <= MaxSeats; i++ {
		p.freeSeatsIdx.Get(i).ForEach(
			func(id int, car *Car) {
				p.logger.Printf("[%s] freeSeatsIdx[%d]: Car %d, Seats %d", caller, i, id, car.Seats)
			},
		)
	}
}

func (p *Pooling) inspectJourneyIndex(caller string) {
	p.journeyIndex.ForEach(
		func(id int, journey *Journey) {
			p.logger.Printf("[%s] journeyIndex: Group %d, Car %d", caller, id, journey.Car.Id)
		},
	)
}

func (p *Pooling) inspectDropoffRetry(caller string) {
	p.dropoffRetry.ForEach(
		func(id int, retry int) {
			p.logger.Printf("[%s] dropoffRetry: Group %d, Retry %d", caller, id, retry)
		},
	)
}
