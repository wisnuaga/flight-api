package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	testifymock "github.com/stretchr/testify/mock"
	"github.com/wisnuaga/flight-api/internal/delivery/http/handler"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/tests/mock"
)

func setupTestRouter(h *handler.FlightHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/api/v1/search", h.Search)
	return r
}

func TestFlightHandler_Search(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		usecaseMock := new(mock.MockFlightUsecase)
		h := handler.NewFlightHandler(&handler.FlightHandlerUsecases{
			FlightUsecase: usecaseMock,
		})
		r := setupTestRouter(h)

		usecaseMock.On("Search", testifymock.Anything, testifymock.AnythingOfType("*entity.SearchRequest")).Return(&entity.SearchResult{
			Flights: []*entity.Flight{
				{
					ID:           "F1",
					Provider:     "Garuda",
					FlightNumber: "GA123",
					Origin:       entity.Location{Airport: "CGK", Time: time.Now().UTC(), Timezone: time.UTC},
					Destination:  entity.Location{Airport: "DPS", Time: time.Now().Add(2 * time.Hour).UTC(), Timezone: time.UTC},
				},
			},
			Meta: &entity.SearchMeta{TotalFlights: 1},
		}, nil).Once()

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
		usecaseMock.AssertExpectations(t)
	})

	t.Run("invalid json", func(t *testing.T) {
		usecaseMock := new(mock.MockFlightUsecase)
		h := handler.NewFlightHandler(&handler.FlightHandlerUsecases{
			FlightUsecase: usecaseMock,
		})
		r := setupTestRouter(h)

		req, _ := http.NewRequest("POST", "/api/v1/search", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("invalid date bounds fallback", func(t *testing.T) {
		usecaseMock := new(mock.MockFlightUsecase)
		h := handler.NewFlightHandler(&handler.FlightHandlerUsecases{
			FlightUsecase: usecaseMock,
		})
		r := setupTestRouter(h)

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
		usecaseMock := new(mock.MockFlightUsecase)
		h := handler.NewFlightHandler(&handler.FlightHandlerUsecases{
			FlightUsecase: usecaseMock,
		})
		r := setupTestRouter(h)

		usecaseMock.On("Search", testifymock.Anything, testifymock.AnythingOfType("*entity.SearchRequest")).Return(nil, errors.New("upstream failed")).Once()

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
		usecaseMock.AssertExpectations(t)
	})
}
