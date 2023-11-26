package app

import (
	"car-pooling-challenge/internal/domain"
	"net/http"
)

func statusHandler(
	w http.ResponseWriter, r *http.Request, p *domain.Pooling,
) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	p.InspectStates("statusHandler")

	w.WriteHeader(http.StatusOK)
}
