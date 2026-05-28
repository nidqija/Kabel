# Kabel
## An OSBAPI-Compliant Service Broker for Local Database Orchestration

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=flat-square&logo=docker)](https://www.docker.com/)

---

## 🚀 Overview

**Kabel** (Malay for "Cable") is a high-performance service broker built in Go that automates the lifecycle of development databases. By implementing the Open Service Broker API (OSBAPI) specification, Kabel allows developers to provision, bind, and manage a collection of database services—including PocketBase, PostgreSQL, and SQLite—through a single, standardized interface.

It effectively transforms a local Docker environment into a **private "Database-as-a-Service" (DBaaS) marketplace**, eliminating setup friction and enabling developers to go from a blank project to a fully provisioned, bound, and connected database in under 5 seconds.

---

## 🎯 Key Features


- **Multi-Service Support**: PocketBase, PostgreSQL, and SQLite out of the box
- **Dynamic Port Management**: Zero-conflict automatic port allocation
- **Standardized Binding**: Auto-generates connection strings and credentials
- **CLI + TUI Dashboard**: Cobra-based CLI with Bubble Tea monitoring interface
- **Idempotent Operations**: Local state management for reliability and recovery
- **Infrastructure as Code**: Bridges application development and system administration

---

## 🏗️ Core Technical Pillars

### 1. Standardized Orchestration (OSBAPI)

Instead of writing custom scripts for each database, Kabel implements the global OSBAPI standard. This decouples the "Request" (the intent to have a database) from the "Implementation" (the Docker container).

#### Key Capabilities:
- **Catalog Management**: Dynamically exposes a menu of service plans
- **Lifecycle Automation**: Handles the full `Provision → Bind → Unbind → Deprovision` flow
- **Service Instance Tracking**: Maintains state across broker restarts

### 2. Multi-Service Driver Architecture

Kabel utilizes a **Registry Pattern** in Go to manage a diverse collection of database engines. Each driver handles the specific nuances of its engine:

#### Supported Drivers:
- **PocketBase**: Automates container deployment and API port mapping
- **PostgreSQL**: Manages user creation, database isolation, and connection string generation
- **SQLite**: Orchestrates local volume mounting and file-path provisioning

### 3. Dynamic Service Binding

The project solves the "Connection String" problem. When an application "binds" to a Kabel-managed service, the broker dynamically generates credentials and network details, injecting them into the developer's environment via a custom CLI and `.env` automation.

---

## 🛠️ Technology Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Language** | Go (Golang) 1.21+ | High-performance concurrency & cloud-native tooling |
| **Orchestration** | Docker SDK for Go | Direct interaction with Docker Engine API |
| **CLI Framework** | Cobra | Command-line interface builder |
| **State Management** | Local persistence | Idempotency & recovery mechanisms |