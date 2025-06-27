# Go Code Runner Sandbox

This project demonstrates how to safely execute untrusted Go code using Docker-based sandboxing techniques. It is intended for developers and security engineers who want to explore container isolation, runtime constraints, and secure execution environments.

## Overview

This tool reads a hardcoded Go source file from the local filesystem, then executes it in a containerized, hardened Docker environment using the official `golang:latest` image. Output from the execution is returned via a simple HTTP API.

The project is useful as a foundation for:
- Safe dynamic analysis of unknown or untrusted code
- Code runner microservices
- Teaching container hardening techniques
- Red/blue team labs involving sandbox bypasses or container escapes

## Features

- Runs Go code inside a temporary container
- Enforces strict isolation:
  - Read-only root filesystem
  - No network access
  - CPU and memory limits
  - Process count limits
  - All Linux capabilities dropped
  - Limited writable path (`/tmp`) with `--tmpfs`
- Adjustable request timeout for container execution
- Accepts only a hardcoded absolute file path to the `.go` file
- Returns standard output and error output via HTTP response

## Architecture

1. Go web server exposes a `/run` POST endpoint
2. When triggered, the server mounts a local Go source file into a Docker container
3. Container executes the file using `go run`, with runtime constraints
4. Output is captured and returned to the requester

## Usage

1. Ensure Docker is installed and running
2. Adjust swap or memory if encountering container kills (2 GB recommended)
3. Create a test Go file in seperate directory :

    ```go
    // snippet.go
    package main

    import "fmt"

    func main() {
        fmt.Println("Hello from inside Docker!")
    }

4. Replace with your path to file:

    ```go
    snippetPath := "/PATH/TO/FILE/snippet.go"

5. Run the sandbox and start server
    - go run .
    or
    - go run sandbox.go

6. In seperate terminal POST curl request to http://localhost:8080/run
    -  curl http://localhost:8080/run -X POST


## Snippets
    - Three test files are included to demonstrate different use cases and execution behaviors. See snippets.txt for a summary of each test and its purpose.
