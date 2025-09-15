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
    2. Middleware Stack*→ Validates method, authenticates JWT, applies RBAC and rate limiting, and logs the request.
    3. Handler → Executes business logic or queries database.
    4. Response → Returns data or error to the client.



## Modular Folder Layout

    project-root/
    ├── cmd/ # main.go entry point
    ├── internal/
    │ ├── api/ # route handlers
    │ ├── auth/ # JWT creation, verification
    │ ├── db/ # MySQL and Redis clients
    │ ├── middleware/ # logging, auth, rate limiting
    │ └── validators/ # request validation
    └── docs/ # this documentation


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
