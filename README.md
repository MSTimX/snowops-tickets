# ticket-service

ticket-service is the SnowOps microservice responsible for managing tickets, assignments, and related workflows.

## Quick start

### Запуск ticket-service

```bash
go run ./cmd/ticket-service
```

Сервис слушает на порту `:8080` и предоставляет endpoint `/healthz`.

### Запуск всех сервисов

1. **Auth Service** (порт 7080):
   ```bash
   cd snowops-auth-service
   go run ./cmd/auth-service
   ```

2. **Operations Service** (порт 7081):
   ```bash
   cd snowops-operations-service
   go run ./cmd/operations-service
   ```

3. **Roles Service** (порт 7070):
   ```bash
   cd snowops-roles
   go run ./cmd/Snowops-roles
   ```

4. **Ticket Service** (порт 8080):
   ```bash
   go run ./cmd/ticket-service
   ```

### Порты сервисов

- **Auth Service**: `7080`
- **Operations Service**: `7081`
- **Roles Service**: `7070`
- **Ticket Service**: `8080`

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
JWT_ACCESS_SECRET=supersecret

# External services
AUTH_SERVICE_URL=http://localhost:7080
ROLES_SERVICE_URL=http://localhost:7070
OPERATIONS_SERVICE_URL=http://localhost:7081
AI_SERVICE_URL=

```

