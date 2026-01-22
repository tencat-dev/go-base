# AGENT SYSTEM ARCHITECTURE CONSTITUTION

## 1. Purpose

This document defines the **official architecture rules and design principles** of the project.
It serves as a **single source of truth** for all contributors to ensure:

* Architectural consistency
* Long-term maintainability
* High testability
* Technology independence

This project is built on **Clean Architecture** and **Hexagonal Architecture (Ports & Adapters)**, implemented using **go-kratos** as the application framework.

---

## 2. Architectural Philosophy

### 2.1 Clean Architecture

The system strictly follows Clean Architecture principles:

* **Independence of frameworks**: Frameworks are replaceable details
* **Independence of UI**: Business logic does not depend on delivery mechanisms
* **Independence of database**: Domain and use cases are unaware of storage
* **Testability**: Business rules are testable without external systems

Dependency direction is **always inward**.

---

### 2.2 Hexagonal Architecture (Ports & Adapters)

* **Ports** define contracts the core depends on
* **Adapters** implement ports to interact with external systems
* **Domain & Application** remain isolated from infrastructure

---

## 3. Logical Architecture Layers

```
┌───────────────────────────┐
│ Presentation              │  HTTP / gRPC / CLI
├───────────────────────────┤
│ Application               │  Use Cases / Services
├───────────────────────────┤
│ Domain                    │  Entities / Value Objects
├───────────────────────────┤
│ Ports                     │  Repositories / Event / External APIs
├───────────────────────────┤
│ Infrastructure            │  DB / Cache / MQ / Frameworks
└───────────────────────────┘
```

### Layer Responsibilities

| Layer          | Responsibility                                       |
| -------------- | ---------------------------------------------------- |
| Presentation   | Request validation, DTO mapping, response formatting |
| Application    | Orchestrates business rules and workflows            |
| Domain         | Core business rules and invariants                   |
| Ports          | Interfaces owned by the core                         |
| Infrastructure | Technology-specific implementations                  |

---

## 4. Agent Concept (Conceptual Only)

> **Agent** represents a *conceptual responsibility*, **not** a runtime process, package, or goroutine.

Agents must always be implemented **within the boundaries of the architecture layers**.

### Agent Mapping

| Agent          | Responsibility                                | Architectural Location |
| -------------- | --------------------------------------------- | ---------------------- |
| User Agent     | Authentication, authorization, user lifecycle | Application / Ports    |
| Service Agent  | Business orchestration                        | Application            |
| Data Agent     | Data access contracts                         | Ports                  |
| Event Agent    | Event publishing & consuming contracts        | Ports                  |
| Database Agent | Database implementation                       | Infrastructure         |
| Cache Agent    | Caching implementation                        | Infrastructure         |
| Message Agent  | Messaging / queue implementation              | Infrastructure         |
| Config Agent   | Configuration & bootstrap                     | Infrastructure         |

---

## 5. Domain Rules

### 5.1 Entities

* Plain Go structs
* No framework dependencies
* No database annotations
* No infrastructure imports

```go
type User struct {
    ID        uuid.UUID
    Email     string
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

---

### 5.2 UUID Policy

* All primary identifiers use **UUID v7**
* UUIDs are generated **in the Application layer**
* IDs are immutable once assigned
* Domain must not rely on database-generated identifiers

```go
type IDGenerator interface {
    New() uuid.UUID
}
```

---

## 6. Application Layer Rules

* Contains use cases and business workflows
* Depends only on:

    * Domain
    * Ports
* Never depends on Infrastructure

```go
type UserUseCase struct {
    repo UserRepository
    ids  IDGenerator
}
```

---

## 7. Ports (Interface Layer)

Ports are **owned by the core**, not infrastructure.

```go
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
}
```

---

## 8. Infrastructure Rules

Infrastructure is a **detail**, never a dependency of the core.

### 8.1 Database

* Default database: **PostgreSQL 18+**
* Driver: `github.com/jackc/pgx/v5`
* ORM / SQL builder: `github.com/stephenafamo/bob`
* Migrations: `github.com/golang-migrate/migrate/v4`

```go
type PostgreSQLUserRepository struct {
    db *pgxpool.Pool
}
```

---

### 8.2 Migration Structure

```
migrations/
├── 00001_create_users_table.up.sql
├── 00001_create_users_table.down.sql
├── 00002_add_index_email.up.sql
├── 00002_add_index_email.down.sql
```

---

### 8.3 Cache & Messaging

* Cache: Redis 8+ or Valkey
* Client: `github.com/redis/go-redis/v9`
* Messaging: Redis-based queue or stream

---

## 9. Configuration

* Configuration is loaded at bootstrap
* Uses go-kratos config system
* Environment-specific overrides allowed
* Domain and Application layers must never read env vars directly

---

## 10. Architectural Rules (Mandatory)

The following rules are **non-negotiable**:

* Domain MUST NOT import:

    * Infrastructure
    * Frameworks
    * Database drivers
* Application MUST NOT import Infrastructure
* Infrastructure MUST NOT be referenced by Domain or Application
* Cross-module communication must occur via Ports or Events

Violations are considered **architectural defects**.

---

## 11. Code Quality Standards

* Prefer small, focused functions
* Optimize for readability and clarity
* Follow Go idioms and `gofmt`
* Use meaningful names
* Apply SOLID principles
* Avoid premature abstraction
* Eliminate dead code
* Minimize nil pointer risks
* Optimize for low-resource environments (≥ 1 CPU, 1GB RAM)

### Testing

* ≥ 80% coverage for Application & Domain layers
* Prefer unit tests over integration tests
* Use mocks/fakes for ports

---

## 12. Benefits

* Technology-agnostic core
* Highly testable business logic
* Clear separation of concerns
* Scalable and maintainable design
* Suitable for long-term, enterprise-grade systems
