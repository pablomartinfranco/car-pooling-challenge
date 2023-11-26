package domain

import (
	"errors"
	"fmt"
)

func (p *Pooling) JourneyTrigger(group *Group) error {

	if ok := p.IsGroupSizeValid(group); ok {
		return errors.New("invalid number of people")
	}

	if ok := p.IsGroupWaiting(group.Id); ok {
		return fmt.Errorf("group %d already waiting", group.Id)
	}

	if _, ok := p.IsGroupInJourney(group.Id); ok {
		return fmt.Errorf("group %d already in journey", group.Id)
	}

	p.waitingQueue.Enqueue(&group)

	p.journeyEvent.Signal()

	return nil
}
