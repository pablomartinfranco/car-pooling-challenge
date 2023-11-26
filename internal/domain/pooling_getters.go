package domain

import (
	"log"
)

func (p *Pooling) GetLogger() *log.Logger {
	return p.logger
}

func (p *Pooling) CancelContext() {
	p.cancel()
}

func (p *Pooling) IsGroupSizeValid(g *Group) bool {
	return g.People < MinGroup || g.People > MaxGroup
}

func (p *Pooling) IsGroupWaiting(id int) bool {
	var hasGroupId = func(g *Group) bool {
		return g.Id == id
	}
	if ok := p.waitingQueue.Any(hasGroupId); ok {
		return true
	}
	return false
}

func (p *Pooling) IsGroupInJourney(id int) (*Journey, bool) {
	journey, ok := p.journeyIndex.Lookup(id)
	return journey, ok
}
