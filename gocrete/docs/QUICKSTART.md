# Gocrete Quick Start Guide

This guide will get you up and running with Gocrete in minutes.

## Installation

### From Source (Recommended for now)

```bash
# Clone the repository
git clone https://github.com/TRiZKy/gocrete.git
cd gocrete

# Build the binary
go build -o gocrete cmd/gocrete/main.go

# Install to your PATH
sudo mv gocrete /usr/local/bin/

# Verify installation
gocrete --help
```

### Future: Go Install (Once Published)

```bash
go install github.com/TRiZKy/gocrete/cmd/gocrete@latest
```

## Your First Project

### 1. Create a Simple API

```bash
gocrete init my-api --module github.com/yourusername/my-api
cd my-api
```

This creates a minimal Go project with:
- HTTP server (Chi router)
- Structured logging
- Configuration management
- Health check endpoints

### 2. Start the Server

```bash
go mod download
go run cmd/server/main.go
```

Output:
```
{"time":"...","level":"INFO","msg":"Starting server","port":8080,"env":"development"}
{"time":"...","level":"INFO","msg":"Server listening","addr":":8080"}
```

### 3. Test It

```bash
# Health check
curl http://localhost:8080/health
{"status":"healthy"}

# Readiness check
curl http://localhost:8080/ready
{"status":"ready"}
```

## Adding a Database

### PostgreSQL

```bash
# Start from scratch with PostgreSQL
gocrete init blog-api \
  --module github.com/yourusername/blog-api \
  --db postgres \
  --docker

cd blog-api

# Start PostgreSQL
docker-compose up -d postgres

# Run your server
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/blog_api?sslmode=disable"
go run cmd/server/main.go
```

Or add to an existing project:

```bash
cd existing-project
gocrete add db --type postgres
```

### MongoDB

```bash
# Create with MongoDB
gocrete init user-service \
  --module github.com/yourusername/user-service \
  --db mongo \
  --docker

cd user-service

# Start MongoDB
docker-compose up -d mongo

# Run your server
export MONGO_URL="mongodb://localhost:27017"
export MONGO_DB="user_service"
go run cmd/server/main.go
```

## Adding OpenAPI

### Code Generation (Recommended)

```bash
# Create your API spec first
cat > api.yaml << 'EOF'
openapi: 3.0.0
info:
  title: My API
  version: 1.0.0
paths:
  /api/v1/users:
    get:
      summary: List users
      responses:
        '200':
          description: Success
EOF

# Generate project with OpenAPI
gocrete init my-api \
  --module github.com/yourusername/my-api \
  --openapi gen \
  --spec ./api.yaml

cd my-api

# Install oapi-codegen (for code generation)
go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest

# Generate API code
make api-gen

# Implement handlers in internal/api/handlers/
```

### Manual Implementation

```bash
gocrete init my-api \
  --module github.com/yourusername/my-api \
  --openapi manual

cd my-api

# Implement your handlers in internal/api/handlers/handlers.go
```

## Full-Featured Project

Create a production-ready API with everything:

```bash
gocrete init ecommerce-api \
  --module github.com/yourcompany/ecommerce-api \
  --router chi \
  --db postgres \
  --migrations goose \
  --openapi gen \
  --spec ./api.yaml \
  --docker

cd ecommerce-api

# Start all services
docker-compose up -d

# Run migrations (install goose first if needed)
goose -dir migrations postgres $DATABASE_URL up

# Start server
go run cmd/server/main.go
```

## Project Structure Explained

```
my-api/
├── cmd/
│   └── server/
│       └── main.go          # Entry point - starts the server
│
├── internal/                # Private application code
│   ├── config/              # Loads environment variables
│   │   └── config.go
│   ├── logger/              # Structured logging setup
│   │   └── logger.go
│   ├── http/                # HTTP server and routes
│   │   └── server.go
│   ├── errors/              # Error response helpers
│   │   └── errors.go
│   ├── db/                  # Database layer (if enabled)
│   │   ├── postgres/
│   │   └── mongo/
│   └── api/                 # API handlers (if OpenAPI enabled)
│       ├── generated/       # Generated code (don't edit)
│       └── handlers/        # Your implementations
│
├── migrations/              # Database migrations (if enabled)
├── api/                     # OpenAPI specs (if enabled)
├── docker-compose.yml       # Services (if Docker enabled)
├── Dockerfile               # App container (if Docker enabled)
├── go.mod                   # Go dependencies
└── README.md                # Project documentation
```

## Common Workflows

### Development

```bash
# Start dependencies
docker-compose up -d

# Run with hot reload (install air first)
air

# Or run normally
go run cmd/server/main.go
```

### Testing

```bash
# Run tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/db/postgres/...
```

### Building

```bash
# Build binary
go build -o bin/server cmd/server/main.go

# Run binary
./bin/server
```

### Deployment

```bash
# Build Docker image
docker build -t my-api:latest .

# Run with Docker Compose
docker-compose up -d

# Or deploy to your platform (Kubernetes, ECS, etc.)
```

## Customization

### Adding Your Own Logic

1. **Add handlers:**
   ```go
   // internal/api/handlers/users.go
   func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
       // Your logic here
   }
   ```

2. **Add routes:**
   ```go
   // internal/http/server.go
   r.Get("/api/v1/users", handlers.GetUsers)
   ```

3. **Add middleware:**
   ```go
   // internal/http/server.go
   r.Use(myCustomMiddleware)
   ```

### Environment Variables

Create a `.env` file:

```bash
ENVIRONMENT=development
PORT=8080
LOG_LEVEL=debug
DATABASE_URL=postgres://localhost:5432/mydb
```

Load it in your code or use a tool like `direnv`.

## Router Comparison

### When to use Chi (Default)
- You want standard library compatibility
- Building traditional REST APIs
- Need excellent middleware ecosystem

### When to use Gin
- You need maximum performance
- Large existing ecosystem
- Familiar with Express.js

### When to use Fiber
- You need the absolute fastest performance
- Building microservices
- Familiar with Express.js

## Next Steps

### Learn More
- Read the [full README](README.md)
- Check out [examples](examples/)
- Explore the [architecture](docs/ARCHITECTURE.md)

### Add More Features
- JWT authentication (coming soon)
- Caching with Redis (coming soon)
- Metrics with Prometheus (coming soon)

### Customize
- Modify generated code to fit your needs
- Add your business logic
- Extend the structure

### Deploy
- Build Docker images
- Deploy to Kubernetes, ECS, or your platform
- Set up CI/CD pipelines

## Troubleshooting

### "go: command not found"
Install Go 1.22 or higher from https://go.dev/dl/

### "port already in use"
Change the port:
```bash
export PORT=3000
go run cmd/server/main.go
```

### "database connection failed"
Check your DATABASE_URL and ensure PostgreSQL/MongoDB is running:
```bash
docker-compose up -d postgres  # or mongo
```

### "module not found"
Run:
```bash
go mod download
go mod tidy
```

## Tips

1. **Start Simple** - Begin with just `gocrete init`, add features later
2. **Use Docker** - Makes development easier with `--docker` flag
3. **Version Control** - Initialize git immediately after generation
4. **Read the Code** - Generated code is meant to be understood and modified
5. **Add Tests** - Start with test files early in development

## Getting Help

- Check [examples/](examples/) for common use cases
- Read [ARCHITECTURE.md](docs/ARCHITECTURE.md) for internals
- Open an issue on GitHub for bugs
- Start a discussion for questions

## What's Next?

Now that you have a running project:

1. **Implement your business logic** in `internal/`
2. **Add tests** alongside your code
3. **Configure CI/CD** for your repository
4. **Deploy** to your preferred platform
5. **Iterate** and grow your application

---

Happy coding! 🚀
