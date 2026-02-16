# Gocrete - Complete Documentation

## What You Get

This is the **complete, fully working Gocrete** with ALL issues fixed:

✅ All routers work (Chi, Gin, Fiber)  
✅ All databases work (PostgreSQL, MongoDB)  
✅ Docker works with Go 1.26+  
✅ All imports correct  
✅ All build errors fixed  

## Installation

```bash
# Extract the archive
tar -xzf gocrete-v5-final.tar.gz
cd gocrete-v5

# Build Gocrete
go build -o gocrete cmd/gocrete/main.go

# Optional: Install globally
sudo mv gocrete /usr/local/bin/
```

## Usage

### Command: `gocrete init`

```bash
gocrete init <project-name> --module <module-path> [options]
```

**Required:**
- `--module` - Go module path (e.g., github.com/user/project)

**Optional:**
- `--router` - chi (default), gin, or fiber
- `--db` - none (default), postgres, or mongo
- `--openapi` - none (default), gen, or manual
- `--spec` - OpenAPI spec path (required if openapi=gen)
- `--docker` - Include Docker configuration
- `--migrations` - none (default) or goose
- `--force` - Overwrite existing directory

### Examples

**1. Simple API (Chi router):**
```bash
gocrete init my-api --module github.com/user/my-api
cd my-api
go run cmd/server/main.go
```

**2. API with PostgreSQL:**
```bash
gocrete init blog-api \
  --module github.com/user/blog-api \
  --router chi \
  --db postgres \
  --migrations goose \
  --docker

cd blog-api
docker-compose up -d postgres
go run cmd/server/main.go
```

**3. High-Performance Fiber + MongoDB:**
```bash
gocrete init product-service \
  --module github.com/company/product-service \
  --router fiber \
  --db mongo \
  --docker

cd product-service
docker-compose up -d mongo
export MONGO_URL="mongodb://localhost:27017"
export MONGO_DB="product_service"
go run cmd/server/main.go
```

**4. Gin with OpenAPI:**
```bash
gocrete init api \
  --module github.com/user/api \
  --router gin \
  --db postgres \
  --openapi gen \
  --spec ./api.yaml \
  --docker

cd api
make api-gen  # Generate code from OpenAPI spec
docker-compose up
```

### Command: `gocrete add`

Add modules to existing projects:

```bash
# Add database
gocrete add db --type postgres

# Add OpenAPI
gocrete add openapi --mode gen --spec api.yaml

# Add Docker
gocrete add docker
```

## Generated Project Structure

```
your-project/
├── cmd/
│   └── server/
│       └── main.go              # Entry point with graceful shutdown
├── internal/
│   ├── config/
│   │   └── config.go            # Environment-based configuration
│   ├── logger/
│   │   └── logger.go            # Structured JSON logging
│   ├── http/
│   │   └── server.go            # HTTP server with middleware
│   ├── errors/
│   │   └── errors.go            # Error response helpers
│   ├── db/                      # (if database enabled)
│   │   ├── postgres/            # pgx connection pool + repos
│   │   └── mongo/               # mongo-driver client + repos
│   └── api/                     # (if OpenAPI enabled)
│       ├── generated/           # Generated code (don't edit)
│       └── handlers/            # Your implementations
├── migrations/                  # (if migrations enabled)
│   └── 00001_initial.sql
├── api/                         # (if OpenAPI enabled)
│   └── openapi.yaml
├── docker-compose.yml           # (if Docker enabled)
├── Dockerfile                   # (if Docker enabled)
├── Makefile                     # (if OpenAPI gen enabled)
├── .gitignore
├── .env.example
├── go.mod
└── README.md
```

## Features by Router

### Chi (Default)
```bash
gocrete init app --module github.com/user/app --router chi
```

**Pros:**
- Standard library compatible (`http.Handler`)
- Excellent middleware ecosystem
- Lightweight and fast
- Great for REST APIs

**Use when:**
- Building traditional REST APIs
- Want standard library compatibility
- Need rich middleware support

### Gin
```bash
gocrete init app --module github.com/user/app --router gin
```

**Pros:**
- Very high performance
- Large ecosystem
- JSON validation built-in
- Express-like API

**Use when:**
- Need high performance
- Want a proven, popular framework
- Coming from Node.js/Express

### Fiber
```bash
gocrete init app --module github.com/user/app --router fiber
```

**Pros:**
- Maximum performance (FastHTTP-based)
- Zero allocation in hot paths
- Express-like API
- Great for microservices

**Use when:**
- Need absolute maximum performance
- Building high-throughput services
- Want zero-allocation framework

**Note:** Fiber uses FastHTTP, not net/http (by design for performance)

## Database Options

### PostgreSQL
```bash
gocrete init app --module github.com/user/app --db postgres
```

**Includes:**
- pgx v5 connection pool
- Optimized pool settings
- Repository pattern examples
- Transaction support ready
- Docker Compose service
- Goose migration support

**Example usage:**
```go
// Generated in internal/db/postgres/repository.go
repo := postgres.NewUserRepository(db)
user, err := repo.GetByID(ctx, 1)
users, err := repo.List(ctx)
user, err := repo.Create(ctx, "user@example.com")
```

### MongoDB
```bash
gocrete init app --module github.com/user/app --db mongo
```

**Includes:**
- Official mongo-driver
- Collection-based repositories
- BSON/ObjectID support
- Aggregation pipeline ready
- Docker Compose service

**Example usage:**
```go
// Generated in internal/db/mongo/repository.go
repo := mongo.NewUserRepository(db)
user, err := repo.GetByID(ctx, "507f1f77bcf86cd799439011")
users, err := repo.List(ctx)
user, err := repo.Create(ctx, "user@example.com")
```

## OpenAPI Integration

### Code Generation Mode
```bash
gocrete init app \
  --module github.com/user/app \
  --openapi gen \
  --spec ./api.yaml
```

**Includes:**
- OpenAPI spec in `api/openapi.yaml`
- Code generation with oapi-codegen
- Separate generated/manual code
- Makefile for regeneration

**Workflow:**
1. Edit `api/openapi.yaml`
2. Run `make api-gen`
3. Implement handlers in `internal/api/handlers/`
4. Regenerate anytime with `make api-gen`

### Manual Mode
```bash
gocrete init app \
  --module github.com/user/app \
  --openapi manual
```

**Includes:**
- Example handler structure
- Router setup
- No code generation

**Workflow:**
1. Write handlers in `internal/api/handlers/`
2. Add routes in `internal/http/server.go`
3. Full manual control

## Docker

### Development
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app

# Restart service
docker-compose restart app

# Stop all
docker-compose down
```

### Building
```bash
# Build image
docker build -t my-app:latest .

# Run image
docker run -p 8080:8080 my-app:latest
```

### docker-compose.yml

Generated docker-compose includes:
- App service (multi-stage build)
- Database service (if enabled)
- Volume persistence
- Health checks
- Proper networking

## Environment Variables

Projects use environment-based configuration:

```bash
# Required
ENVIRONMENT=development
PORT=8080
LOG_LEVEL=info

# PostgreSQL (if enabled)
DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=disable

# MongoDB (if enabled)
MONGO_URL=mongodb://host:27017
MONGO_DB=database_name
```

Create `.env` file:
```bash
cp .env.example .env
# Edit .env with your values
```

Or export directly:
```bash
export DATABASE_URL="postgres://..."
go run cmd/server/main.go
```

## Testing Generated Projects

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/db/postgres/...

# Verbose
go test -v ./...

# Build check
go build ./...
```

## Common Workflows

### Workflow 1: REST API Development

```bash
# 1. Generate
gocrete init my-api \
  --module github.com/user/my-api \
  --router chi \
  --db postgres \
  --docker

# 2. Start database
cd my-api
docker-compose up -d postgres

# 3. Run migrations
goose -dir migrations postgres $DATABASE_URL up

# 4. Develop
go run cmd/server/main.go

# 5. Test
curl http://localhost:8080/health
```

### Workflow 2: Microservice Development

```bash
# 1. Generate
gocrete init user-service \
  --module github.com/company/user-service \
  --router fiber \
  --db mongo \
  --docker

# 2. Start
cd user-service
docker-compose up

# 3. Develop & Test
# Service runs in container, edit code locally
# Changes require rebuild: docker-compose up --build
```

### Workflow 3: OpenAPI-First

```bash
# 1. Create spec
cat > api.yaml << 'EOF'
openapi: 3.0.0
info:
  title: My API
  version: 1.0.0
paths:
  /api/v1/users:
    get:
      responses:
        '200':
          description: Success
EOF

# 2. Generate project
gocrete init api \
  --module github.com/user/api \
  --openapi gen \
  --spec ./api.yaml

# 3. Generate code
cd api
make api-gen

# 4. Implement handlers
# Edit internal/api/handlers/

# 5. Run
go run cmd/server/main.go
```

## Troubleshooting

### Build Errors

**"undefined: fiber"**
```bash
go mod tidy
```

**"cannot find package"**
Check your module path in imports matches `go.mod`

**Go version mismatch in Docker**
Update Dockerfile:
```dockerfile
FROM golang:alpine AS builder
ENV GOTOOLCHAIN=auto
```

### Runtime Errors

**"port already in use"**
```bash
export PORT=3000
go run cmd/server/main.go
```

**"cannot connect to database"**
Check DATABASE_URL or MONGO_URL is set correctly

**"failed to ping database"**
Make sure database is running:
```bash
docker-compose up -d postgres  # or mongo
```

## What's Fixed in v5

### From Previous Versions

- ✅ v2: Fixed go:embed paths, type conflicts
- ✅ v3: Added Go 1.26+ Docker support  
- ✅ v4: Basic Fiber router support

### New in v5

- ✅ **Complete Fiber fix**: Proper imports, no net/http
- ✅ **Conditional imports**: Clean code for each router
- ✅ **Tested all routers**: Chi, Gin, Fiber all verified working

## Project Philosophy

Gocrete is:
- ✅ A code generator (not a framework)
- ✅ A starting point (not a dependency)
- ✅ Opinionated (but not restrictive)
- ✅ Modular (only generate what you need)
- ✅ Production-ready (real-world defaults)

Gocrete is NOT:
- ❌ A framework
- ❌ A runtime dependency  
- ❌ An ORM
- ❌ A magic abstraction layer

## Best Practices

1. **Version control from day one**
   ```bash
   cd my-project
   git init
   git add .
   git commit -m "Initial commit from Gocrete"
   ```

2. **Use environment files**
   ```bash
   cp .env.example .env
   # Add .env to .gitignore (already done)
   ```

3. **Write tests early**
   ```bash
   # Add tests alongside code
   # internal/db/postgres/repository_test.go
   ```

4. **Use migrations**
   ```bash
   # Track schema changes in migrations/
   # Never modify old migrations
   ```

5. **Document your API**
   ```bash
   # Keep api/openapi.yaml up to date
   # Or use OpenAPI gen mode
   ```

## Future Development

After generating:

1. **Add your business logic** in `internal/`
2. **Add routes** in `internal/http/server.go`
3. **Add handlers** in `internal/api/handlers/`
4. **Add tests** with `_test.go` files
5. **Add middleware** as needed
6. **Add CI/CD** (GitHub Actions, etc.)

## Support & Resources

- **Documentation**: This file, `docs/` directory
- **Examples**: `examples/README.md`
- **Architecture**: `docs/ARCHITECTURE.md`
- **Contributing**: `CONTRIBUTING.md`

## License

MIT License

---

## Quick Reference

```bash
# Build Gocrete
go build -o gocrete cmd/gocrete/main.go

# Generate minimal
gocrete init app --module github.com/user/app

# Generate full
gocrete init app \
  --module github.com/user/app \
  --router chi \
  --db postgres \
  --openapi gen \
  --spec api.yaml \
  --migrations goose \
  --docker

# Add to existing
gocrete add db --type postgres
gocrete add docker

# Run generated project
cd app
go run cmd/server/main.go
```

**All routers work. All features work. Production ready.** 🚀
