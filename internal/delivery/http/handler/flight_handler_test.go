package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wisnuaga/flight-api/internal/delivery/http/handler"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
)

type mockFlightUsecase struct {
	mockSearch func(ctx context.Context, req *entity.SearchRequest) (*entity.SearchResult, error)
}

func (m *mockFlightUsecase) Search(ctx context.Context, req *entity.SearchRequest) (*entity.SearchResult, error) {
	if m.mockSearch != nil {
		return m.mockSearch(ctx, req)
	}
	return nil, nil
}

func setupTestRouter(h *handler.FlightHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/api/v1/search", h.Search)
	return r
}

func TestFlightHandler_Search(t *testing.T) {
	usecaseMock := &mockFlightUsecase{}
	h := handler.NewFlightHandler(&handler.FlightHandlerUsecases{
		FlightUsecase: usecaseMock,
	})

	r := setupTestRouter(h)

	t.Run("success", func(t *testing.T) {
		usecaseMock.mockSearch = func(ctx context.Context, req *entity.SearchRequest) (*entity.SearchResult, error) {
			return &entity.SearchResult{
				Flights: []*entity.Flight{
					{
						ID:            "F1",
						Provider:      "Garuda",
						FlightNumber:  "GA123",
						DepartureTime: time.Now(),
						ArrivalTime:   time.Now().Add(2 * time.Hour),
					},
				},
				Meta: &entity.SearchMeta{TotalFlights: 1},
			}, nil
		}

		body := map[string]interface{}{
			"origin":         "CGK",
			"destination":    "DPS",
			"departure_date": "2025-12-15",
			"passengers":     1,
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/api/v1/search", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/search", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("invalid date bounds fallback", func(t *testing.T) {
		body := map[string]interface{}{
			"origin":         "CGK",
			"destination":    "DPS",
			"departure_date": "invalid-date",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/search", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("usecase error", func(t *testing.T) {
		usecaseMock.mockSearch = func(ctx context.Context, req *entity.SearchRequest) (*entity.SearchResult, error) {
			return nil, errors.New("upstream failed")
		}

		body := map[string]interface{}{
			"origin":         "CGK",
			"destination":    "DPS",
			"departure_date": "2025-12-15",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/search", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}
