# Notaria 178 API

Backend REST API for a Mexican Notary Office management system. Built with Go, following Hexagonal Architecture (Ports & Adapters) and Domain-Driven Design principles.

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Technology Stack](#technology-stack)
3. [Prerequisites](#prerequisites)
4. [Environment Configuration](#environment-configuration)
5. [Database Setup](#database-setup)
6. [Running the Server](#running-the-server)
7. [Project Structure](#project-structure)
8. [Modules](#modules)
9. [API Endpoints](#api-endpoints)
10. [Authentication and Authorization](#authentication-and-authorization)
11. [Caching (Redis)](#caching-redis)
12. [Real-Time Notifications (SSE)](#real-time-notifications-sse)
13. [Cross-Module Integration](#cross-module-integration)
14. [Race Detector](#race-detector)
15. [End-to-End Tests](#end-to-end-tests)
16. [Contributing](#contributing)

---

## Architecture Overview

Each bounded context (module) is organized into three layers:

```
module/
  domain/         -- Entities, repository interfaces (ports), domain events
  app/            -- Use cases, DTOs, business rules
  infra/          -- PostgreSQL repositories, Gin controllers, routes, dependency wiring
```

- **Domain Layer**: Pure Go, zero framework imports. Defines entities and repository interfaces.
- **Application Layer**: Orchestrates business logic through use cases. Consumes domain ports.
- **Infrastructure Layer**: Implements ports with concrete adapters (PostgreSQL, Redis, BCrypt, JWT, SSE).

Cross-module communication is achieved exclusively through domain-level interfaces and adapter bridges located in `internal/integration/adapters/`. No module imports another module's domain directly.

---

## Technology Stack

| Component          | Technology                       |
|--------------------|----------------------------------|
| Language           | Go 1.25.5                        |
| HTTP Framework     | Gin v1.11.0                      |
| Database           | PostgreSQL (lib/pq v1.11.2)      |
| Cache              | Redis (go-redis/v9 v9.18.0)      |
| Authentication     | JWT (golang-jwt/v4)              |
| Password Hashing   | BCrypt (golang.org/x/crypto)     |
| CORS               | gin-contrib/cors                 |
| UUID               | google/uuid v1.6.0               |
| Environment        | godotenv v1.5.1                  |

---

## Prerequisites

- Go 1.25.5 or higher
- PostgreSQL 14 or higher
- Redis 7 or higher (optional, the server starts without it)
- Git

---

## Environment Configuration

Create a `.env` file in the project root:

```env
# Required
JWT_SECRET=your_secure_jwt_secret
DATABASE_URL=postgres://user:password@localhost:5432/notaria178_db?sslmode=disable

# Optional
PORT=8080
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

| Variable         | Required | Description                                    |
|------------------|----------|------------------------------------------------|
| `JWT_SECRET`     | Yes      | Secret key for JWT token signing               |
| `DATABASE_URL`   | Yes      | PostgreSQL connection string                   |
| `PORT`           | No       | HTTP server port (default: 8080)               |
| `REDIS_ADDR`     | No       | Redis address (host:port). If empty, no cache  |
| `REDIS_PASSWORD` | No       | Redis password (empty for no auth)             |
| `REDIS_DB`       | No       | Redis database index (default: 0)              |

---

## Database Setup

Execute the schema against your PostgreSQL instance:

```bash
psql -U your_user -f schema.sql
```

The schema creates:
- Enums: `user_role`, `user_status`, `work_status`, `document_category`, `notification_type`
- 12 tables: `branches`, `users`, `attendances`, `clients`, `act_catalogs`, `works`, `work_acts`, `work_collaborators`, `documents`, `work_comments`, `notifications`, `audit_logs`

---

## Running the Server

```bash
# Install dependencies
go mod tidy

# Run
go run .

# Or build and run
go build -o notaria178 .
./notaria178
```

The server starts at `http://localhost:8080` (or the port defined in `PORT`).

---

## Project Structure

```
Notaria178_API/
  main.go                          -- Entry point and dependency wiring
  go.mod / go.sum                  -- Go module definition
  schema.sql                       -- PostgreSQL DDL
  run_with_race_detector.sh        -- Race detector script (Linux/macOS)
  run_with_race_detector.bat       -- Race detector script (Windows)
  tests/
    e2e/                           -- Python end-to-end tests
  internal/
    core/
      postgresql.go                -- DB pool (10 open, 5 idle connections)
      cors.go                      -- CORS middleware configuration
      cache/
        redis_adapter.go           -- CachePort interface + Redis implementation
      dtos/
        pagination.go              -- PaginationRequest, DateRangeRequest, PaginatedResponse
      domain/ports/
        cache_repository.go        -- Legacy cache interface (unused)
    middleware/
      auth.go                      -- JWT extraction + role-based authorization
    integration/
      adapters/
        audit_adapter.go           -- AuditLogger port -> audit use case
        notification_adapter.go    -- Notifier port -> notification use case
    user/                          -- User management, login, profiles
    attendance/                    -- Employee check-in/check-out
    act/                           -- Act catalog management (cached)
    client/                        -- Client registry
    branch/                        -- Branch/office management
    work/                          -- Work dossiers (expedientes), status machine
    document/                      -- File upload/download with OS-aware storage
    notification/                  -- Notifications + SSE real-time streaming
    audit/                         -- Audit log recording and search
```

---

## Modules

### User
Employee and admin management. Handles registration, login (JWT), profile viewing and editing.

### Attendance
Clock-in/clock-out system for employees. Supports multiple shifts per day.

### Act (Cached)
Catalog of notarial act types (e.g., Sale, Power of Attorney). First module with Redis cache — search results are cached for 24 hours with automatic invalidation on create, update, or status toggle.

### Client
Registry of notary clients with name, RFC, phone, and email.

### Branch
Office/branch management. Users are assigned to a branch.

### Work
Core module. Manages legal dossiers (expedientes) with a status machine:
`PENDING -> IN_PROGRESS -> READY_FOR_REVIEW -> APPROVED / REJECTED`

Supports collaborators, comments, and cross-module integration with audit logging and notifications on status changes.

### Document
File uploads tied to works and clients. OS-aware storage paths (Linux/Windows). Supports categories: DRAFT_DEED, FINAL_DEED, CLIENT_REQUIREMENT, OTHER.

### Notification
User notifications with real-time delivery via Server-Sent Events (SSE). Types: NEW_COMMENT, ASSIGNMENT, STATUS_CHANGE, SYSTEM.

### Audit
Immutable audit trail. Records user, action, entity, entity ID, and JSONB details (before/after snapshots). Search endpoint restricted to admin roles.

---

## API Endpoints

### User `/api/v1/users`

| Method | Path                       | Description         | Auth          |
|--------|----------------------------|---------------------|---------------|
| POST   | `/api/v1/users/login`      | Login (returns JWT) | Public        |
| GET    | `/api/v1/users/profile`    | Get own profile     | Authenticated |
| PATCH  | `/api/v1/users/profile`    | Update own profile  | Authenticated |
| POST   | `/api/v1/users/create`     | Create employee     | Admin         |
| PATCH  | `/api/v1/users/update/:id` | Update employee     | Admin         |
| GET    | `/api/v1/users/search`     | Search employees    | Admin         |

### Attendance `/api/v1/attendance`

| Method | Path                                  | Description              | Auth          |
|--------|---------------------------------------|--------------------------|---------------|
| POST   | `/api/v1/attendance/check`            | Clock in / clock out     | Authenticated |
| GET    | `/api/v1/attendance/history`          | Own attendance history   | Authenticated |
| GET    | `/api/v1/attendance/admin/history/:id`| Employee attendance      | Admin         |

### Act `/api/v1/acts`

| Method | Path                       | Description         | Auth          |
|--------|----------------------------|---------------------|---------------|
| GET    | `/api/v1/acts/search`      | Search acts (cached)| Authenticated |
| POST   | `/api/v1/acts/create`      | Create act          | Admin         |
| PATCH  | `/api/v1/acts/update/:id`  | Update act          | Admin         |
| PATCH  | `/api/v1/acts/status/:id`  | Toggle act status   | Admin         |

### Client `/api/v1/clients`

| Method | Path                          | Description      | Auth          |
|--------|-------------------------------|------------------|---------------|
| GET    | `/api/v1/clients/search`      | Search clients   | Authenticated |
| POST   | `/api/v1/clients/create`      | Create client    | Authenticated |
| PATCH  | `/api/v1/clients/update/:id`  | Update client    | Authenticated |

### Branch `/api/v1/branches`

| Method | Path                           | Description      | Auth          |
|--------|--------------------------------|------------------|---------------|
| GET    | `/api/v1/branches/search`      | Search branches  | Authenticated |
| POST   | `/api/v1/branches/create`      | Create branch    | Admin         |
| PATCH  | `/api/v1/branches/update/:id`  | Update branch    | Admin         |

### Work `/api/v1/works`

| Method | Path                                          | Description           | Auth          |
|--------|-----------------------------------------------|-----------------------|---------------|
| GET    | `/api/v1/works/search`                        | Search works          | Authenticated |
| GET    | `/api/v1/works/:id`                           | Get work detail       | Authenticated |
| POST   | `/api/v1/works/create`                        | Create work           | Managers      |
| PATCH  | `/api/v1/works/update/:id`                    | Update work           | Managers      |
| PATCH  | `/api/v1/works/status/:id`                    | Update work status    | Managers      |
| GET    | `/api/v1/works/:id/comments`                  | List work comments    | Authenticated |
| POST   | `/api/v1/works/:id/comments`                  | Add comment           | Authenticated |
| POST   | `/api/v1/works/:id/collaborators`             | Add collaborator      | Managers      |
| DELETE | `/api/v1/works/:id/collaborators/:userId`     | Remove collaborator   | Managers      |

### Document `/api/v1/documents`

| Method | Path                                   | Description          | Auth          |
|--------|----------------------------------------|----------------------|---------------|
| POST   | `/api/v1/documents/upload`             | Upload document      | Authenticated |
| GET    | `/api/v1/documents/work/:work_id`      | List work documents  | Authenticated |
| GET    | `/api/v1/documents/download/:id`       | Download document    | Authenticated |

### Notification `/api/v1/notifications`

| Method | Path                                   | Description               | Auth          |
|--------|-----------------------------------------|--------------------------|---------------|
| GET    | `/api/v1/notifications`                | List my notifications     | Authenticated |
| GET    | `/api/v1/notifications/stream`         | SSE real-time stream      | Authenticated |
| PATCH  | `/api/v1/notifications/:id/read`       | Mark as read              | Authenticated |
| PATCH  | `/api/v1/notifications/read-all`       | Mark all as read          | Authenticated |

### Audit `/api/v1/audit`

| Method | Path                       | Description          | Auth  |
|--------|----------------------------|----------------------|-------|
| GET    | `/api/v1/audit/search`     | Search audit logs    | Admin |

**Total: 30 endpoints (1 public, 29 protected)**

---

## Authentication and Authorization

All protected endpoints require a `Bearer` token in the `Authorization` header:

```
Authorization: Bearer <jwt_token>
```

The JWT payload contains `userID`, `userRole`, and `branchID`, which are extracted by the auth middleware and injected into the Gin context.

### Roles

| Role          | Description                                     |
|---------------|-------------------------------------------------|
| SUPER_ADMIN   | Full system access across all branches           |
| LOCAL_ADMIN   | Full access within their assigned branch         |
| DRAFTER       | Can manage works, documents, and comments        |
| DATA_ENTRY    | Basic data input capabilities                    |

Role checks are enforced at the route level using `middleware.RequireRoles(...)`.

---

## Caching (Redis)

Redis is an **optional** dependency. If `REDIS_ADDR` is not set or Redis is unreachable, the server starts normally and all queries go directly to PostgreSQL.

### Cache Architecture

The cache layer follows the Ports & Adapters pattern:

- **Port**: `cache.CachePort` interface (Set, Get, Invalidate, InvalidatePrefix)
- **Adapter**: `cache.RedisCache` using go-redis/v9

### Cache Strategy (Act Module)

| Operation      | Behavior                                                    |
|----------------|-------------------------------------------------------------|
| Search acts    | Check Redis first (key: `acts:search:<hash>`). On miss, query DB and cache result with 24h TTL. |
| Create act     | Invalidate all keys with prefix `acts:search:`.             |
| Update act     | Invalidate all keys with prefix `acts:search:`.             |
| Toggle status  | Invalidate all keys with prefix `acts:search:`.             |

All Redis errors are handled gracefully — on failure, the operation falls back to the database without interrupting the request.

---

## Real-Time Notifications (SSE)

The notification module includes a Server-Sent Events hub for real-time delivery:

1. Client connects to `GET /api/v1/notifications/stream` with a valid JWT.
2. The SSEHub registers the client channel keyed by user ID.
3. When a notification is created (e.g., from a work status change), `CreateNotificationUseCase` calls `Broadcast()`.
4. The SSEHub pushes the event to the connected client's channel.
5. On disconnect, the client channel is removed and closed.

The SSEHub uses `sync.RWMutex` for concurrent access safety and buffered channels with non-blocking sends to prevent slow clients from blocking the hub.

---

## Cross-Module Integration

The `work` module triggers audit logs and notifications when a work status changes. This is achieved without circular imports:

```
work/domain/events/
  external_ports.go       -- Defines AuditLogger and Notifier interfaces

internal/integration/adapters/
  audit_adapter.go        -- Bridges AuditLogger to audit.LogActionUseCase
  notification_adapter.go -- Bridges Notifier to notification.CreateNotificationUseCase
```

Adapters are created in `main.go` and injected into the work module at startup.

---

## Race Detector

Go's built-in race detector can be used to verify concurrent safety (particularly for the SSE hub):

**Linux / macOS:**
```bash
chmod +x run_with_race_detector.sh
./run_with_race_detector.sh
```

**Windows:**
```cmd
run_with_race_detector.bat
```

These scripts compile and run the binary with `-race` enabled. Any data race will be reported to stderr at runtime.

---

## End-to-End Tests

A Python-based E2E test suite validates the complete API flow.

### Setup

```bash
pip install -r requirements-qa.txt
```

### Run

```bash
pytest tests/e2e/ -v
```

### Features

- 6 ordered tests covering login, CRUD, and cross-module flows
- Automatic PDF report generation (via conftest.py hook)
- Results saved to `tests/e2e/` directory

---

## Contributing

1. Follow the hexagonal architecture pattern established in existing modules.
2. Domain and application layers must not import infrastructure packages.
3. Use `context.Context` in all repository methods and database calls.
4. Use inline role string literals in the application layer instead of importing user entities.
5. Run `go build ./...` before committing to verify compilation.
6. Run with the race detector periodically to catch concurrency issues.
