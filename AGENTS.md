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

The system strictly follows Clean Architecture principles integrated with Domain-Driven Design (DDD):

* **Independence of frameworks**: Frameworks are replaceable details
* **Independence of UI**: Business logic does not depend on delivery mechanisms
* **Independence of database**: Domain and use cases are unaware of storage
* **Testability**: Business rules are testable without external systems
* **Domain-centric design**: Focus on domain models and business logic as the core of the system
* **DDD building blocks**: Utilize Entities, Value Objects, Aggregates, Domain Services, and Domain Events within the architecture

Dependency direction is **always inward**.

---

### 2.2 Hexagonal Architecture (Ports & Adapters)

* **Ports** define contracts the core depends on
* **Adapters** implement ports to interact with external systems
* **Domain & Application** remain isolated from infrastructure
* **DDD Integration**: Ports often represent domain interfaces that align with DDD patterns

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

## 5. Domain Rules

### 5.1 DDD Building Blocks

The system implements Domain-Driven Design building blocks within the Domain layer:

* **Entities**: Objects with distinct identities that run through lifecycles
* **Value Objects**: Objects that are distinguished by their attribute values
* **Aggregates**: Cluster of domain objects treated as a single unit
* **Domain Services**: Operations that don't naturally belong to an Entity or Value Object
* **Domain Events**: Record significant occurrences within the domain
* **Repositories**: Provide collection-like access to aggregate roots
* **Factories**: Encapsulate complex object creation logic

### 5.2 Entities

* Plain Go structs
* No framework dependencies
* No database annotations
* No infrastructure imports

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
* Migrations: `github.com/pressly/goose/v3`

```go
type PostgreSQLUserRepository struct {
    db *pgxpool.Pool
}
```

### 8.2 Dependency Injection

* DI Framework: `github.com/goforj/wire`

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
* Absolutely adhere to the declared libraries that have been pre-declared

Violations are considered **architectural defects**.

---

## 11. Code Quality Standards

* Prefer small, focused functions
* Functions must not exceed 60 lines of code
* Optimize for readability and clarity
* Follow Go idioms and `gofmt`
* Use meaningful names
* Apply SOLID principles
* Avoid premature abstraction
* Eliminate dead code
* Minimize nil pointer risks
* Code as a senior developer with 20+ years of experience
* Optimize for low-resource environments (≥ 1 CPU, 1GB RAM)

### Testing

* ≥ 80% coverage for Application & Domain layers
* Prefer unit tests over integration tests
* Always ensure unit tests are written before developing subsequent features
* Use mocks/fakes for ports

---

## 12. Benefits

* Technology-agnostic core
* Highly testable business logic
* Clear separation of concerns
* Scalable and maintainable design
* Suitable for long-term, enterprise-grade systems

---

## 13. Security

### 13.1 Security Principles

* **Defense in depth**: Multiple layers of security controls
* **Least privilege**: Components operate with minimal required permissions
* **Secure defaults**: Security configurations are enabled by default
* **Fail securely**: Systems default to secure state on failures
* **Zero trust**: Verify all requests regardless of origin

### 13.2 Authentication & Authorization

* All agents must implement proper authentication mechanisms
* Authorization checks must occur at the Application layer
* User Agent handles authentication flows and token management
* Role-based access control (RBAC) for service-to-service communication

### 13.3 Data Protection

* Sensitive data encryption at rest and in transit
* Input validation and sanitization at Presentation layer
* Secure secrets management through configuration layer
* Audit logging for sensitive operations

### 13.4 Communication Security

* All inter-service communication must use TLS
* API endpoints must implement rate limiting
* Proper CORS policies for web interfaces
* Secure headers for HTTP responses

---

## 14. Domain Organization for Microservices

### 14.1 Bounded Contexts

* Organize domains around business capabilities and responsibilities
* Each bounded context should have clear boundaries and responsibilities
* Minimize coupling between different bounded contexts
* Establish explicit interfaces for communication between contexts

### 14.2 Domain Separation Principles

* **Single Responsibility**: Each domain should have one clear purpose
* **Cohesion**: Related functionality should be grouped within the same domain
* **Autonomy**: Domains should be able to evolve independently
* **Data Ownership**: Each domain owns its data and exposes it through well-defined APIs

### 14.3 Microservices Readiness

* Design domains to be independently deployable
* Ensure domains have minimal runtime dependencies
* Define clear contract-based interfaces between domains
* Plan for eventual migration to separate services while maintaining monolithic structure initially
