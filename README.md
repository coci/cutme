# CutMe - URL Shortener Service

A lightweight, high-performance URL shortener service built with Go, using Cassandra for persistence and Redis for distributed ID generation.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)
- [Building](#building)
- [Development](#development)
- [Project Structure](#project-structure)
- [TODO](#todo)
- [Contributing](#contributing)

## 🔍 Overview

CutMe is a URL shortening service that generates unique, short codes for long URLs using HashIDs. It follows clean architecture principles with hexagonal/ports-and-adapters design pattern.

## ✨ Features

- **URL Shortening**: Convert long URLs into short, memorable codes
- **URL Resolution**: Retrieve original URLs from short codes
- **Distributed ID Generation**: Uses Redis for atomic counter increments
- **Scalable Storage**: Cassandra for distributed data persistence
- **Configurable Hash Generation**: Customizable alphabet and minimum length
- **Structured Logging**: Using Uber Zap for production-ready logging
- **Clean Architecture**: Separation of concerns with ports and adapters pattern

## 🏗 Architecture

The project follows hexagonal architecture:

- **Core/Domain**: Business logic and entities
- **Ports**: Interfaces for external dependencies
- **Adapters**: Implementations of ports (API handlers, repositories)
- **Services**: Use case implementations
- **Infrastructure**: Configuration and external integrations

## 📦 Prerequisites

- **Go**: 1.25.1 or higher
- **Docker & Docker Compose**: For running Redis and Cassandra
- **Make**: For build automation (optional but recommended)

## 🚀 Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/coci/cutme.git
   cd cutme
   ```

2. **Install Go dependencies**:
   ```bash
   go mod download
   ```

3. **Start infrastructure services** (Redis & Cassandra):
   ```bash
   docker-compose -f docker/docker-compose.yml up -d
   ```

4. **Wait for Cassandra to be ready** (usually 30-60 seconds):
   ```bash
   docker logs cassandra -f

   ```

5. **Initialize Cassandra schema**:
   ```bash
   docker exec -it cassandra cqlsh -u cassandra -p cassandra
   ```

   Then run:
   ```cql
   CREATE KEYSPACE IF NOT EXISTS test 
   WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

   
   CREATE TABLE IF NOT EXISTS links (
       code text PRIMARY KEY,
       link text,
       created_at bigint,
       expires_at bigint
   );
   ```

## ⚙️ Configuration

### Configuration File

The main configuration is in `configs/config.yaml` file.

### Environment Variables

You can override config values using environment variables. Create a `.env` file.

## 🏃 Running the Application

### Using Make (Recommended)
Run with tests and quality checks
```shell
make run
```
Or just build and run
```shell
make build ./out/cutme -config ./configs/config.yaml
```

### Manual Run
```shell
go run cmd/main.go -config ./configs/config.yaml
```

## 🧪 Testing

### Run All Tests

```shell
make test
```

### Run Tests with Coverage

```shell

make coverage
```

### Clean Build Artifacts

```shell

make clean
```

## 🛠 Development
### Available Make Commands
```text
make check-quality - Run linting, formatting, and vetting
make lint - Run golangci-lint
make fmt - Format code
make vet - Run go vet
make tidy - Tidy go modules
make test - Run tests
make coverage - Generate coverage report
make build - Build the binary
make run - Build and run the application
make clean - Clean build artifacts
make vendor - Vendor dependencies
make all - Run all quality checks, tests, and build
```

## 📁 Project Structure

```text
cutme/
├── cmd/
│   └── main.go                 # Application entry point
├── configs/
│   └── config.yaml             # Configuration file
├── docker/
│   └── docker-compose.yml      # Docker services setup
├── internal/
│   ├── adapters/
│   │   ├── api/
│   │   │   └── shortener_handler.go     # HTTP handlers
│   │   └── repositories/
│   │       ├── cassandra_link_repository.go    # Link persistence
│   │       └── redis_id_generator_repository.go # ID generation
│   ├── core/
│   │   ├── domain/
│   │   │   ├── errors.go       # Domain errors
│   │   │   └── link.go         # Link entity
│   │   └── ports/
│   │       ├── hash_id_generator.go    # ID generator interface
│   │       ├── link_repository.go      # Repository interface
│   │       ├── logger.go               # Logger interface
│   │       └── shortener.go            # Service interface
│   ├── infra/
│   │   └── config/
│   │       ├── config.go       # Config structure
│   │       └── loader.go       # Config loader
│   ├── services/
│   │   ├── shortener_service.go  # Business logic
│   │   └── zap_logger.go         # Logger implementation
│   └── test/
│       └── shortener_service_test.go  # Service tests
├── .gitignore
├── go.mod
├── go.sum
├── makefile
└── README.md
```