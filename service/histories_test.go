package service

import (
	"context"
	"errors"
	"testing"
	"time"
)

type TestCase struct {
	Name          string
	Symbol        string
	Days          int
	Period        string
	ExpectedError error
}

func TestGetHistories(t *testing.T) {
	svc := NewHistories("https://api.coingecko.com/api/v3", "CG-jztzX6oPA8z2ekw5Do7wZwXW")
	testCases := []TestCase{
		{
			Name:          "ValidInput",
			Symbol:        "bitcoin",
			Days:          1,
			Period:        "daily",
			ExpectedError: nil,
		},
		{
			Name:          "NegativeDays",
			Symbol:        "eth",
			Days:          -1,
			Period:        "daily",
			ExpectedError: errors.New("Out of range"),
		},
		{
			Name:          "NegativeDays",
			Symbol:        "",
			Days:          1,
			Period:        "daily",
			ExpectedError: errors.New("symbol is required"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()
			_, err := svc.GetHistories(ctx, tc.Symbol, tc.Days, tc.Period)
			if err != tc.ExpectedError {
				t.Errorf("unexpected error, got: %v, want: %v", err, tc.ExpectedError)
			}
		})
	}
}
