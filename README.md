# URL Shortener

A small URL-shortening HTTP service written in Go. It generates a random
10-character code for a submitted URL, stores the mapping in PostgreSQL, and
redirects requests made to that code.

## Features

- Create a short code for a valid URL.
- Redirect from `/<short_code>` to the stored original URL.
- Health check endpoint.
- PostgreSQL persistence with schema migrations and `sqlc`-generated queries.
- Cryptographically secure, alphanumeric short-code generation.

## Technology

- Go 1.25
- [Gin](https://github.com/gin-gonic/gin) HTTP framework
- PostgreSQL 17 (used by the included Docker Compose configuration)
- [`pgx`](https://github.com/jackc/pgx) PostgreSQL driver
- [`sqlc`](https://sqlc.dev/) generated database access code
- [`golang-migrate`](https://github.com/golang-migrate/migrate) migration files and CLI

## Requirements

For local development, install:

- Go 1.25 or later
- Docker Engine with the Docker Compose plugin, or a reachable PostgreSQL
  database
- The `migrate` CLI to apply database migrations

To regenerate the database package from the SQL sources, install `sqlc`.

## Run locally

### 1. Start PostgreSQL

The included Compose file starts PostgreSQL only (it does not start the Go
application):

```sh
docker compose up -d postgres
```

The database is published on `localhost:5433`, with database `url_shortener`
and the development credentials `postgres` / `postgres`, as configured in
[`docker-compose.yml`](./docker-compose.yml). These credentials are for local
development only.

### 2. Configure the application

Copy the example environment file and set the values:

```sh
cp app.env.example app.env
```

Set `DB_SOURCE` to a PostgreSQL connection URL and `PORT` to the address the
HTTP server should listen on. For the included Compose database, use:

```dotenv
DB_SOURCE=postgres://postgres:postgres@localhost:5433/url_shortener?sslmode=disable
PORT=:8000
```

The app loads `app.env` at startup when that file is present. It can also read
these values from the process environment. `DB_SOURCE` must point to a
reachable database. `PORT` is passed to Gin as the listen address; `:8000`
listens on all interfaces on port 8000.

The example file currently also contains `DB_DRVIER` (misspelled). The
application does not read a driver setting; it uses `pgx` directly. Only
`DB_SOURCE` and `PORT` are used by the current application.

### 3. Apply the database migration

With `migrate` installed and the connection URL available in your shell, apply
the schema:

```sh
set -a
. ./app.env
set +a
migrate -path migrations/ -database "$DB_SOURCE" -verbose up
```

The initial migration creates the `urls` table and its unique short-code
constraint. The matching `down` migration drops that table. **Running the down
migration deletes all stored URL mappings.**

The Makefile also defines `make migrateup` and `make migratedown`; those
targets use a connection string configured directly in the Makefile rather
than reading `DB_SOURCE`. Check that connection string before using either
target.

### 4. Start the service

From the repository root:

```sh
go run .
```

Or build and run a binary:

```sh
go build -o url-shortener .
./url-shortener
```

The server listens on the configured `PORT` (for example, `:8000`). The
database must be running and the schema migration must have completed before
creating or resolving URLs.

## HTTP API

### Health check

```http
GET /health
```

Returns HTTP `200`:

```json
{
  "message": "Healthy"
}
```

This endpoint checks that the HTTP handler is serving requests; it does not
perform a database health check.

### Create a short URL

```http
POST /url
Content-Type: application/json
```

Request body:

```json
{
  "original_url": "https://example.com/some/long/path"
}
```

`original_url` is required and validated as a URL by Gin. A successful request
returns HTTP `200` and the stored URL record:

```json
{
  "id": 1,
  "short_code": "a1B2c3D4eF",
  "original_url": "https://example.com/some/long/path",
  "created_at": "2025-01-01T12:00:00Z",
  "expires_at": null
}
```

The code is a random 10-character string from the digits, lowercase letters,
and uppercase letters. The database enforces uniqueness. A database error,
including a code collision, is returned as HTTP `500`; code generation does
not currently retry collisions.

Example:

```sh
curl -i -X POST http://localhost:8000/url \
  -H 'Content-Type: application/json' \
  -d '{"original_url":"https://example.com/some/long/path"}'
```

### Redirect to the original URL

```http
GET /<short_code>
```

For example:

```sh
curl -i http://localhost:8000/a1B2c3D4eF
```

When the code is found, the server responds with HTTP `302 Found` and a
`Location` header containing the original URL.

**Current limitation:** the redirect handler does not stop after a database
lookup error. Consequently, behavior for an unknown code is not a reliable
not-found response. Callers should not depend on a specific status or response
body for missing codes until that handler is corrected.

## Configuration

| Variable | Required | Description |
| --- | --- | --- |
| `DB_SOURCE` | Yes for database operations | PostgreSQL connection URL passed to `pgx.Connect`. |
| `PORT` | Recommended | HTTP listen address passed to Gin, such as `:8000`. |

No `DB_DRIVER` setting is used. The current template's `DB_DRVIER` key is a
typo and has no effect.

## Database and generated code

The schema is defined in
[`migrations/000001_init.up.sql`](./migrations/000001_init.up.sql). It creates
`urls` with:

- An auto-incrementing `id`.
- A unique `short_code` of up to 10 characters.
- The required `original_url`.
- A `created_at` timestamp defaulting to the current time.
- A nullable `expires_at` timestamp.

SQL statements are maintained in
[`internal/database/queries/urls.sql`](./internal/database/queries/urls.sql).
The generated Go code is in [`internal/database/sqlc/`](./internal/database/sqlc/),
configured by [`sqlc.yaml`](./sqlc.yaml). Regenerate it after changing the
queries or schema:

```sh
sqlc generate
```

Do not edit the generated files directly. A delete query exists in the
database layer, but the HTTP API does not currently expose a delete endpoint.
Although the schema has an `expires_at` column, URL creation and redirects do
not set or enforce expiration.

## Docker

Build the application image:

```sh
docker build -t url-shortener .
```

The image listens on port 8000. To run it, provide the environment and publish
the port:

```sh
docker run --rm \
  --env-file app.env \
  -p 8000:8000 \
  url-shortener
```

Make sure `DB_SOURCE` is reachable **from inside the container**. In
particular, `localhost` inside the application container refers to that
container, not the host or the PostgreSQL container. If using Compose, the
application would need to join the database service's network and use the
database service name as its hostname; the current Compose file does not define
an application service.

## Development commands

Run all Go tests:

```sh
go test ./...
```

Build the application:

```sh
go build -o url-shortener .
```

Regenerate SQL access code:

```sh
sqlc generate
```

The current unit tests cover short-code length, allowed characters, and
duplicates across a small number of generated codes.

## Project layout

```text
.
├── main.go                          # Loads configuration and starts DB/server
├── internal/
│   ├── handler/handler.go           # Gin routes and HTTP handlers
│   ├── service/url_creation.go      # Random short-code generation
│   ├── service/url_creation_test.go # Short-code unit tests
│   ├── util/utils.go                # PostgreSQL connection and query setup
│   └── database/
│       ├── queries/urls.sql         # Source SQL queries for sqlc
│       └── sqlc/                    # Generated models and query methods
├── migrations/                      # PostgreSQL schema migrations
├── Dockerfile                       # Multi-stage application image build
├── docker-compose.yml               # Local PostgreSQL service
├── Makefile                         # Migration and sqlc convenience targets
└── Jenkinsfile                      # Test, build, and Docker publish pipeline
```

## CI pipeline

The Jenkins pipeline checks out the repository, runs `go test ./...`, builds
the Go binary, builds the Docker image, and pushes the `latest` image to the
configured Docker registry. Registry credentials are expected to be supplied
through Jenkins credentials; they are not part of the application's local
configuration.

## Current scope and limitations

- There is no authentication, authorization, rate limiting, or custom
  application-level request logging.
- There is no URL listing, update, or delete HTTP endpoint.
- Expiration is represented in the schema but is not currently assigned or
  enforced.
- Duplicate short-code insertion is not retried.
- Unknown-code lookup handling in the redirect route needs correction, as
  described above.
