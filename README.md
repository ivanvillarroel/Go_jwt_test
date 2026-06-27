# Go JWT Gin API

REST API built with Go, Gin, an in-memory SQLite database, dummy user data, and JWT authentication.

This Go implementation uses SQLite in memory to provide the same local and test-friendly behavior without an external database service.

## Run

```bash
go run ./cmd/api
```

## Login

```bash
curl -X POST http://localhost:8080/api/auth/token \
  -H 'Content-Type: application/json' \
  -d '{"username":"api-user","password":"password"}'
```

## List users

Requires a valid JWT with the `users:read` permission.

```bash
curl http://localhost:8080/api/users \
  -H "Authorization: Bearer <token>"
```

## Validate token

```bash
curl http://localhost:8080/api/auth/valid \
  -H "Authorization: Bearer <token>"
```

## Read token

```bash
curl http://localhost:8080/api/auth/read \
  -H "Authorization: Bearer <token>"
```

## Test

```bash
go test ./...
```
