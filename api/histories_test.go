package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sfvn-hbduy/api"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockHistoriesService struct{}

type TestCase struct {
	Name        string
	QueryParams map[string]string
	Expected    int
}

func (m *mockHistoriesService) GetHistories(c context.Context, symbol string, days int, period string) (int, any) {
	return http.StatusOK, "mock response"
}

func TestGetHistories(t *testing.T) {
	r := gin.Default()
	historiesService := &mockHistoriesService{}
	api.APIHistoriesHandler(r, historiesService)

	testCases := []TestCase{
		{
			Name: "ValidInput",
			QueryParams: map[string]string{
				"start_date": "2024-02-23",
				"end_date":   "2024-02-24",
				"symbol":     "bitcoin",
			},
			Expected: http.StatusOK,
		},
		{
			Name: "MissingSymbol",
			QueryParams: map[string]string{
				"start_date": "2024-01-01",
				"end_date":   "2024-01-10",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "InvalidDates",
			QueryParams: map[string]string{
				"start_date": "2024-02-10",
				"end_date":   "2024-02-01",
				"symbol":     "bitcoin",
			},
			Expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/v1/get_histories", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			query := req.URL.Query()
			for key, value := range tc.QueryParams {
				query.Add(key, value)
			}
			req.URL.RawQuery = query.Encode()

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.Expected {
				t.Errorf("unexpected status code: got %d, want %d", w.Code, tc.Expected)
			}
		})
	}
}
