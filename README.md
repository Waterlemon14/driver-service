# Driver Management Service
A lightweight Golang backend service handling driver creation, retrieval, and suspension.

## Project Structure
```
driver-service/
├── cmd/
│   └── server/
│       └── main.go       # Application entry point
├── internal/
│   ├── domain/           # Entities and Interfaces
│   ├── service/          # Business Logic
│   ├── adapter/          # Concrete implementations (SQLite, MemCache, ChanMQ)
│   └── handler/          # HTTP Transport
├── go.mod
└── README.md
```

## Tech Stack
- Language: Go (1.22+)
- Database: SQLite 
- Cache: In-Memory sync.Map with Time To Live (TTL)
- Message Queue: Go Channels (Built-in support for Async Events)

> [!Note]
> Caching was mocked using in-memory adapters to ensure this project runs immediately with go run without requiring Docker or environment configuration. The interfaces (DriverRepository, DriverCache, EventQueue) are designed so implementations can be swapped cleanly.
>
> Database, Caching, and Concurrency mechanisms should be hot-swappable as they are decoupled from the main implementation and you only need to implement the adapter for a specific dependency then use that adapter in main.go

## How to Run
1. Dependencies:
```
go mod init driver-service
go get github.com/mattn/go-sqlite3
go get github.com/google/uuid
go mod tidy
```

2. Start Server:
```
go run ./cmd/server
```

## API Usage
1. Create Driver
```
curl -X POST http://localhost:8080/drivers \
  -d '{"name": "John Doe", "phone": "123-456", "license_number": "L12345"}'
```

2. List Drivers (Cached)
```
curl "http://localhost:8080/drivers?page=1&limit=5"
```
> [!Note]
> Observe stdout logs to see cache hits vs DB queries. Get metrics from curl output and compare

3. Suspend Driver
```
curl -X POST http://localhost:8080/drivers/{id-from-create}/suspend \
  -d '{"reason": "Expired License"}'
```

> [!Note]
> To check return values better, you can install jq (run `brew install jq` for macOS) then run the tests above with `| jq .` at the end
> Sample `curl "http://localhost:8080/drivers?page=1&limit=5" | jq .`

## Strategy
- Caching:
  - Strategy: TTL (Time To Live) + Invalidation on Write.
  - List Endpoint: Caches the result of specific page and limit combinations for 30 seconds.
  - Invalidation: When a driver is Created or Suspended, the list cache is cleared to ensure data consistency.
- Message Queue:
  - Fire-and-forget events (driver.created, driver.suspended) are pushed to a buffered channel. A background goroutine worker consumes these logs to stdout, demonstrating async decoupling.
