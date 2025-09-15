# Tiered Service Backend
> A Go API demonstrating secure authentication, RBAC, and observability.

Overview

     - A Go-based API demonstrating modern authentication, role-based access control (RBAC), and secure session management. The project implements JWT access and refresh tokens with rotation and revocation, middleware-enforced RBAC for protected endpoints, payload-based rate limiting to simulate tiered restrictions, and structured logging using Zap. User data is stored in MySQL, while refresh tokens are managed in Redis. While not a full production backend, this project showcases secure backend practices, middleware architecture, and integration of authentication, authorization, and observability features in a single, well-structured codebase.

Tech Stack

    - Language: Go
    - Database: MySQL (user auth)
    - Cache/Session Store: Redis (refresh tokens)
    - Libraries: zap, mysql, go-redis, jwt-go, net/http

Features

    - Modular Architecture
        - Clean, maintainable folder layout separating concerns: middleware, auth, db, routes, validators, and config.
        - Easy to extend with new endpoints, middleware, or services without breaking existing code.

    - Authentication & Sessions
        - JWT access + refresh tokens
        - MySQL for user credentials and roles
        - Redis for refresh token storage

    - Role-Based Access Control (RBAC)
        - Granular endpoint protection based on user roles
        - Middleware-driven for reusability

    - Rate Limiting
        - Request-based and payload-based throttling
        - Tiered control to simulate upload restrictions

    - Structured Logging
        - Zap logger middleware for clean, JSON-formatted logs
        - Captures request method, path, status, latency, and authenticated user details

    - Database Integration
        - MySQL for persistent user and role data
        - Redis for fast token management and future token-based request limiting

Project Structure

    /cmd            // main.go entry point
    /config         // configuration and environment files
    /db             // database connectors and migrations
    /middleware     // reusable middleware (logging, auth, rate limiting, etc.)
    /auth           // JWT creation, verification, context helpers
    /api            // route definitions and handlers
    /server         // server bootstrap and router wiring
    /validators     // request method and payload validation
    /docs           // security write-ups, screenshots, and code snippets

Setup & Running Locally

    1. Clone the repo

    2.Environment variables
        - Create a .env file in /config with:
            REDIS_ADDR=  //localhost:6379 default
            PORT=  //port for sql db
            DATABASE_URL=  //sql db
            DB_DRIVER=mysql //mysql for go
            JWT_SECRET_KEY= //string

    3.Run The Server
        - go run cmd/main.go


Security & Architecture Write-Ups

    See the /docs folder for:
        - Security architecture overview (Markdown with screenshots & code snippets)
        - JWT & RBAC design
        - Logging and observability strategy
        - Rate limiting configuration

    Security Controls:
        - JWT authentication, RBAC, rate-limited payloads,
        and IP-based login lockout (5 failed attempts ⇒ 1-hour block)..