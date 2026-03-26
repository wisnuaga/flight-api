# Flight Search API Documentation

## Overview

The Flight API aggregates real-time flight data from multiple airline providers (Garuda, Lion Air, Batik Air, AirAsia) into a single, normalized search endpoint. It handles provider failures gracefully, deduplicates codeshare flights, caches results per provider, and applies configurable filtering and sorting.

---

## Endpoint

### `POST /flights/search`

Search for available flights across all configured providers. Supports both one-way and round-trip searches.

---

## Request

**Content-Type:** `application/json`

### Body Fields

| Field | Type | Required | Description |
|---|---|---|---|
| `origin` | string | Yes | Departure airport IATA code (e.g. `"CGK"`) |
| `destination` | string | Yes | Arrival airport IATA code (e.g. `"DPS"`) |
| `departure_date` | string | Yes | Departure date in `YYYY-MM-DD` format |
| `return_date` | string | No | Return date in `YYYY-MM-DD` format. Omit for one-way search, include for round-trip search |
| `passengers` | integer | No | Number of passengers (default: 0 = any) |
| `cabin_class` | string | No | Filter by cabin class (e.g. `"economy"`, `"business"`) |
| `min_price` | number | No | Minimum price filter (currency matches provider, typically IDR) |
| `max_price` | number | No | Maximum price filter |
| `max_stops` | integer | No | Maximum number of stops (0 = direct only) |
| `departure_start` | string | No | Earliest departure time (RFC3339 or `YYYY-MM-DD HH:MM`) |
| `departure_end` | string | No | Latest departure time (RFC3339 or `YYYY-MM-DD HH:MM`) |
| `arrival_start` | string | No | Earliest arrival time |
| `arrival_end` | string | No | Latest arrival time |
| `airline_codes` | array of strings | No | Filter by IATA code (e.g. `"GA"`) **or** provider name (e.g. `"Garuda Indonesia"`) |
| `max_duration` | integer | No | Maximum flight duration in **minutes** |
| `sort_by` | string | No | Sort field: `price` \| `duration` \| `departure_time` \| `arrival_time` \| `best_value` (default: `price`) |
| `sort_order` | string | No | Sort direction: `asc` \| `desc` (default: `asc`) |

> **Note:** All time filters are compared against UTC-normalised departure/arrival times. Provide times with a timezone offset (e.g. `2025-12-15T06:00:00+07:00`) for accurate filtering, or bare UTC timestamps.
>
> **Round-Trip:** Include `return_date` to search for round-trip itineraries. The API will combine outbound and return flights, enforcing a minimum 90-minute layover and maximum 24-hour layover between flights.

---

## Response

**Content-Type:** `application/json`

For one-way searches, the response contains:

```json
{
  "search_criteria": { ... },
  "metadata": { ... },
  "flights": [ ... ]
}
```

For round-trip searches, the response contains:

```json
{
  "search_criteria": { ... },
  "metadata": { ... },
  "round_trip_itineraries": [ ... ]
}
```

### `search_criteria`

| Field | Type | Description |
|---|---|---|
| `origin` | string | Origin airport code from request |
| `destination` | string | Destination airport code from request |
| `departure_date` | string | Departure date from request |
| `return_date` | string | Return date from request (round-trip only) |
| `passengers` | integer | Passenger count from request |
| `cabin_class` | string | Cabin class from request |

### `metadata`

| Field | Type | Description |
|---|---|---|
| `total_results` | integer | Number of flights returned (after dedup + filters) or itineraries (round-trip) |
| `providers_queried` | integer | Total number of configured providers |
| `providers_succeeded` | integer | Number of providers that returned data |
| `providers_failed` | integer | Number of providers that failed or timed out |
| `search_time_ms` | integer | Total end-to-end search time in milliseconds |
| `cache_hit` | boolean | `true` if at least one provider result was served from cache |

### `flights[]` (One-Way Search)

Each element in the `flights` array:

| Field | Type | Description |
|---|---|---|
| `id` | string | Unique flight identifier (`<flight_number>_<provider>`) |
| `provider` | string | Airline display name (e.g. `"Garuda Indonesia"`) |
| `airline.name` | string | Airline display name |
| `airline.code` | string | IATA airline code (e.g. `"GA"`) |
| `flight_number` | string | Flight number (e.g. `"GA-404"`) |
| `departure.airport` | string | Departure airport IATA code |
| `departure.city` | string | Departure city name |
| `departure.datetime` | string | Local departure time in RFC3339 format (with original TZ offset) |
| `departure.timestamp` | integer | Departure time as Unix epoch seconds (UTC) |
| `arrival.airport` | string | Arrival airport IATA code |
| `arrival.city` | string | Arrival city name |
| `arrival.datetime` | string | Local arrival time in RFC3339 format (with original TZ offset) |
| `arrival.timestamp` | integer | Arrival time as Unix epoch seconds (UTC) |
| `duration.total_minutes` | integer | Total flight duration in minutes |
| `duration.formatted` | string | Human-readable duration (e.g. `"2h 30m"`) |
| `stops` | integer | Number of stops (0 = direct) |
| `price.amount` | integer | Price as an integer in the currency unit |
| `price.currency` | string | Currency code (e.g. `"IDR"`) |
| `price.formatted` | string | Human-readable price (e.g. `"IDR 750000"`) |
| `available_seats` | integer | Number of seats available |
| `cabin_class` | string | Cabin class (e.g. `"economy"`, `"business"`) |
| `aircraft` | string \| null | Aircraft type if provided by the source |
| `amenities` | array | List of amenities (may be empty) |
| `baggage.carry_on` | string | Carry-on baggage allowance |
| `baggage.checked` | string | Checked baggage allowance |

### `round_trip_itineraries[]` (Round-Trip Search)

Each element in the `round_trip_itineraries` array represents a complete round-trip combination:

| Field | Type | Description |
|---|---|---|
| `outbound_flight` | object | Outbound flight details (see `flights[]` structure above) |
| `return_flight` | object | Return flight details (see `flights[]` structure above) |
| `total_price.amount` | integer | Sum of outbound and return flight prices |
| `total_price.currency` | string | Currency code (e.g. `"IDR"`) |
| `total_price.formatted` | string | Human-readable total price (e.g. `"IDR 1500000"`) |
| `total_duration_minutes` | integer | Sum of both flight durations plus layover time |
| `layover_minutes` | integer | Layover duration between return flight departure and outbound arrival |

---

## Error Responses

| HTTP Status | Cause |
|---|---|
| `400 Bad Request` | Invalid JSON body or unparseable `departure_date` or `return_date` |
| `500 Internal Server Error` | All providers failed and usecase returned an error |

### Error Body

```json
{
  "error": "<description>"
}
```

---

## Behaviour Details

### Provider Aggregation
Flights are fetched from all configured providers **concurrently**. Each provider gets up to **500 ms** total (shared across 3 retry attempts with exponential backoff: 50 ms, 100 ms, 200 ms). If a provider fails all retries, its error is counted in `providers_failed` and results continue without it.

### Deduplication
Codeshare flights are deduplicated using a key of `origin_airport + destination_airport + departure_5min_bucket + arrival_5min_bucket`. When duplicates exist, the **cheaper flight** is kept.

### Caching
Results per `(provider, origin, destination, date, passengers)` tuple are cached in-memory for **5 minutes**. Cached results bypass the provider HTTP call entirely. `metadata.cache_hit = true` when at least one provider was served from cache.

### Round-Trip Combination
For round-trip searches, outbound and return flights are combined with the following rules:
- Minimum layover: 90 minutes (1.5 hours)
- Maximum layover: 24 hours
- Total price: Sum of outbound and return flight prices
- Results are ranked by lowest total price by default

### Timezone Handling
All internal comparisons (filters, sort) use **UTC**. Departure/arrival timestamps in the response are returned in the provider's **original timezone** (via RFC3339 `datetime`) alongside a UTC Unix `timestamp`. This means a `06:00 WIB` departure displays as `2025-12-15T06:00:00+07:00`.

### Best-Value Ranking
When `sort_by=best_value`, each flight receives a score `= (normalised_price × price_weight) + (normalised_duration × duration_weight)`. Scores are normalised across the result set to `[0, 1]`. Flights with the **lowest score** rank highest (most value). Default weights are both `1.0`.

---

## Examples

### One-Way Search

#### Request

```bash
curl -s -X POST http://localhost:8080/flights/search \
  -H 'Content-Type: application/json' \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departure_date": "2025-12-15",
    "passengers": 1,
    "max_stops": 0,
    "max_price": 1000000,
    "sort_by": "price",
    "sort_order": "asc"
  }'
```

#### Response (truncated)

```json
{
  "search_criteria": {
    "origin": "CGK",
    "destination": "DPS",
    "departure_date": "2025-12-15",
    "passengers": 1,
    "cabin_class": ""
  },
  "metadata": {
    "total_results": 5,
    "providers_queried": 4,
    "providers_succeeded": 4,
    "providers_failed": 0,
    "search_time_ms": 312,
    "cache_hit": false
  },
  "flights": [
    {
      "id": "QZ520_AirAsia",
      "provider": "AirAsia",
      "airline": { "name": "AirAsia", "code": "QZ" },
      "flight_number": "QZ520",
      "departure": {
        "airport": "CGK",
        "city": "Jakarta",
        "datetime": "2025-12-15T04:45:00+07:00",
        "timestamp": 1734216300
      },
      "arrival": {
        "airport": "DPS",
        "city": "Denpasar",
        "datetime": "2025-12-15T08:25:00+08:00",
        "timestamp": 1734222300
      },
      "duration": { "total_minutes": 100, "formatted": "1h 40m" },
      "stops": 0,
      "price": { "amount": 650000, "currency": "IDR", "formatted": "IDR 650000" },
      "available_seats": 67,
      "cabin_class": "economy",
      "aircraft": null,
      "amenities": [],
      "baggage": { "carry_on": "", "checked": "" }
    }
  ]
}
```

### Round-Trip Search

#### Request

```bash
curl -s -X POST http://localhost:8080/flights/search \
  -H 'Content-Type: application/json' \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departure_date": "2025-12-15",
    "return_date": "2025-12-22",
    "passengers": 1,
    "sort_by": "price",
    "sort_order": "asc"
  }'
```

#### Response (truncated)

```json
{
  "search_criteria": {
    "origin": "CGK",
    "destination": "DPS",
    "departure_date": "2025-12-15",
    "return_date": "2025-12-22",
    "passengers": 1,
    "cabin_class": ""
  },
  "metadata": {
    "total_results": 12,
    "providers_queried": 4,
    "providers_succeeded": 4,
    "providers_failed": 0,
    "search_time_ms": 425,
    "cache_hit": false
  },
  "round_trip_itineraries": [
    {
      "outbound_flight": {
        "id": "GA401_Garuda",
        "provider": "Garuda",
        "airline": { "name": "Garuda Indonesia", "code": "GA" },
        "flight_number": "GA401",
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
        "duration": { "total_minutes": 120, "formatted": "2h" },
        "stops": 0,
        "price": { "amount": 750000, "currency": "IDR", "formatted": "IDR 750000" },
        "available_seats": 50,
        "cabin_class": "economy",
        "aircraft": "B737",
        "amenities": ["meals", "wifi"],
        "baggage": { "carry_on": "7kg", "checked": "20kg" }
      },
      "return_flight": {
        "id": "GA402_Garuda",
        "provider": "Garuda",
        "airline": { "name": "Garuda Indonesia", "code": "GA" },
        "flight_number": "GA402",
        "departure": {
          "airport": "DPS",
          "city": "Denpasar",
          "datetime": "2025-12-22T14:00:00+08:00",
          "timestamp": 1735030800
        },
        "arrival": {
          "airport": "CGK",
          "city": "Jakarta",
          "datetime": "2025-12-22T17:00:00+07:00",
          "timestamp": 1735035600
        },
        "duration": { "total_minutes": 120, "formatted": "2h" },
        "stops": 0,
        "price": { "amount": 750000, "currency": "IDR", "formatted": "IDR 750000" },
        "available_seats": 50,
        "cabin_class": "economy",
        "aircraft": "B737",
        "amenities": ["meals", "wifi"],
        "baggage": { "carry_on": "7kg", "checked": "20kg" }
      },
      "total_price": { "amount": 1500000, "currency": "IDR", "formatted": "IDR 1500000" },
      "total_duration_minutes": 1686,
      "layover_minutes": 1446
    }
  ]
}
```

### Filter by airline IATA code

```bash
curl -s -X POST http://localhost:8080/flights/search \
  -H 'Content-Type: application/json' \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departure_date": "2025-12-15",
    "airline_codes": ["GA", "QZ"],
    "sort_by": "best_value"
  }'
```

---

## Providers

| Provider Name | IATA Code | Notes |
|---|---|---|
| Garuda Indonesia | `GA` | Supplies city names per airport, baggage, amenities, segments |
| Lion Air | `JT` | Supplies IANA timezone names (`Asia/Jakarta`) separately |
| Batik Air | `ID` | Compact offset format (`+0700`) in datetime strings |
| AirAsia | `QZ` | RFC3339 offset (`+07:00`) embedded in datetime strings |
