package app

import (
	"bytes"
	"car-pooling-challenge/internal/domain"
	"car-pooling-challenge/pkg/config"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func Test_dropoffHandler(t *testing.T) {
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

	aGroup := &domain.Group{
		Id:     1,
		People: 6,
	}

	aPoolingService.JourneyTrigger(aGroup)

	time.Sleep(1 * time.Second)

	form := url.Values{}

	form.Add("ID", "1")

	formDataOk := form.Encode()

	requestOk := httptest.NewRequest(http.MethodPost, "/dropoff", bytes.NewBufferString(formDataOk))

	requestOk.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	form = url.Values{}

	form.Add("ID", "2")

	formDataNotFound := form.Encode()

	requestNotFound := httptest.NewRequest(http.MethodPost, "/dropoff", bytes.NewBufferString(formDataNotFound))

	requestNotFound.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	form = url.Values{}

	form.Add("ID", "//\\n3")

	formDataBadRequest := form.Encode()

	requestBadRequest := httptest.NewRequest(http.MethodPost, "/dropoff", bytes.NewBufferString(formDataBadRequest))

	requestBadRequest.Header.Set("Content-Type", "application/json")

	tests := []struct {
		name     string
		args     args
		expected int
	}{
		{
			name: "When dropoff a group that is in journey Then return 200 status code",
			args: args{
				w: httptest.NewRecorder(),
				r: requestOk,
				p: aPoolingService,
			},
			expected: http.StatusOK,
		},
		{
			name: "When dropoff a group that is not in journey Then return 404 status code",
			args: args{
				w: httptest.NewRecorder(),
				r: requestNotFound,
				p: aPoolingService,
			},
			expected: http.StatusNotFound,
		},
		{
			name: "When dropoff a group that is invalid Then return 400 status code",
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
			dropoffHandler(tt.args.w, tt.args.r, tt.args.p)
			response := tt.args.w.(*httptest.ResponseRecorder)
			if response.Code != tt.expected {
				t.Errorf("dropoffHandler() = %v, want %v", response.Code, tt.expected)
			}
		})
	}
}
