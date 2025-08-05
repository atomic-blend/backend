# Docker Build Instructions

This monorepo contains multiple services that can be built using Docker. Each service has its own Dockerfile at the root level that can access the entire monorepo and use the `go.work` file for proper module resolution.

## Available Services

- **auth** - Authentication service
- **mail** - Mail service  
- **mail-server** - SMTP server
- **productivity** - Productivity service

## Building Individual Services

### Auth Service
```bash
docker build -f Dockerfile.auth -t auth-service .
```

### Mail Service
```bash
docker build -f Dockerfile.mail -t mail-service .
```

### Mail Server
```bash
docker build -f Dockerfile.mail-server -t mail-server-service .
```

### Productivity Service
```bash
docker build -f Dockerfile.productivity -t productivity-service .
```

## Building All Services

You can build all services at once using a script:

```bash
#!/bin/bash
docker build -f Dockerfile.auth -t auth-service .
docker build -f Dockerfile.mail -t mail-service .
docker build -f Dockerfile.mail-server -t mail-server-service .
docker build -f Dockerfile.productivity -t productivity-service .
```

## Key Features

- **Monorepo Support**: Each Dockerfile copies the entire monorepo and uses `go.work` for proper module resolution
- **Local Module Resolution**: Can resolve local modules like `grpc` defined in the monorepo
- **Multi-stage Builds**: Uses multi-stage builds for smaller final images
- **Alpine Base**: Uses Alpine Linux for minimal image size

## Development vs Production

For development, you can still use the individual service Dockerfiles in each directory with hot-reloading tools like `air`. The root-level Dockerfiles are optimized for production builds.

## Running Services

After building, you can run the services:

```bash
# Run auth service
docker run -p 8080:8080 auth-service

# Run mail service  
docker run -p 8080:8080 mail-service

# Run mail server
docker run -p 1025:1025 mail-server-service

# Run productivity service
docker run -p 8080:8080 productivity-service
``` 