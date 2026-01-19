# AGENT SYSTEM DESIGN

## Project Architecture Overview

This project follows the **Clean Architecture** and **Hexagonal Architecture** patterns to ensure maintainability, testability, and independence from external frameworks or databases.

### Architecture Layers

```
┌─────────────────┐
│   Presentation  │ ← HTTP handlers, CLI interfaces
├─────────────────┤
│   Use Cases     │ ← Business logic, application services  
├─────────────────┤
│   Entities      │ ← Domain models, business rules
├─────────────────┤
│   Interfaces    │ ← Repository interfaces, ports
├─────────────────┤
│   Infrastructure│ ← Database implementations, external services
└─────────────────┘
```

## Core Principles

### 1. Clean Architecture
- **Independence**: Framework-independent, UI-independent, database-independent
- **Testable**: Business rules can be tested without UI, database, web server, or any external element
- **Maintainable**: Changes in one area don't cascade to others

### 2. Hexagonal Architecture (Ports & Adapters)
- **Ports**: Define contracts (interfaces) that the domain exposes and consumes
- **Adapters**: Implement the ports to connect the domain with external systems
- **Domain**: Core business logic sits at the center, isolated from external concerns

### 3. UUID v7 Implementation
All primary keys and unique identifiers use **UUID version 7**:
- Provides temporal ordering capability
- Ensures global uniqueness
- Improves database performance compared to random UUIDs

### 4. Interface-Driven Development
Every data access layer implements interfaces to enable:
- Easy database switching (PostgreSQL → MySQL → MongoDB)
- Mock implementations for testing
- Dependency inversion principle compliance

## Agent Types

### Primary Agents
- **User Agent**: Manages user authentication, authorization, and profiles
- **Data Agent**: Handles data persistence and retrieval operations  
- **Service Agent**: Orchestrates business logic and use cases
- **Event Agent**: Manages event publishing and subscription

### Infrastructure Agents
- **Database Agent**: Implements repository interfaces with specific DB technology
- **Cache Agent**: Provides caching mechanisms
- **Message Agent**: Handles messaging and queue operations

## Implementation Guidelines

### Entity Layer
```go
type User struct {
    ID        uuid.UUID
    CreatedAt time.Time
    UpdatedAt time.Time
    // ... business fields
}
```

### Interface Layer
```go
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
}
```

### Infrastructure Layer
```go
// Can be swapped: PostgreSQLUserRepository → MySQLUserRepository → MongoUserRepository
type PostgreSQLUserRepository struct {
    db *sql.DB
}

func (r *PostgreSQLUserRepository) Save(ctx context.Context, user *User) error {
    // Implementation
}
```

## Benefits

1. **Technology Agnostic**: Switch databases without changing business logic
2. **Testable**: Mock interfaces for unit testing
3. **Scalable**: Independent components can scale separately
4. **Maintainable**: Clear separation of concerns
5. **Future-Proof**: Easy to adapt to new requirements or technologies