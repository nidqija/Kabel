# Kabel
## An OSBAPI-Compliant Service Broker for Local Database Orchestration

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=flat-square&logo=docker)](https://www.docker.com/)

---

## 🚀 Overview

**Kabel** is a service broker written in Go that automates provisioning and binding of local development databases. It implements the Open Service Broker API (OSBAPI) and manages services such as PocketBase and PostgreSQL using Docker containers.

The project is focused on developer ergonomics: fast provisioning, automatic credential generation, and a consistent CLI for lifecycle operations.

---

## 🎯 Key Features

- Multi-service support: PocketBase, PostgreSQL, and SQLite (via embedded drivers)
- Automatic port and container lifecycle management using Docker
- CLI powered by `cobra` for provisioning, binding, and instance management
- Idempotent local state to survive restarts and enable safe retries

---

## 🏗️ Core Technical Pillars

- OSBAPI-compliant broker to decouple requests from implementations
- Driver/strategy pattern to support multiple database engines (PocketBase, Postgres, SQLite)
- Dynamic binding that produces ready-to-use connection strings and credentials

---

## 🛠️ Technology Stack

| Component | Technology | Notes |
|-----------|-----------|-------|
| Language | Go 1.25.0 | See `go.mod` for module versions |
| Orchestration | Docker SDK for Go (`github.com/docker/docker`) | Interacts with Docker Engine |
| Database Engines | PocketBase (`github.com/pocketbase/pocketbase`), PostgreSQL (`lib/pq`), SQLite (via `modernc.org/sqlite`) | PocketBase used for lightweight local app DBs |
| CLI | Cobra (`github.com/spf13/cobra`) | Command-line tooling and subcommands |

---

## Project Layout

- `main.go` — application entrypoint
- `go.mod` — Go module file (Go 1.25)
- `Command/` — CLI commands and wiring (Cobra)
- `Engine/` — engine abstractions and Docker interaction
- `Strategy/` — strategy implementations for PocketBase, Postgres, etc.
- `data/my-pb-app/` — sample PocketBase app schema and `types.d.ts`

---

## Quick Start (Development)

Prerequisites:

- Go 1.25+ installed
- Docker Engine running and accessible to your user

Build and run locally:

```powershell
go build -v ./...
./kabel
```

Or run directly:

```powershell
go run main.go
```

Common commands (via the CLI):

- `kabel catalog` — list available service plans
- `kabel provision <plan>` — provision a service instance
- `kabel bind <instance-id>` — bind and receive credentials

Refer to the `Command/` directory for command implementations.

---

## Development Notes

- Configuration and state are stored locally; be careful when removing state files as they affect idempotency.
- PocketBase artifacts and types live in `data/my-pb-app/` — useful as an example app
- Module versions and indirect deps are stored in `go.mod`/`go.sum`

---

## License

MIT — see `LICENSE`.
