# Flight Search & Aggregation System

A robust, high-performance flight search aggregation backend built in Go that supports both one-way and round-trip flight searches. This system concurrently searches multiple airline mock APIs (Garuda Indonesia, Lion Air, Batik Air, AirAsia), handles data inconsistencies, implements exponential backoff strategies, and performs intelligent flight matching with realistic layover constraints.

## Core Features

### Flight Search
- **One-Way & Round-Trip Searches**: Seamlessly support both search types with automatic detection
- **Concurrent Provider Aggregation**: Parallel fan-out fetching from all providers with strict 500ms timeouts
- **Robust Deduplication**: Remove codeshare duplicates using 5-minute departure/arrival bucketing
- **Exponential Backoff Retry**: Automatic retry logic with 3 attempts and expanding backoff (50ms, 100ms, 200ms)
- **Result Caching**: 5-minute in-memory caching per provider/route/date combination
- **Graceful Degradation**: Continue serving results even if some providers fail

### Round-Trip Intelligence
- **Automatic Itinerary Combination**: Intelligently combines outbound and return flights
- **Layover Validation**: Enforces 90-minute minimum and 24-hour maximum layovers
- **Combined Pricing**: Returns total price and duration for complete itineraries
- **Consistent Interface**: Single unified search method handles both one-way and round-trip

### Advanced Filtering & Sorting
- **Price Filtering**: Min/max price filters with decimal precision
- **Time Filtering**: Departure/arrival time ranges (UTC-based for consistency)
- **Stops Filtering**: Direct flights only or maximum stops allowed
- **Airline Filtering**: Filter by IATA code or provider name
- **Duration Filtering**: Maximum flight duration in minutes
- **Multiple Sort Options**: Price, duration, departure time, arrival time, or best-value ranking
- **Best-Value Ranking**: Normalized scoring combining price and duration

## Project Structure

```
flight-api/
├── cmd/
│   └── http/
│       └── main.go                 # Entry point, initializes providers and HTTP server
│
├── internal/
│   ├── command/
│   │   ├── flight_filter.go        # Filtering logic implementation
│   │   ├── flight_sorter.go        # Sorting logic implementation
│   │   ├── round_trip_combiner.go  # Round-trip itinerary combination
│   │   └── *_test.go               # Comprehensive unit tests
│   │
│   ├── config/
│   │   └── config.go               # Configuration management
│   │
│   ├── delivery/
│   │   └── http/
│   │       ├── handler/            # HTTP request handlers
│   │       ├── router/             # Route definitions and integration tests
│   │       └── dto/                # Request/response data structures
│   │
│   ├── domain/
│   │   └── entity/                 # Core business entities
│   │       ├── flight.go           # Flight entity and validation
│   │       ├── search.go           # Search request/result entities
│   │       ├── location.go         # Location with timezone
│   │       ├── layover.go          # Layover calculation
│   │       └── sort.go             # Sort parameters
│   │
│   ├── infra/
│   │   ├── cache/                  # In-memory caching implementation
│   │   └── provider/               # Airline provider integrations
│   │       ├── airasia/            # AirAsia client
│   │       ├── batikair/           # Batik Air client
│   │       ├── garuda/             # Garuda Indonesia client
│   │       ├── lionair/            # Lion Air client
│   │       └── registry.go         # Provider registry
│   │
│   ├── port/                       # Port/Interface definitions
│   │   ├── cache.go                # Cache interface
│   │   ├── command.go              # Command interfaces
│   │   ├── flight_provider.go      # Provider interface
│   │   └── usecase.go              # Usecase interface
│   │
│   ├── usecase/                    # Business logic orchestration
│   │   ├── flight_usecase.go       # Main search orchestration
│   │   └── flight_usecase_test.go  # Usecase tests
│   │
│   ├── util/                       # Utility functions
│   │   ├── flight.go               # Flight utilities
│   │   ├── price.go                # Price formatting
│   │   ├── time_parser.go          # Time parsing utilities
│   │   └── mock_loader.go          # Mock data loading
│   │
│   └── test_helper/                # Test utilities
│
├── tests/
│   ├── factory/                    # Mock API response JSON files
│   └── mock/                       # Mock implementations for testing
│
├── docs/
│   └── api.md                      # Detailed API documentation
│
├── Dockerfile                      # Container image definition
├── Makefile                        # Build and development tasks
├── go.mod                          # Go module dependencies
├── go.sum                          # Go module checksums
└── README.md                       # This file
```

## Architecture & Design Patterns

### Clean Architecture Principles
The codebase strictly follows Clean Architecture with clear separation of concerns:

1. **Domain Layer** (`internal/domain/entity/`): Core business rules and data structures
2. **Port Layer** (`internal/port/`): Interfaces defining contracts between layers
3. **Usecase Layer** (`internal/usecase/`): Orchestration of business logic
4. **Adapter Layer** (`internal/delivery/`, `internal/infra/`): External interface implementations
5. **Command Layer** (`internal/command/`): Specific business operations (filtering, sorting, combination)

### Key Design Decisions

#### Single Search Method for Both Search Types
Rather than maintaining separate `Search()` and `SearchRoundTrip()` methods, the system uses a single unified interface:
- The `ReturnDate` field in `SearchRequest` determines search type
- `nil` ReturnDate → one-way search
- Set ReturnDate → round-trip search
- **Benefit**: Simpler API, less code duplication, easier to maintain

#### Round-Trip Combination Logic
The `RoundTripCombiner` is implemented as a separate command:
- Takes all outbound and return flights as inputs
- Creates all valid combinations (O(n×m) combinations)
- Validates layover constraints (90 min - 24 hours)
- Returns only valid itineraries
- **Benefit**: Testable, reusable, follows command pattern like filters/sorters

#### Layover Constraints
Hard-coded realistic constraints:
- **Minimum 90 minutes**: Allows deplaning, immigration, baggage claim, re-boarding
- **Maximum 24 hours**: Prevents unrealistic multi-day layovers
- **Benefit**: Ensures user-friendly results without configuration complexity

#### In-Memory Caching per Provider
Caching key: `(provider, origin, destination, date, passengers)`
- 5-minute TTL per provider
- Isolated per provider to reduce coupling
- Cache hits bypass network requests entirely
- **Benefit**: Reduces redundant API calls, improves response times

#### Provider Isolation with Retry Logic
Each provider has independent retry strategy:
- 3 total attempts with exponential backoff: 50ms, 100ms, 200ms
- Failure in one provider doesn't affect others
- Results aggregate despite partial failures
- **Benefit**: Fault tolerance, graceful degradation

#### Deduplication Strategy
Uses 5-minute bucketing for departure and arrival times:
```
key = origin_airport + destination_airport + 
      departure_time.truncate(5m) + arrival_time.truncate(5m)
```
When duplicates found, keeps the cheaper flight
- **Benefit**: Catches codeshare flights without false positives, price-optimized

### Testing Strategy

#### Unit Tests
- **Command tests**: Filter and sort logic, round-trip combination
- **Entity tests**: Flight validation and normalization
- **Provider tests**: Individual provider response parsing
- **Usecase tests**: Aggregation, deduplication, error resilience

#### Integration Tests
- Full API endpoint testing with mock providers
- Response format validation
- Metadata accuracy verification

#### Test Coverage
```bash
go test ./internal/... -v        # Run all tests
go test ./internal/command -v    # Command layer tests
go test ./internal/usecase -v    # Usecase layer tests
go test ./internal/delivery/... -v  # HTTP handler tests
```

## Setup & Running the Service

### Prerequisites
- Go 1.19 or higher
- curl (for testing API)

### Installation & Setup

1. **Navigate to project directory**:
```bash
cd flight-api
```

2. **Install dependencies**:
```bash
go mod tidy
```

3. **Run the server**:
```bash
go run cmd/http/main.go
```

The server starts on `http://localhost:8080` by default.

### Testing

Run the complete test suite:
```bash
go test ./internal/... -v
```

Run specific test categories:
```bash
go test ./internal/command -v        # Filtering, sorting, round-trip tests
go test ./internal/usecase -v        # Business logic tests
go test ./internal/delivery/... -v   # Handler and integration tests
```

Build the application:
```bash
go build -o flight-api ./cmd/http
./flight-api
```

## API Usage Examples

### Health Check
```bash
curl http://localhost:8080/healthz
```

### One-Way Search
```bash
curl -X POST http://localhost:8080/flights/search \
  -H 'Content-Type: application/json' \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departure_date": "2025-12-15",
    "passengers": 1,
    "max_stops": 0,
    "sort_by": "price"
  }'
```

### Round-Trip Search
```bash
curl -X POST http://localhost:8080/flights/search \
  -H 'Content-Type: application/json' \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departure_date": "2025-12-15",
    "return_date": "2025-12-22",
    "passengers": 1,
    "sort_by": "price"
  }'
```

### With Filters
```bash
curl -X POST http://localhost:8080/flights/search \
  -H 'Content-Type: application/json' \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departure_date": "2025-12-15",
    "passengers": 1,
    "min_price": 500000,
    "max_price": 1500000,
    "airline_codes": ["GA", "JT"],
    "max_stops": 0,
    "sort_by": "best_value"
  }'
```

## Docker Support

Build and run with Docker:

```bash
docker build -t flight-api .
docker run -p 8080:8080 flight-api
```

## Documentation

For detailed API specification including request/response formats, field descriptions, and error handling:
- See [docs/api.md](docs/api.md) - Complete API documentation with examples

For technical implementation details:
- See [ROUND_TRIP_IMPLEMENTATION.md](ROUND_TRIP_IMPLEMENTATION.md) - Round-trip feature implementation
- See [DESIGN_CHOICES.md](DESIGN_CHOICES.md) - Detailed design decisions and architecture

## Performance Characteristics

- **Concurrent Requests**: Handles hundreds of concurrent flight searches
- **Search Latency**: ~300-500ms for fresh searches, <50ms for cached results
- **Memory Usage**: ~10MB base + ~5MB per 10,000 cached flight objects
- **Provider Timeout**: 500ms total per provider (with retries)
- **Deduplication**: O(N) single-pass algorithm
- **Caching**: 5-minute TTL, automatic eviction

## Response Examples

### One-Way Search Response
```json
{
  "search_criteria": {
    "origin": "CGK",
    "destination": "DPS",
    "departure_date": "2025-12-15",
    "passengers": 1,
    "cabin_class": "economy"
  },
  "metadata": {
    "total_results": 15,
    "providers_queried": 4,
    "providers_succeeded": 4,
    "providers_failed": 0,
    "search_time_ms": 445,
    "cache_hit": false
  },
  "flights": [
    {
      "id": "GA401_Garuda",
      "provider": "Garuda Indonesia",
      "airline": {
        "name": "Garuda Indonesia",
        "code": "GA"
      },
      "flight_number": "GA-401",
      "departure": {
        "airport": "CGK",
        "city": "Jakarta",
        "datetime": "2025-12-15T06:00:00+07:00",
        "timestamp": 1734224400
      },
      "arrival": {
        "airport": "DPS",
        "city": "Denpasar",
        "datetime": "2025-12-15T09:00:00+08:00",
        "timestamp": 1734227400
      },
      "duration": {
        "total_minutes": 120,
        "formatted": "2h"
      },
      "stops": 0,
      "price": {
        "amount": 750000,
        "currency": "IDR",
        "formatted": "IDR 750000"
      },
      "available_seats": 50,
      "cabin_class": "economy",
      "aircraft": "B737",
      "amenities": ["meals", "wifi"],
      "baggage": {
        "carry_on": "7kg",
        "checked": "20kg"
      }
    }
  ]
}
```

### Round-Trip Search Response
```json
{
  "search_criteria": {
    "origin": "CGK",
    "destination": "DPS",
    "departure_date": "2025-12-15",
    "return_date": "2025-12-22",
    "passengers": 1,
    "cabin_class": "economy"
  },
  "metadata": {
    "total_results": 12,
    "providers_queried": 4,
    "providers_succeeded": 4,
    "providers_failed": 0,
    "search_time_ms": 680,
    "cache_hit": false
  },
  "round_trip_itineraries": [
    {
      "outbound_flight": { /* flight details */ },
      "return_flight": { /* flight details */ },
      "total_price": {
        "amount": 1500000,
        "currency": "IDR",
        "formatted": "IDR 1500000"
      },
      "total_duration_minutes": 1686
    }
  ]
}
```

## Error Handling

The system handles various error scenarios gracefully:

- **Provider Failures**: Continues serving results from successful providers
- **Timeout**: Gracefully handles provider timeouts with exponential backoff
- **Invalid Input**: Returns 400 Bad Request with error message
- **No Results**: Returns 200 OK with empty results array
- **Internal Error**: Returns 500 Internal Server Error with error message

All errors are logged with trace IDs for debugging.

## Technologies Used

- **Language**: Go 1.19+
- **HTTP Framework**: Gin
- **JSON Parsing**: Standard library
- **Decimal Precision**: github.com/shopspring/decimal
- **Assertions (Testing)**: github.com/stretchr/testify
- **Logging**: log/slog (Go 1.21+ standard library)

## License

This project is part of the BookCabin technical interview assessment.
