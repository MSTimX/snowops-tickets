# ticket-service

ticket-service is the SnowOps microservice responsible for managing tickets, assignments, and related workflows.

## Quick start

```bash
go run ./cmd/ticket-service
```

The service listens on `:8080` and exposes a `/health` endpoint.

## Environment

Создайте `.env` на основе примера ниже или выставите переменные окружения вручную:

```env
# HTTP
HTTP_PORT=8080
GIN_MODE=debug

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=snowops_tickets
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Almaty

# JWT
JWT_SECRET=dev-secret-change-me

# External services
AUTH_SERVICE_URL=http://localhost:8081
ROLES_SERVICE_URL=http://localhost:8082
OPERATIONS_SERVICE_URL=http://localhost:8083
AI_SERVICE_URL=

```

