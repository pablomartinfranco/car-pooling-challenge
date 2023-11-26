package app

import (
	"car-pooling-challenge/internal/domain"
	"car-pooling-challenge/pkg/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_statusHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
		p *domain.Pooling
	}

	var aConfig = &config.Config{}

	var aPoolingService = domain.NewPooling(aConfig)

	tests := []struct {
		name     string
		args     args
		expected int
	}{
		{
			name: "When method is not GET Then StatusMethodNotAllowed",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/status", nil),
				p: aPoolingService,
			},
			expected: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusHandler(tt.args.w, tt.args.r, tt.args.p)
			if tt.args.w.Code != tt.expected {
				t.Errorf("StatusCode = %v, want %v, %s", tt.args.w.Code, tt.expected, tt.name)
			}
		})
	}
}
