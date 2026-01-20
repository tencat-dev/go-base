# AGENT SYSTEM DESIGN

## Project Architecture Overview

This project follows the **Clean Architecture** and **Hexagonal Architecture** patterns built on top of **gofr** (github.com/gofr-dev/gofr) to ensure maintainability, testability, and independence from external frameworks or databases.

### Architecture Layers

```
┌─────────────────────┐
│   Presentation      │ ← HTTP handlers, gRPC handlers, CLI interfaces
├─────────────────────┤
│   Use Cases         │ ← Business logic, application services
├─────────────────────┤
│   Entities          │ ← Domain models, business rules
├─────────────────────┤
│   Interfaces        │ ← Repository interfaces, ports
├─────────────────────┤
│   Infrastructure    │ ← Database implementations, external services
└─────────────────────┘
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

### 3. UUID Implementation
All primary keys and unique identifiers use UUID with the default library **github.com/google/uuid**:
- Provides temporal ordering capability (for UUID v7)
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
- **Database Agent**: Implements repository interfaces with specific DB technology (default: PostgreSQL)
- **Cache Agent**: Provides caching mechanisms (default: Redis via github.com/redis/go-redis/v9)
- **Message Agent**: Handles messaging and queue operations
- **Config Agent**: Manages environment variables and configuration using github.com/caarlos0/env/v11

## Implementation Guidelines

### Entity Layer
Entities should be plain Go structs with minimal dependencies:
```go
// User represents a user entity with UUID v7 identifier
type User struct {
    ID        uuid.UUID
    CreatedAt time.Time
    UpdatedAt time.Time
    Email     string
    Name      string
}
```

### Interface Layer
Repository interfaces define contracts for data access:
```go
// UserRepository defines the contract for user data operations
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
}
```

### Infrastructure Layer
Concrete implementations adhere to interfaces while encapsulating technology-specific details:
```go
// PostgreSQLUserRepository implements UserRepository using PostgreSQL
type PostgreSQLUserRepository struct {
    db *sql.DB
}

func (r *PostgreSQLUserRepository) Save(ctx context.Context, user *User) error {
    query := `INSERT INTO users (id, email, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
    _, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Name, user.CreatedAt, user.UpdatedAt)
    return err
}
```

### Code Quality Standards
- Keep functions under 20 lines when possible (max 60 lines)
- Limit files to approximately 70 lines for optimal readability
- Follow Go idiomatic patterns and fmt standards
- Write comprehensive tests for all business logic
- Use meaningful variable and function names
- Apply consistent error handling patterns
- Eliminate dead code and unnecessary dependencies
- Apply SOLID principles for maintainable design
- Maintain >= 80% test coverage for business logic
- Prioritize unit tests for faster feedback and isolation
- Commit messages should be concise (under 60 characters when possible)
- Minimize nil pointer risks with safe pointer handling
- Ensure secure and performant pointer usage

## Benefits

1. **Technology Agnostic**: Switch databases without changing business logic
2. **Testable**: Mock interfaces for unit testing
3. **Scalable**: Independent components can scale separately
4. **Maintainable**: Clear separation of concerns
5. **Future-Proof**: Easy to adapt to new requirements or technologies