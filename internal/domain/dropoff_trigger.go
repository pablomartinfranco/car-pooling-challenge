package domain

import "fmt"

func (p *Pooling) DropoffTrigger(id int) error {

	journey, ok := p.journeyIndex.Lookup(id)
	if !ok {
		return fmt.Errorf("group %d not in journey", id)
	}

	var hasGroupId = func(group *Group) bool {
		return group.Id == id
	}
	if ok := p.dropoffQueue.Any(hasGroupId); ok {
		return fmt.Errorf("group %d already in dropoffQueue", id)
	}

	p.dropoffQueue.Enqueue(&journey.Group)

	p.dropoffEvent.Signal()

	return nil
}
