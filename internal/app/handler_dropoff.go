package app

import (
	"car-pooling-challenge/internal/domain"
	"log"
	"net/http"
	"strconv"
)

func dropoffHandler(
	w http.ResponseWriter, r *http.Request, p *domain.Pooling,
) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	groupId, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := p.DropoffTrigger(groupId); err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
