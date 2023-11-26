package app

import (
	"bytes"
	"car-pooling-challenge/internal/domain"
	"car-pooling-challenge/pkg/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_carsHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
		p *domain.Pooling
	}
	aConfig := &config.Config{}

	aPoolingService := domain.NewPooling(aConfig)

	goodCars := []byte(`[ { "id": 1, "seats": 4 }, { "id": 2, "seats": 6 } ]`)

	aGoodRequest := httptest.NewRequest(http.MethodPut, "/cars", bytes.NewBuffer(goodCars))
	aGoodRequest.Header.Set("Content-Type", "application/json")

	badCarsMin := []byte(`[ { "id": 1, "seats": 3 }, { "id": 2, "seats": 6 } ]`)

	aBadRequestMin := httptest.NewRequest(http.MethodPut, "/cars", bytes.NewBuffer(badCarsMin))
	aBadRequestMin.Header.Set("Content-Type", "application/json")

	badCarsMax := []byte(`[ { "id": 1, "seats": 1 }, { "id": 2, "seats": 7 } ]`)

	aBadRequestMax := httptest.NewRequest(http.MethodPut, "/cars", bytes.NewBuffer(badCarsMax))
	aBadRequestMax.Header.Set("Content-Type", "application/json")

	aBadRequestMethod := httptest.NewRequest(http.MethodPost, "/cars", bytes.NewBuffer(goodCars))
	aBadRequestMethod.Header.Set("Content-Type", "application/json")

	aBadRequestContent := httptest.NewRequest(http.MethodPut, "/cars", bytes.NewBuffer(goodCars))
	aBadRequestContent.Header.Set("Content-Type", "application/xml")

	tests := []struct {
		name     string
		args     args
		expected int
	}{
		{
			name: "When cars from 4 to 6 seats Then StatusOk",
			args: args{
				w: httptest.NewRecorder(),
				r: aGoodRequest,
				p: aPoolingService,
			},
			expected: http.StatusOK,
		},
		{
			name: "When cars below 1 seats Then StatusBadRequest",
			args: args{
				w: httptest.NewRecorder(),
				r: aBadRequestMin,
				p: aPoolingService,
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "When cars above 6 seats Then StatusBadRequest",
			args: args{
				w: httptest.NewRecorder(),
				r: aBadRequestMax,
				p: aPoolingService,
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "When method is not PUT Then StatusMethodNotAllowed",
			args: args{
				w: httptest.NewRecorder(),
				r: aBadRequestMethod,
				p: aPoolingService,
			},
			expected: http.StatusMethodNotAllowed,
		},
		{
			name: "When ContentType is not applicationJson Then StatusBadRequest",
			args: args{
				w: httptest.NewRecorder(),
				r: aBadRequestContent,
				p: aPoolingService,
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			carsHandler(tt.args.w, tt.args.r, tt.args.p)
			if tt.args.w.Code != tt.expected {
				t.Errorf("StatusCode = %v, want %v, %s", tt.args.w.Code, tt.expected, tt.name)
			}
		})
	}
}
