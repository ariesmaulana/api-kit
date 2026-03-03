# ARS Kit API

A REST API starter kit built with Go, Echo framework, and PostgreSQL.

## Tech Stack

- Go 1.23.4
- Echo v4 (HTTP framework)
- PostgreSQL (via pgx/v5)
- Zerolog (structured logging)
- Swagger (API documentation)

## Prerequisites

- Go 1.23.4 or higher
- PostgreSQL 12 or higher
- Make (for Makefile commands)

## Quick Start

```bash
# Install dependencies
make install-deps

# Copy and configure environment variables
cp .env.example .env
# Edit .env with your database credentials

# Set up database tables
psql -U <username> -d <database> -f database/sql/users.sql
psql -U <username> -d <database> -f database/sql/books.sql

# Optional: Add seed data
psql -U <username> -d <database> -f database/sql/seed_data.sql

# Run the application
make run
```

Server starts on `http://localhost:8181` (or configured PORT).

## Building

For comprehensive build instructions, deployment guides, and troubleshooting, see **[BUILD.md](./BUILD.md)**.

Quick build commands:

```bash
# Build for current OS
make build

# Build for Linux
make build-linux

# Build with tests
./build.sh --test

# See all build options
make help
./build.sh --help
```

## API Endpoints

- `GET /health` - Health check
- User endpoints (see Swagger docs)

## API Documentation

Access Swagger documentation at `/docs/api.html` when the server is running.

### Generating Swagger Documentation

To regenerate Swagger documentation after API changes:

```bash
# Install swag CLI (if not already installed)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs from annotations in code
swag init -g src/main.go -o docs

# Generate html file
npx @redocly/cli build-docs docs/swagger.yaml -o docs/api.html
```

The generated files will be placed in the `docs/` directory.

## Testing

```bash
# Run all tests
make test

# Run with verbose output
make test-verbose

# Run with coverage report
make test-coverage
```

The testing suite provides database-isolated test execution with automatic schema creation and cleanup. 

For detailed testing setup and documentation, see:
- **[BUILD.md](./BUILD.md#testing)** - Test setup and configuration
- **[testing/README.md](./testing/README.md)** - Testing framework documentation

## Project Structure

```
api-kit/
├── config/          # Configuration management
├── database/        # Database connection and SQL schemas
├── docs/            # API documentation (Swagger)
├── src/
│   ├── app/         # Application modules
│   │   └── user/    # User feature
│   └── main.go      # Application entry point
└── testing/         # Test suite and helpers
```

## License

Not specified
