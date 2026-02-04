# Driver Management Service

## Overview

Build a small **Driver Management backend service** using **Golang**.

This exercise evaluates:

- system design thinking
- clean architecture
- database usage
- caching strategy
- message queue integration
- maintainable and idiomatic Go code

The goal is to simulate a small production-style backend.  
Focus on **clarity, correctness, and simplicity** rather than feature quantity.

## Scope

You only need to implement **three APIs**.

Each API intentionally tests different backend capabilities:

| API            | Focus                      |
| -------------- | -------------------------- |
| Create Driver  | DB + Cache + Message Queue |
| List Drivers   | Pagination + DB + Cache    |
| Suspend Driver | Business logic             |

No frontend is required.

## Technical Constraints

- Language: Go
- HTTP: any framework or standard library
- Database: any (SQLite / Postgres / MySQL / etc.)
- Cache: any (Redis / in-memory / etc.)
- Message Queue: any (Kafka / NATS / RabbitMQ / simple async queue / etc.)

All technology and architecture choices are up to you.

Infrastructure can be simplified or mocked if necessary.

## Driver Entity

Minimum fields:

- id
- name
- phone
- license number
- status (active / suspended)
- suspend_reason (nullable)
- created_at

You may add additional fields if needed.

# Required APIs

## 1. Create Driver (Hard)

### Endpoint

POST /drivers

### Behavior

- create driver in database
- update or invalidate cache
- publish an event to a message queue (e.g., `driver.created`)

### Notes

- event payload can be simple
- queue implementation may be lightweight or mocked
- focus on demonstrating async/event design

## 2. List Drivers (Medium)

### Endpoint

GET /drivers?page=&limit=

### Behavior

- return paginated list of drivers
- read data from database
- implement caching for list results
- cache strategy is up to you (TTL or invalidation)

### Response should include

- list of drivers
- page
- limit
- total count
- total pages

## 3. Suspend Driver (Simple)

### Endpoint

POST /drivers/{id}/suspend

### Request Body

```json
{
  "reason": "license expired"
}
```

### Behavior

- mark driver as suspended
- store suspend reason
- persist change to database
- idempotent (safe to call multiple times)

## General Requirements

Your service should:

- use JSON for requests and responses
- return appropriate HTTP status codes
- validate inputs
- handle errors gracefully
- support concurrent requests safely
- avoid inconsistent or corrupted data
- run locally with minimal setup

## Code Quality Expectations

We evaluate:

- project structure
- separation of concerns
- readability
- idiomatic Go practices
- correct use of DB/cache/MQ
- simplicity and maintainability

Implementation details are your choice.

## Explicitly Out of Scope

- authentication/authorization
- deployment/Docker/Kubernetes
- monitoring/metrics
- frontend/UI
- features beyond the three APIs

Do not over-engineer.

## Deliverables

Please submit:

1. Source code (GitHub public repo or zip)
2. README including:
   - how to run
   - chosen tech stack
   - design decisions
   - cache strategy
   - message queue approach
   - trade-offs or assumptions

The project should start with a simple command (e.g., go run ./cmd/server).

## Time Allocation

Expected time: 1-2 hours
A clean, simple, and correct solution is preferred over a complex one.

## AI Usage

You may freely use:

- AI tools
- documentation
- online resources

We are interested in how you structure and reason about the system, not memorization.

