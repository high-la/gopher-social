# GopherSocial API

A production-ready social feed backend built with Go, designed with clean architecture, scalability, and real-world backend patterns.

This project simulates a social platform where users can create content, follow others, and interact through a personalized feed — while handling concerns like authentication, caching, rate limiting, and observability.

---

# Project Overview

GopherSocial API is a RESTful backend system that supports:

- User registration, activation, and authentication
- Social interactions (follow users, posts, comments)
- Personalized feed generation
- Performance optimizations using caching
- Secure and scalable API design

The system is built to reflect real production backend systems, not just CRUD operations.

---

# Architecture

The project follows a Clean Layered Architecture with clear separation of concerns:

```
HTTP Layer (Handlers + Middleware)
        ↓
Application Logic (Services)
        ↓
Data Access (Store / Repository)
        ↓
PostgreSQL + Redis
```

---

### Key Architectural Decisions

- cmd/api → application entry point and HTTP layer
- internal/ → core business logic (protected from external use)
- Repository pattern (internal/store) → abstracts DB operations
- Interface-driven design → enables mocking and testing
- Context propagation → request timeouts and cancellation

# Tech Stack

- **Language:** Go
- **HTTP Server:** net/http
- **Database:** PostgreSQL
- **Cache:** Redis
- **Auth:** JWT-based authentication
- **Email Service:** SendGrid
- **Docs:** Swagger (OpenAPI)
- **Containerization:** Docker
- **CI/CD:** GitHub Actions

---

# Project Structure

```
cmd/
  api/        → HTTP handlers, middleware, routes
  migrate/    → DB migrations & seeding

internal/
  auth/       → JWT auth, token handling
  store/      → database access (users, posts, comments, followers)
  db/         → DB connection & setup
  mailer/     → email sending (SendGrid)
  ratelimiter → fixed-window rate limiting
  env/        → environment configuration

docs/         → Swagger API documentation
scripts/      → helper scripts
```

---

# Core Features

### Authentication & Authorization

- JWT-based authentication (internal/auth/jwt.go)
- Token validation middleware
- Role-based access control (RBAC)
- Secure route protection

### Database Layer

- PostgreSQL with structured queries
- Repository pattern for abstraction
- SQL migrations & seeding support
- Pagination support
- Optimized queries with joins and indexing

### Social Features

- User registration & activation (email-based)
- Follow/unfollow system
- Post creation, update, delete
- Comment system
- Personalized user feed

### Performance & Caching

- Redis caching layer (internal/store/cache)
- Cached user/profile data
- Reduced DB load for frequent reads

### Rate Limiting

- Custom fixed-window rate limiter implementation
- Protects API from abuse and traffic spikes
- Implemented in internal/ratelimiter

### Email System

- Email sending via SendGrid
- Account activation flow
- Template-based email system

### Observability & Reliability

- Structured logging
- Graceful shutdown support
- Health check endpoint
- Timeout handling via context

### Testing

- Unit tests for handlers and services
- Mocking via interfaces (mocks.go)
- Test utilities for cleaner test setup

### API Documentation

- Swagger/OpenAPI integration
- Auto-generated API docs in /docs

### API Highlights

Users

- Register & activate account
- Authenticate via JWT
- Follow other users

Posts

- Create, update, delete posts
- Add comments
- Fetch posts

Feed

- Personalized feed based on followers
- Pagination, filtering, and sorting

---

# Key Engineering Concepts Demonstrated

This project showcases:

- Clean architecture in Go
- Repository pattern & abstraction
- JWT authentication & RBAC
- Redis caching strategies
- Rate limiting implementation
- Context-aware request handling
- Testable and modular design
- Production-ready API patterns
