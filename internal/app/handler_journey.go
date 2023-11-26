package app

import (
	"car-pooling-challenge/internal/domain"
	"encoding/json"
	"log"
	"net/http"
)

func journeyHandler(
	w http.ResponseWriter, r *http.Request, p *domain.Pooling,
) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var group *domain.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := p.JourneyTrigger(group); err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
