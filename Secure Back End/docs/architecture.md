# Architecture Overview

    This document describes the overall design and modular structure of the Tiered Service Backend.


## High-Level Design

    The application is built as a modular Go service with clearly separated responsibilities:

        - API Layer – Exposes HTTP endpoints via the Go `net/http` standard library.
        - Middleware Layer – Handles authentication, logging, rate limiting, and request validation.
        - Auth & RBAC – Issues and verifies JWT access/refresh tokens and enforces role-based access control.
        - Persistence
        - MySQL for user data, roles, and credentials (passwords stored as bcrypt hashes).
        - Redis for refresh-token rotation and fast key-value storage.



## Data Flow

    1. Client Request → Hits API route.
    2. Middleware Stack*→ Validates method, authenticates JWT, applies RBAC and limiters, and logs the request.
    3. Handler → Executes business logic or queries database.
    4. Response → Returns data or error to the client.



## Modular Folder Layout

    project-root/
├── cmd/
│   └── main.go                # Entry point of the app
│
├── config/                    # Configuration files/env loaders
│
├── internal/                  # Core application logic (not exposed externally)
│   ├── api/                   # Handlers for routes (snake, breeding, admin, etc.)
│   ├── auth/                  # JWT, password hashing, refresh rotation
│   ├── helpers/               # Utility funcs (e.g., DecodeBody, ClientIP)
│   ├── middleware/            # Logging, security headers, limiters, validators
│   ├── server/                # Server setup, router wiring, middleware wrapping
│   └── validators/            # Request validation logic
│
├── logs/                      # Centralized logging (zap setup, LogEvent helpers)
│
├── models/                    # Structs for requests/responses (User, Snake, BreedingEvent, etc.)
│
├── repo/                      # Data persistence layer
│   ├── auth_db/                    # Database connection setup (auth_db, data_db, redis)
│   └── data_db/               # SQL queries (authdb, datadb, migrations)
│   └── redi_db_/
|
├── go.mod
└── go.sum



This structure makes it easy to extend or swap components (e.g., adding a new DB client or middleware) without touching unrelated code.



## Diagram

See the included **`images/architecture-diagram.png`** for a visual representation of how the API, middleware, MySQL, and Redis layers interact.



### Key Design Principles

    - Separation of Concerns – Each package handles a single responsibility.
    - Security First – JWT-based session management, RBAC, and rate limiting are baked in.
    - Observability – Centralized structured logging with Zap for easier debugging and monitoring.


## Deployment & Security Notes

    Planned production deployments will run behind a reverse proxy such as **Nginx** or **Caddy** with **TLS termination**.
    The proxy will:

        - Handle HTTPS/TLS encryption and automatic certificate renewal (e.g., via Let’s Encrypt).
        - Forward only validated traffic to the Go application.
        - Provide an additional security layer for rate limiting, caching, and header hardening.

    This design keeps the Go service focused on application logic while the proxy manages SSL/TLS handshakes and edge-security concerns.

_This document provides the foundation for understanding how all moving parts fit together and where to add new features or enhancements._
