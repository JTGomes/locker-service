# Bloq.it Challenge - Locker Service API
Software Engineering Challenge, by Bloqit.

##  Installation & Setup

### Requirements

- Go 1.25+
- PostgreSQL 16+
- Docker & Docker Compose (optional)

### Environment Variables

Create a local `.env` from the example and fill the variables:
    ```bash
    cp .env-example .env
    ```


### Option 1: Local Execution
You need a reachable PostgreSQL instance. The quickest option is to run the
database in Docker and the app on your machine:

```bash
# 1. Start Postgres (and keep it running in the background)
docker compose up -d postgres

# 2. Make sure .env points at it (DB_HOST=localhost, DB_PORT=5432)

# 3. Download dependencies
make install

# 4. Run the API — migrations are applied automatically on startup
make run
```

The API will be available at `http://localhost:8080`

### Option 2: Docker Compose

```bash
make docker-run
```

This starts:
- API at `http://localhost:8080`
- PostgreSQL at `localhost:5432`

To stop:
```bash
make docker-down
```

## API Endpoints

### Health Check
- `GET /health` - Check API status

**Status codes**

| Code | Meaning                                                              |
| ---- | ------------------------------------------------------------------- |
| 200  | OK                                                                  |
| 201  | Created                                                            |
| 400  | Validation error / malformed or unknown JSON fields                |
| 404  | Resource not found                                                 |
| 409  | Conflict (illegal status transition or concurrent update)         |
| 500  | Internal server error                                              |

**Pagination** query parameters (on all list endpoints): `limit` (1–100, default `20`) and
`offset` (≥ 0, default `0`).

---

### Endpoints

| Method | Path                               | Description                       |
| ------ | ---------------------------------- | --------------------------------- |
| GET    | `/health`                          | Health check                      |
| POST   | `/api/v1/bloq`                     | Create a bloq                     |
| GET    | `/api/v1/bloq`                     | List bloqs                        |
| GET    | `/api/v1/bloq/{id}`                | Get a bloq by ID                  |
| DELETE | `/api/v1/bloq/{id}`                | Delete a bloq                     |
| POST   | `/api/v1/locker`                   | Create a locker                   |
| GET    | `/api/v1/locker`                   | List lockers                      |
| GET    | `/api/v1/locker/{id}`              | Get a locker by ID                |
| DELETE | `/api/v1/locker/{id}`              | Delete a locker                   |
| POST   | `/api/v1/rent`                     | Create a rent                     |
| GET    | `/api/v1/rent/{id}`                | Get a rent by ID                  |
| POST   | `/api/v1/rent/{id}/allocate`       | Allocate locker to rent           |
| POST   | `/api/v1/rent/{id}/dropoff`        | Register dropoff                  |
| POST   | `/api/v1/rent/{id}/pickup`         | Register pickup                   |

---

#### Create a bloq

```bash
curl -X POST http://localhost:8080/api/v1/bloq \
  -H 'Content-Type: application/json' \
  -d '{"title":"Downtown Bloq","address":"Main Street, 123"}'
```

`201 Created`

```json
{
  "id": "b0c1a2d3-4e5f-6789-abcd-0123456789ab",
  "title": "Downtown Bloq",
  "address": "Main Street, 123"
}
```

#### List bloqs

```bash
curl 'http://localhost:8080/api/v1/bloq?limit=20&offset=0'
```

`200 OK`

```json
{
  "data": [
    {
      "id": "b0c1a2d3-4e5f-6789-abcd-0123456789ab",
      "title": "Downtown Bloq",
      "address": "Main Street, 123"
    }
  ],
  "total": 1
}
```

#### Get a bloq

```bash
curl http://localhost:8080/api/v1/bloq/b0c1a2d3-4e5f-6789-abcd-0123456789ab
```

`200 OK` or `404 Not Found`.

```json
{
  "id": "b0c1a2d3-4e5f-6789-abcd-0123456789ab",
  "title": "Downtown Bloq",
  "address": "Main Street, 123"
}
```

#### Create a locker

```bash
curl -X POST http://localhost:8080/api/v1/locker \
  -H 'Content-Type: application/json' \
  -d '{"bloq_id":"b0c1a2d3-4e5f-6789-abcd-0123456789ab","status":"closed"}'
```

`201 Created`

```json
{
  "id": "c1d2e3f4-5678-90ab-cdef-1234567890ab",
  "bloq_id": "b0c1a2d3-4e5f-6789-abcd-0123456789ab",
  "status": "closed",
  "is_occupied": false,
  "created_at": "2026-06-25T12:00:00Z",
  "updated_at": "2026-06-25T12:00:00Z"
}
```

#### List lockers

```bash
curl 'http://localhost:8080/api/v1/locker?limit=20&offset=0'
```

`200 OK`

```json
{
  "data": [
    {
      "id": "c1d2e3f4-5678-90ab-cdef-1234567890ab",
      "bloq_id": "b0c1a2d3-4e5f-6789-abcd-0123456789ab",
      "status": "closed",
      "is_occupied": false,
      "created_at": "2026-06-25T12:00:00Z",
      "updated_at": "2026-06-25T12:00:00Z"
    }
  ],
  "total": 1
}
```

#### Create a rent

```bash
curl -X POST http://localhost:8080/api/v1/rent \
  -H 'Content-Type: application/json' \
  -d '{"weight":1800,"size":"M"}'
```

`201 Created` — rent starts with `created` status

```json
{
  "id": "d2e3f4a5-6789-0abc-def1-234567890abc",
  "locker_id": null,
  "weight": 1800,
  "size": "M",
  "status": "created",
  "created_at": "2026-06-25T12:00:00Z",
  "updated_at": "2026-06-25T12:00:00Z",
  "dropped_off_at": null,
  "picked_up_at": null
}
```

#### Allocate a locker

```bash
curl -X POST http://localhost:8080/api/v1/rent/d2e3f4a5-6789-0abc-def1-234567890abc/allocate \
  -H 'Content-Type: application/json' \
  -d '{"bloq_id":"b0c1a2d3-4e5f-6789-abcd-0123456789ab"}'
```

`200 OK` on success, `409 Conflict` on allocation race or invalid transition.

```json
{
  "id": "d2e3f4a5-6789-0abc-def1-234567890abc",
  "locker_id": "c1d2e3f4-5678-90ab-cdef-1234567890ab",
  "weight": 1800,
  "size": "M",
  "status": "waiting_dropoff",
  "created_at": "2026-06-25T12:00:00Z",
  "updated_at": "2026-06-25T12:05:00Z"
}
```


### Lockers
- `POST /api/v1/locker` - Create a new locker
- `GET /api/v1/locker` - List lockers (with filters and pagination)
- `GET /api/v1/locker/:id` - Get specific locker
- `DELETE /api/v1/locker/:id` - Delete locker



**Available filters for listing:**
- `bloq_id` - Filter by bloq
- `status` - Filter by status (open/closed)
- `is_occupied` - Filter by occupancy (true/false)
- `limit` - Items per page (default: 10)
- `offset` - Pagination offset (default: 0)

### Rents
- `POST /api/v1/rent` - Create a new rent
- `GET /api/v1/rent/:id` - Get specific rent
- `POST /api/v1/rent/:id/allocate` - Allocate locker to rent
- `POST /api/v1/rent/:id/dropoff` - Register dropoff
- `POST /api/v1/rent/:id/pickup` - Register pickup

**Request example (Create Rent):**
```bash
curl -X POST http://localhost:8080/api/v1/rent \
  -H "Content-Type: application/json" \
  -d '{
    "weight": 2.5,
    "size": "M"
  }'
```

**Response example (201 Created):**
```json
{
  "id": "d2e3f4a5-6789-0abc-def1-234567890abc",
  "locker_id": null,
  "weight": 2.5,
  "size": "M",
  "status": "created",
  "created_at": "2026-06-25T12:00:00Z",
  "updated_at": "2026-06-25T12:00:00Z",
  "dropped_off_at": null,
  "picked_up_at": null
}
```

**Response example (200 OK after allocation):**
```json
{
  "id": "d2e3f4a5-6789-0abc-def1-234567890abc",
  "locker_id": "c1d2e3f4-5678-90ab-cdef-1234567890ab",
  "weight": 2.5,
  "size": "M",
  "status": "waiting_dropoff",
  "created_at": "2026-06-25T12:00:00Z",
  "updated_at": "2026-06-25T12:05:00Z",
  "dropped_off_at": null,
  "picked_up_at": null
}
```

## Development

### Testing

```bash
# Run all tests
make test

# Run tests with race condition detection
make test-race
```

### Migrations

```bash
# Apply migrations
make migrate-up

# Undo last migration
make migrate-down
```


## Available Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8080 | API port |
| DB_HOST | localhost | PostgreSQL host |
| DB_PORT | 5432 | PostgreSQL port |
| POSTGRES_DB | bloq_db | Database name |
| POSTGRES_USER | postgres | PostgreSQL user |
| POSTGRES_PASSWORD | - | PostgreSQL password |
| DB_MAX_CONNS | 20 | Maximum number of open database connections in the pool |
| DB_MIN_CONNS | 5 | Minimum number of idle connections to keep available |
| DB_MAX_CONN_LIFETIME | 3600 | Maximum lifetime for each connection in seconds before it is recycled |
| DB_MAX_CONN_IDLE_TIME | 900 | Maximum idle time in seconds before a connection is closed |
| DB_HEALTHCHECK_PERIOD | 60 | Interval in seconds between database health checks |

