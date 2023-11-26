package domain

import "errors"

func (p *Pooling) CarsTrigger(cars []*Car) error {
	if err := validateCars(cars); err != nil {
		p.logger.Printf("Error: %v", err)
		return err
	}
	p.Reset()
	registerCars(cars, p)
	return nil
}

func registerCars(cars []*Car, p *Pooling) {
	for _, car := range cars {
		p.freeSeatsIdx.Get(car.Seats).Insert(car.Id, &car)
		p.logger.Printf("[CarsTrigger] Car %d registered with %d seats", car.Id, car.Seats)
	}
	p.logger.Printf("[CarsTrigger] Total of %d cars registered", len(cars))
}

func validateCars(cars []*Car) error {
	for _, car := range cars {
		if ok := validateCarSeats(car); !ok {
			return errors.New("invalid number of seats")
		}
	}
	return nil
}

func validateCarSeats(c *Car) bool {
	return c.Seats < MinSeats || c.Seats > MaxSeats
}
