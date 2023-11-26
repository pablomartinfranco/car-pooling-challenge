package app

import (
	"bytes"
	"car-pooling-challenge/internal/domain"
	"car-pooling-challenge/pkg/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func Test_locateHandler(t *testing.T) {
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

	requestOk := httptest.NewRequest(http.MethodPost, "/locate", bytes.NewBufferString(formDataOk))

	requestOk.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	form = url.Values{}

	form.Add("ID", "2")

	formDataNotFound := form.Encode()

	requestNotFound := httptest.NewRequest(http.MethodPost, "/locate", bytes.NewBufferString(formDataNotFound))

	requestNotFound.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	form = url.Values{}

	form.Add("ID", "//\\n3")

	formDataBadRequest := form.Encode()

	requestBadRequest := httptest.NewRequest(http.MethodPost, "/locate", bytes.NewBufferString(formDataBadRequest))

	requestBadRequest.Header.Set("Content-Type", "application/json")

	tests := []struct {
		name     string
		args     args
		expected int
	}{
		{
			name: "When locate group that is in journey Then return 200 status code",
			args: args{
				w: httptest.NewRecorder(),
				r: requestOk,
				p: aPoolingService,
			},
			expected: http.StatusOK,
		},
		{
			name: "When locate group not in journey and not waiting Then return 404 status code",
			args: args{
				w: httptest.NewRecorder(),
				r: requestNotFound,
				p: aPoolingService,
			},
			expected: http.StatusNotFound,
		},
		{
			name: "When locate group not in journey and not waiting Then return 404 status code",
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
			locateHandler(tt.args.w, tt.args.r, tt.args.p)
			response := tt.args.w.(*httptest.ResponseRecorder)
			if response.Code != tt.expected {
				t.Errorf("locateHandler() = %v, want %v", response.Code, tt.expected)
			}
		})
	}
}

func Test_locateHandlerReturnsJson(t *testing.T) {
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

	formData := form.Encode()

	payload, _ := json.Marshal(someCars[1])

	jsonCar := string(payload)

	request := httptest.NewRequest(http.MethodPost, "/locate", bytes.NewBufferString(formData))

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response := httptest.NewRecorder()

	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			name: "When locate group that is in journey Then get car from journey",
			args: args{
				w: response,
				r: request,
				p: aPoolingService,
			},
			expected: jsonCar,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locateHandler(tt.args.w, tt.args.r, tt.args.p)
			response := tt.args.w.(*httptest.ResponseRecorder)
			result := response.Body.String()
			if result != tt.expected {
				t.Errorf("locateHandler() = %v, want %v", result, tt.expected)
			}
		})
	}
}
