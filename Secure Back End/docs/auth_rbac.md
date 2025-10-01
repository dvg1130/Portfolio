# Authentication & RBAC

    This document explains how user authentication, session management, and role-based access control (RBAC) are implemented in the Tiered Service Backend.



## Overview

    The system uses JSON Web Tokens (JWT) for stateless authentication and Role-Based Access Control (RBAC) to protect resources at different privilege levels.



## Authentication Flow

    1. User Registration & Login
        - Users register and login with a username and password.
        - Passwords are hashed with **bcrypt** before storage and when verifying login attempts.

        Lockout Logic
            - Before validating credentials, the login handler checks a Redis key for IP-based lockout.
              Five consecutive failures trigger a 1-hour block on that IP.

    2. Token Issuance
        - On successful login, the server issues:
            - **Access Token**: short-lived (15 min), signed with HS256.
            - **Refresh Token**: stored in **Redis**, expires in 7 days by default.

    3. Accessing Protected Endpoints
        - Client includes the `Authorization: Bearer <access_token>` header in requests.
        - Middleware verifies the token’s signature, expiration, and claims.

    4. Token Refresh
        - If the access token expires, the client calls `/token/refresh` with its refresh token.
        - The server validates the refresh token, issues a new access token, and rotates the refresh token (old one is deleted in Redis).

    5. Logout
        - `/logout` deletes the user’s refresh token from Redis.
        - Any attempt to use a deleted or rotated refresh token fails.



## Role-Based Access Control (RBAC)

    Role Assignment
        - Each user is assigned a role in the MySQL database (e.g., `user`, `tier1`, `tier2`, `admin`).
        - The role is embedded in JWT claims when the token is created.

    Middleware Enforcement
        - The `RequiredRole` middleware checks the role claim and ensures only authorized roles can access certain routes.
            - Roles can be adjusted in the database and take effect on the user’s next login/refresh.

    - Example

  ```go
  router.Handle("/admin",
      middleware.AuthMiddleware(
          middleware.RequiredRole("admin")(http.HandlerFunc(h.Admin)),
      ),
  )
  ```

## Security Considerations

   Single-Use Refresh Tokens
        - Rotation prevents reuse and mitigates theft.

   Bcrypt Password Hashing
        - Protects against offline password cracking.

   Context Propagation
        - username and role are attached to context.Context so downstream middleware (e.g., logging) can safely access user identity without re-verifying the token.