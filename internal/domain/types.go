package domain

type Car struct {
	Id    int
	Seats int
}

type Group struct {
	Id     int
	People int
}

type Journey struct {
	Group *Group
	Car   *Car
}
