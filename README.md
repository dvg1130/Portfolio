# Security-Aware Engineering Portfolio

    Welcome to my portfolio—a collection of projects that reflect my approach to backend engineering, adversarial simulation, system resilience, and forensic analysis. These projects are designed not just to demonstrate technical skill, but to model real-world complexity, operational clarity, and ethical system design.

## Projects Overview

### 1. Secure Backend Architecture
A modular, production-grade backend system featuring:
- Hashed password authentication
- JWT access and refresh tokens with rotation and revocation
- Role-based access control (RBAC) via middleware and claims
- Payload-based rate limiting and abuse prevention
- Structured logging with Zap (middleware + semantic events)
- TLS termination strategy, containerization, and deployment notes
- SOC 2-aligned documentation and recovery workflows

 See `/src`, `/docs`, and `/deployment_maintenance` for code and documentation.



### 2. Malware Analysis & Obfuscation Toolkit
A Python + Go toolkit for simulating and analyzing obfuscated payloads:
- Payload encoder using Base64, XOR, AES, and junk injection
- Go-based sandbox API for executing binaries in Docker
- Obfuscated Python scripts mimicking malicious behavior
- Cross-language payload wrapping (PowerShell, JS, shell)
- Binary analysis using string extraction and reverse engineering tools

 See `/malware_analysis` for tooling, binaries, and threat modeling notes.



### 3. Memory Analysis Write-Up
A forensic walkthrough of capturing and analyzing a golden image:
- Memory acquisition techniques and tooling
- Static analysis using Ghidra
- Behavioral mapping and string extraction
- Reflections on volatile evidence and post-compromise visibility

 See `/forensics/memory_analysis.md`



### 4. Cloud Forensics in Simulated IaaS
An end-to-end forensic strategy in a simulated cloud environment:
- Instance compromise detection
- Log analysis and artifact correlation
- IAM abuse modeling and privilege escalation detection
- Recommendations for cloud audit readiness

 See `/forensics/cloud_forensics.md`



## Portfolio Themes

- **Security-first architecture**: Every system is designed with adversarial conditions in mind
- **Operational clarity**: Documentation spans deployment, recovery, and compliance
- **Threat simulation**: Projects model how attackers operate—and how defenders respond
- **System design**: Modular, observable, and resilient infrastructure across environments

## Future Directions

- Extend backend scaffold into a full CRUD or reward-based service
- Build reusable Go modules for security tooling and observability
- Integrate monitoring stacks (Prometheus, Grafana) into deployment strategy
- Expand forensic tooling into live response and timeline reconstruction



Built by Dwayne — backend engineer, security strategist, and systems thinker.




