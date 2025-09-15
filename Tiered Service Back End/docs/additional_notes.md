# Additional Notes

This document collects architectural considerations, deployment ideas, and security strategies that complement the core backend design.

## Separate Database for Non-Auth User Data

    Rationale:
       - Splitting authentication data (credentials, tokens) from general user data (profiles, preferences, analytics) reduces the blast radius of a compromise and simplifies scaling.

    Approach:

       - Primary Auth DB: Stores user credentials, roles, and refresh tokens.

       - Secondary App Data DB: Stores all non-auth data.

    Benefits:

       - Easier maintenance and backup schedules.

       - Stronger security boundaries and targeted access control.

       - Containerization & Production Deployment

    Options:

       - Docker for image-based deployments.

       - Docker Compose or Kubernetes for orchestration and scaling.

    Best Practices:

       - Use minimal base images (e.g., distroless or Alpine) to reduce attack surface.

       - Employ multi-stage builds to keep final images small.

       - Ensure environment variables (secrets, keys) are injected securely at runtime.

    Reverse Proxy & TLS Termination:
       - Use Nginx or Traefik in front of the application to handle HTTPS/TLS termination.

## Client-Side Validation

    Purpose:

       - Prevent obvious malformed requests and reduce load on the server.

    Implementation Examples:

       - Form field checks (length, type, regex) before submission.

       - Real-time feedback to users to catch errors early.

    Limitations:

       - Never a replacement for server-side validation.

       - Acts as a first defense layer to limit automated or accidental misuse.

## Vulnerability Mitigation

       - Considered Techniques:

       - Input validation and sanitation to prevent SQL/NoSQL injection and XSS.

       - Rate limiting to deter brute-force attacks and abuse.

       - Proper error handling to avoid information leaks.

    Current Usage:

       - - Implemented role-based access control (RBAC) and strict JWT validation.

       -  Logging and monitoring to detect suspicious activity.

       - Parameterized sql queries

    Available Options for Future Hardening:

       - Web Application Firewall (WAF).

       - Security headers (e.g., CSP, HSTS).

       - Automated dependency vulnerability scanning (e.g., Dependabot).

## Logging & Production Handling

   During development and testing, logs are written to the terminal (stdout) for simplicity.
   In a production deployment, the zap logger can be configured to:

      Write to Rotating Log Files:
         - Use zap’s file-based core (or a log rotation utility such as `logrotate`) to store structured JSON logs on disk.
      Forward to a Log Server:
         - Ship logs to a centralized service (e.g., ELK stack, Loki/Grafana, or a SIEM) for aggregation, search, and alerting.

      Separation of Concerns:
         - Logs will **not** be stored in the authentication database or Redis.
         - Instead, they should be persisted in a **dedicated logging database or service** to prevent access conflicts and to support long-term retention and   analysis.

      Handling Strategy:
         Two recommended approaches:
            1. Scheduled Task – Periodically move/rotate local log files to long-term storage.
            2. Real-Time Streaming – Send logs directly to a log server or message queue (e.g., Kafka, Fluent Bit) for immediate ingestion.

      This ensures:
         - Integrity and tamper resistance of logs.
         - Centralized monitoring and alerting.
         - Compliance with production best practices.


## Other Considerations

    Backups & Disaster Recovery:
       - Automate database backups and periodically test restoration.

    Observability:
       - Combine metrics, tracing, and logging for full-stack visibility.

    CI/CD Pipeline:
       - Integrate automated tests and security scans before deployment.