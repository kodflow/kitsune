package service

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewExchange(t *testing.T) {
	testCases := []struct {
		name    string
		service string
		answer  bool
	}{
		{"ValidServiceWithAnswer", "exampleService", true},
		{"ValidServiceWithoutAnswer", "exampleService", false},
		{"EmptyServiceWithAnswer", "", true},
		{"EmptyServiceWithoutAnswer", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exchange := NewExchange(tc.service, tc.answer)

			if exchange == nil {
				t.Error("Expected non-nil exchange, got nil")
			}

			if exchange.Service != tc.service {
				t.Errorf("Expected service to be %s, got %s", tc.service, exchange.Service)
			}

			if exchange.Req.Endpoint != tc.service {
				t.Errorf("Expected Req.Endpoint to be %s, got %s", tc.service, exchange.Req.Endpoint)
			}

			if exchange.Answer != tc.answer {
				t.Errorf("Expected Answer to be %t, got %t", tc.answer, exchange.Answer)
			}

			if exchange.Req.Id == "" {
				t.Error("Expected non-empty Req.Id, got empty string")
			}

			_, err := uuid.Parse(exchange.Req.Id)
			if err != nil {
				t.Errorf("Expected Req.Id to be a valid UUID, got error: %v", err)
			}
		})
	}
}
