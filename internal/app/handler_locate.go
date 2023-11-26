package app

import (
	"car-pooling-challenge/internal/domain"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func locateHandler(
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

	// Acceptance tests won't PASS if validating Accept header

	// var accept = r.Header.Get("Accept")
	// if accept == "" || !strings.Contains(accept, "application/json") {
	// 	w.WriteHeader(http.StatusNotAcceptable)
	// 	return
	// }

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

	if ok := p.IsGroupWaiting(groupId); ok {
		err := fmt.Errorf("group %d still waiting", groupId)
		p.GetLogger().Printf("Error: %v", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if j, ok := p.IsGroupInJourney(groupId); ok {
		json, err := json.Marshal(j.Car)
		if err != nil {
			log.Printf("Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(json)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
