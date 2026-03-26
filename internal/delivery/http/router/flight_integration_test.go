package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wisnuaga/flight-api/internal/config"
	"github.com/wisnuaga/flight-api/internal/delivery/http/router"
)

// testDataDir resolves the absolute path to tests/factory/ from any working
// directory the test runner may choose (e.g. the package dir or the repo root).
func testDataDir() string {
	_, file, _, _ := runtime.Caller(0)
	// this file lives at internal/delivery/http/router; walk up 5 levels to repo root
	root := filepath.Join(filepath.Dir(file), "../../../../..")
	return filepath.Join(root, "tests/factory")
}

func newIntegrationConfig() *config.Config {
	dir := testDataDir()
	return &config.Config{
		Providers:      []string{"garuda", "lionair", "batikair", "airasia"},
		GarudaConfig:   config.ProviderConfig{MockPath: filepath.Join(dir, "garuda_search_response.json")},
		LionAirConfig:  config.ProviderConfig{MockPath: filepath.Join(dir, "lion_air_search_response.json")},
		BatikAirConfig: config.ProviderConfig{MockPath: filepath.Join(dir, "batik_air_search_response.json")},
		AirAsiaConfig:  config.ProviderConfig{MockPath: filepath.Join(dir, "airasia_search_response.json")},
	}
}

func doSearch(t *testing.T, body map[string]interface{}) *httptest.ResponseRecorder {
	t.Helper()
	eng := router.Setup(newIntegrationConfig())

	raw, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/flights/search", bytes.NewBuffer(raw))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	return w
}

// baseSearchBody is a minimal valid request covering the CGK→DPS route present
// in all four mock files.
func baseSearchBody() map[string]interface{} {
	return map[string]interface{}{
		"origin":         "CGK",
		"destination":    "DPS",
		"departure_date": "2025-12-15",
		"passengers":     1,
	}
}

func TestFlightSearch_Integration_StatusOK(t *testing.T) {
	w := doSearch(t, baseSearchBody())

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	// Top-level structure
	assert.Contains(t, resp, "flights", "response should contain flights key")
	assert.Contains(t, resp, "metadata", "response should contain metadata key")
	assert.Contains(t, resp, "search_criteria", "response should contain search_criteria key")
}

func TestFlightSearch_Integration_MetadataCorrect(t *testing.T) {
	w := doSearch(t, baseSearchBody())
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	meta, ok := resp["metadata"].(map[string]interface{})
	require.True(t, ok, "metadata should be an object")

	// All 4 providers are configured
	assert.Equal(t, float64(4), meta["providers_queried"], "providers_queried should be 4")
	assert.Greater(t, meta["providers_succeeded"], float64(0), "at least one provider should succeed")
}

func TestFlightSearch_Integration_FlightStructure(t *testing.T) {
	w := doSearch(t, baseSearchBody())
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	flights, ok := resp["flights"].([]interface{})
	require.True(t, ok, "flights should be an array")
	require.NotEmpty(t, flights, "at least one flight should be returned")

	// Check the structure of the first flight
	f, ok := flights[0].(map[string]interface{})
	require.True(t, ok)

	requiredFields := []string{"id", "provider", "airline", "flight_number", "departure", "arrival", "duration", "stops", "price", "available_seats", "cabin_class"}
	for _, field := range requiredFields {
		assert.Contains(t, f, field, "flight should have field: %s", field)
	}

	dep, ok := f["departure"].(map[string]interface{})
	require.True(t, ok, "departure should be an object")
	assert.Contains(t, dep, "airport")
	assert.Contains(t, dep, "city")
	assert.Contains(t, dep, "datetime")  // RFC3339 with TZ offset
	assert.Contains(t, dep, "timestamp") // Unix epoch seconds

	price, ok := f["price"].(map[string]interface{})
	require.True(t, ok, "price should be an object")
	assert.Contains(t, price, "amount")
	assert.Contains(t, price, "currency")
	assert.Contains(t, price, "formatted")
}

func TestFlightSearch_Integration_FilterByMaxStops(t *testing.T) {
	body := baseSearchBody()
	maxStops := 0
	body["max_stops"] = maxStops

	w := doSearch(t, body)
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	flights, ok := resp["flights"].([]interface{})
	require.True(t, ok)

	for _, fl := range flights {
		f := fl.(map[string]interface{})
		assert.Equal(t, float64(0), f["stops"], "all returned flights should be non-stop")
	}
}

func TestFlightSearch_Integration_FilterByPrice(t *testing.T) {
	body := baseSearchBody()
	body["max_price"] = 800000 // IDR

	w := doSearch(t, body)
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	flights, ok := resp["flights"].([]interface{})
	require.True(t, ok)

	for _, fl := range flights {
		f := fl.(map[string]interface{})
		price := f["price"].(map[string]interface{})
		amount := price["amount"].(float64)
		assert.LessOrEqual(t, amount, float64(800000), "price should be <= max_price filter")
	}
}

func TestFlightSearch_Integration_FilterByAirlineIATA(t *testing.T) {
	body := baseSearchBody()
	body["airline_codes"] = []string{"GA"} // Garuda IATA code

	w := doSearch(t, body)
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	flights, ok := resp["flights"].([]interface{})
	require.True(t, ok)
	require.NotEmpty(t, flights, "Garuda flights should be returned when filtering by IATA code GA")

	for _, fl := range flights {
		f := fl.(map[string]interface{})
		airline := f["airline"].(map[string]interface{})
		assert.Equal(t, "GA", airline["code"], "all returned flights should be from Garuda (code GA)")
	}
}

func TestFlightSearch_Integration_SortByPriceAsc(t *testing.T) {
	body := baseSearchBody()
	body["sort_by"] = "price"
	body["sort_order"] = "asc"

	w := doSearch(t, body)
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	flights, ok := resp["flights"].([]interface{})
	require.True(t, ok)
	if len(flights) < 2 {
		t.Skip("not enough flights to verify sort order")
	}

	for i := 1; i < len(flights); i++ {
		prev := flights[i-1].(map[string]interface{})["price"].(map[string]interface{})["amount"].(float64)
		curr := flights[i].(map[string]interface{})["price"].(map[string]interface{})["amount"].(float64)
		assert.LessOrEqual(t, prev, curr, "flights should be sorted by price ascending at index %d", i)
	}
}

func TestFlightSearch_Integration_SortByPriceDesc(t *testing.T) {
	body := baseSearchBody()
	body["sort_by"] = "price"
	body["sort_order"] = "desc"

	w := doSearch(t, body)
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	flights, ok := resp["flights"].([]interface{})
	require.True(t, ok)
	if len(flights) < 2 {
		t.Skip("not enough flights to verify sort order")
	}

	for i := 1; i < len(flights); i++ {
		prev := flights[i-1].(map[string]interface{})["price"].(map[string]interface{})["amount"].(float64)
		curr := flights[i].(map[string]interface{})["price"].(map[string]interface{})["amount"].(float64)
		assert.GreaterOrEqual(t, prev, curr, "flights should be sorted by price descending at index %d", i)
	}
}

func TestFlightSearch_Integration_SortByDuration(t *testing.T) {
	body := baseSearchBody()
	body["sort_by"] = "duration"
	body["sort_order"] = "asc"

	w := doSearch(t, body)
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	flights, ok := resp["flights"].([]interface{})
	require.True(t, ok)
	if len(flights) < 2 {
		t.Skip("not enough flights to verify sort order")
	}

	for i := 1; i < len(flights); i++ {
		prev := flights[i-1].(map[string]interface{})["duration"].(map[string]interface{})["total_minutes"].(float64)
		curr := flights[i].(map[string]interface{})["duration"].(map[string]interface{})["total_minutes"].(float64)
		assert.LessOrEqual(t, prev, curr, "flights should be sorted by duration ascending at index %d", i)
	}
}

func TestFlightSearch_Integration_SortByBestValue(t *testing.T) {
	body := baseSearchBody()
	body["sort_by"] = "best_value"
	body["sort_order"] = "asc"

	w := doSearch(t, body)
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	flights, ok := resp["flights"].([]interface{})
	require.True(t, ok)
	assert.NotEmpty(t, flights, "best_value sort should return results")
}

func TestFlightSearch_Integration_InvalidDate(t *testing.T) {
	body := baseSearchBody()
	body["departure_date"] = "not-a-date"

	w := doSearch(t, body)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFlightSearch_Integration_InvalidJSON(t *testing.T) {
	eng := router.Setup(newIntegrationConfig())

	req, _ := http.NewRequest(http.MethodPost, "/flights/search", bytes.NewBufferString("{bad json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	eng.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
