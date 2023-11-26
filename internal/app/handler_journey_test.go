package app

import (
	"bytes"
	"car-pooling-challenge/internal/domain"
	"car-pooling-challenge/pkg/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_journeyHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
		p *domain.Pooling
	}

	aConfig := &config.Config{
		JourneyWorkerPool: 1,
		DropoffWorkerPool: 1,
	}

	aPoolingService := domain.NewPooling(aConfig)

	aPoolingService.Run()

	someCars := []*domain.Car{
		{
			Id:    1,
			Seats: 4,
		},
		{
			Id:    2,
			Seats: 6,
		},
	}

	aPoolingService.CarsTrigger(someCars)

	time.Sleep(1 * time.Second)

	postDataOk := []byte(`{"id":1,"people":6}`)

	requestOk := httptest.NewRequest("POST", "/journey", bytes.NewBuffer(postDataOk))

	requestOk.Header.Set("Content-Type", "application/json")

	postDataBadRequest := []byte(`{"id":1,"people":7}`)

	requestBadRequest := httptest.NewRequest("POST", "/journey", bytes.NewBuffer(postDataBadRequest))

	requestBadRequest.Header.Set("Content-Type", "application/json")

	tests := []struct {
		name     string
		args     args
		expected int
	}{
		{
			name: "When register a group that is not waiting or in journey Then return 202 status code",
			args: args{
				w: httptest.NewRecorder(),
				r: requestOk,
				p: aPoolingService,
			},
			expected: http.StatusAccepted,
		},
		{
			name: "When register a group that is invalid Then return 400 status code",
			args: args{
				w: httptest.NewRecorder(),
				r: requestBadRequest,
				p: aPoolingService,
			},
			expected: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			journeyHandler(tt.args.w, tt.args.r, tt.args.p)
			response := tt.args.w.(*httptest.ResponseRecorder)
			if response.Code != tt.expected {
				t.Errorf("journeyHandler() = %v, want %v", response.Code, tt.expected)
			}
		})
	}
}
