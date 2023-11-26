package app

import (
	"car-pooling-challenge/internal/domain"
	"encoding/json"
	"log"
	"net/http"
)

func carsHandler(
	w http.ResponseWriter, r *http.Request, p *domain.Pooling,
) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var cars []*domain.Car
	if err := json.NewDecoder(r.Body).Decode(&cars); err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := p.CarsTrigger(cars); err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
