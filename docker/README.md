# Docker Directory

This directory contains production-ready Dockerfiles for all services in the monorepo.

## Contents

- `Dockerfile.auth` - Authentication service
- `Dockerfile.mail` - Mail service
- `Dockerfile.mail-server` - SMTP server
- `Dockerfile.productivity` - Productivity service
- `DOCKER_BUILD.md` - Detailed build instructions

## Purpose

These Dockerfiles are designed for production builds and can:
- Access the entire monorepo
- Use the `go.work` file for proper module resolution
- Resolve local modules like `grpc`
- Create optimized, minimal production images

## Quick Start

```bash
# Build all services
docker build -f docker/Dockerfile.auth -t auth-service .
docker build -f docker/Dockerfile.mail -t mail-service .
docker build -f docker/Dockerfile.mail-server -t mail-server-service .
docker build -f docker/Dockerfile.productivity -t productivity-service .
```

For detailed instructions, see [DOCKER_BUILD.md](./DOCKER_BUILD.md). 