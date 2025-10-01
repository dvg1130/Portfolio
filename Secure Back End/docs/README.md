Directory Structure

    docs/
    │
    ├─ README.md                  # Quick index of the docs folder
    ├─ architecture.md            # System architecture overview
    ├─ security.md                # Security controls & threat model
    ├─ auth_rbac.md               # JWT, refresh tokens, RBAC design
    ├─ logging.md                 # Zap logging & observability
    ├─ rate_limiting.md           # Request & payload rate limiter design
    └─ images/
    ├─ architecture-diagram.png
    ├─ db-schema.png
    ├─ flow-auth-sequence.png
    └─ log-sample.png

This folder contains all technical and security documentation for the **Tiered Service Backend** project.

Use the files below to explore different aspects of the system:

- **[architecture.md](architecture.md)** – High-level system overview and modular layout diagram.
- **[security.md](security.md)** – Security controls, threat model, and mitigations.
- **[auth_rbac.md](auth_rbac.md)** – Detailed explanation of the JWT authentication flow, refresh token rotation, and RBAC.
- **[logging.md](logging.md)** – Zap logging configuration, sample log output, and observability notes.
- **[rate_limiting.md](rate_limiting.md)** – Request and payload rate-limiting design.

The included layout diagram in `architecture.md` illustrates how the Go API, middleware layers, MySQL, and Redis components fit together.
Each document may reference screenshots or code snippets stored in the `images/` subdirectory.
