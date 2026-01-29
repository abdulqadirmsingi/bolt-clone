# Mobility Platform

Real-time mobility and delivery platform (Bolt-style), optimized for **multi-stop rides** and **batched deliveries**. Push-based architecture, no polling.

## Repo layout

- **`mobility-backend/`** — Go backend: gRPC streaming, Fireball (RAMEN-style), H3 spatial indexing, Redis + PostgreSQL.
- **`mobility_app/`** — Flutter app: map-first UI, Kalman + dead reckoning, gRPC streams.
- **`ARCHITECTURE.md`** — System design, diagrams, and implementation notes.

## Quick start

### Backend (Go)

```bash
cd mobility-backend
go mod tidy
go run ./cmd/server
```

Uses `configs/local.yaml` by default. Set `CONFIG_PATH` for dev/prod. Requires PostgreSQL and Redis.

### Frontend (Flutter)

```bash
cd mobility_app
flutter pub get
flutter run
```

### Proto generation (optional)

From `mobility-backend`:

```bash
make proto   # requires protoc, protoc-gen-go, protoc-gen-go-grpc
```

## Principles

- **No polling** — gRPC bidirectional streaming; Fireball pushes only significant location updates.
- **H3 spatial indexing** — O(K² + M) nearby drivers; hexagons for uniform proximity.
- **Live data in Redis** — raw GPS never written to PostgreSQL.
- **Smooth UX** — Kalman filter + dead reckoning on the client for stable map markers.

See **ARCHITECTURE.md** for full design and explanations.
