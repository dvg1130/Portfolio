# Security Overview

    This document describes the security controls and practices implemented in the Tiered Service Backend project.



## Authentication & Session Management

    - JWT Access & Refresh Tokens
        - Access tokens are short-lived (15 minutes) and used for all protected endpoints.
        - Refresh tokens are stored in Redis, rotated on each use, and deleted on logout to prevent reuse.

    - Password Security
        - User passwords are hashed with bcrypt before storage and before login verification.

    - Token Validation
        - All protected routes verify JWT signatures and expiration.
        - Claims include `username` and `role` for role-based access control.



## Authorization & RBAC

    - Role-Based Access Control (RBAC)
        - Roles (e.g., user, tier1, tier2, admin) are embedded in JWT claims.
        - Middleware enforces endpoint access based on role.

    - Granular Rate Limiting
        - Payload-based rate limiter simulates tiered service restrictions.
        - Planned addition of request-based (token-based) rate limiting for further abuse prevention.



## Data & Transport Security

    Database Protection
        - MySQL stores only hashed passwords and minimal user data.
        - Redis keys (refresh tokens) are scoped and short-lived.

    Planned TLS Termination
        - Future production deployments will run behind a reverse proxy (e.g., Nginx or Caddy) with TLS termination to secure traffic in transit and provide hardened headers (HSTS, CSP, etc.).

## IP-Based Login Lockout

    Purpose:
        - Throttle brute-force and credential-stuffing attempts.

    Mechanism
        - Redis tracks failed logins per client IP (`fail:<ip>`).
        - After 5 failed attempts within 1 hour, a `lockout:<ip>` key blocks all further login requests from that IP for 1 hour.

    Reset:
        - Counter is cleared on successful authentication.

    Mitigation:
        - Protects against password-spraying by using the client IP rather than the username.

## Security Headers

    The application now sets several HTTP response headers to strengthen its security posture:

        Header:
        	- Purpose Current Setting

        Strict-Transport-Security:
            - Forces browsers to use HTTPS only and remember the rule for future requests.	max-age=63072000; includeSubDomains

        Content-Security-Policy:
            - Limits the sources from which scripts, images, etc. can be loaded, reducing XSS risk.	e.g. default-src 'self';
                    script-src 'self'
        X-Frame-Options:
            - Prevents click-jacking by disallowing the site to be loaded in an iframe.	DENY

        X-Content-Type-Options:
            -Blocks MIME type sniffing.	nosniff

        Referrer-Policy:
            - Controls how much referrer info is shared.	strict-origin-when-cross-origin

        Permissions-Policy:
            - Restricts powerful browser features.	geolocation=(), camera=()


## Logging & Monitoring

    - Structured Logging with Zap
        - Centralized JSON logs capture HTTP method, path, status, latency, remote IP, and authenticated user/role.
        - Supports integration with external log aggregators or SIEM tools.

    - Audit Trail
        - Login, logout, and refresh events can be monitored via log analysis for anomaly detection.


## Threat Model & Mitigations

| Threat                             | Mitigation                                                            |
|------------------------------------|-----------------------------------------------------------------------|
| Stolen refresh token               | Rotation + deletion on use/logout ensures single-use tokens.          |
| Brute-force login attempts         | Rate limiting (future request-based limiter) and bcrypt hashing.      |
| SQL injection                      | Parameterized queries in database access layer.                       |
| Man-in-the-middle (MITM) attacks   | Planned TLS termination and HSTS headers.                             |
| Privilege escalation               | Strict RBAC middleware checks using claims verified by signed JWTs.    |


## Future Enhancements

    - Deploy behind reverse proxy with TLS and automatic certificate renewal.
    - Integrate monitoring/alerting (Prometheus, Grafana, or cloud-based SIEM).
    - Optional Web Application Firewall (WAF) for additional edge protection.
