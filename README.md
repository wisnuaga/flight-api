# Flight Search & Aggregation System

A robust, high-performance flight search aggregation backend built in Go. This system concurrently searches multiple airline mock APIs (Garuda, Lion Air, Batik Air, AirAsia), handles data inconsistencies, implements exponential backoff strategies, and performs in-line isolated memory filtering.

## Overview Features
- **Concurrent Aggregation**: Parallel fan-out fetching using strict 500ms context timeouts mapped across providers safely.
- **Robust Deduplication**: Removes duplicate code-shares natively using an innovative, absolute UNIX combination mapping.
- **Exponential Backoff**: The system will explicitly retry flaky network streams automatically 3 times independently with expanding timeout cushions (useful for simulating external networks).
- **Graceful Parsing Fallbacks**: Built to read standard ISO 8601 definitions but capable of shifting and dropping broken API arrays reliably without taking the engine offline.
- **Best Value Ranking Strategy**: Evaluates mathematical value heuristics parsing the cost of price vs the time lost scaling (Configurable dynamically).

## Architecture Approach
Explicitly built utilizing strict Clean Architecture principles to keep components thoroughly abstracted without leakage:
- **Domain Layer**: Core structs defining `Flight` and validation pipelines guarding business structures.
- **Provider Layer (Repository)**: Segregated HTTP clients imitating unique physical integration responses mimicking JSON payload anomalies natively.
- **Usecase Layer**: In-place filtering and map-based aggregation slicing memory limits aggressively to serve hundreds of concurrent clients safely.
- **Delivery Handler**: Enforces strict payload bindings generating localized `IDR Formatted` currency payloads easily parsed by a UI array.

## Setup & Running the Service

1. Pull the dependencies locally:
```bash
go mod tidy
```

2. Run the main listener component:
```bash
go run cmd/main.go
```

3. Initiate a sample UI Search string to trigger aggregation mapping:
```bash
curl -X POST http://localhost:8080/api/v1/search \
-H "Content-Type: application/json" \
-d '{
  "origin": "CGK",
  "destination": "DPS",
  "departureDate": "2025-12-15",
  "passengers": 1,
  "cabinClass": "economy"
}'
```

## Evaluating Code Test Coverages
Check edge case behavior against filtering logic and provider fault assumptions immediately natively via standard test sets:
```bash
go test ./internal/... -v
```
