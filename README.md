# Omoide Memo (OSS)

Self-hosted telemetry stack. Go ingestion API, Nuxt dashboard, generated TypeScript SDK — all in one pnpm workspace.

## Packages

- `packages/telemetry-api` — Go ingestion API (GORM + Postgres).
- `packages/telemetry-dashboard` — Nuxt 4 dashboard UI.
- `packages/telemetry-sdk` — Generated TypeScript SDK (workspace package).

## Run with Docker (prod)

Self-contained stack with caddy-docker-proxy as the single ingress on `:80`. Same-origin: dashboard at `/`, API at `/api/*` (path stripped by caddy before proxying).

```
cp .env.example .env
docker compose up --build
```

- Dashboard: `http://localhost`
- API: `http://localhost/api/v1/events`
- Postgres / API port / dashboard port are not exposed on the host — debug via `docker compose exec`.

To stop and wipe Postgres data: `docker compose down -v`.

## Run locally (dev)

Hot-reload dev server for the API + dashboard via `mise` + `mprocs`.

Prereqs: `mise`, `pnpm`, `docker` (for Postgres + the API's testcontainers).

```
mise run dev
```

This:
1. Generates the OpenAPI spec from the API source.
2. Generates the SDK into the dashboard.
3. Starts both processes under `mprocs` (API on `:9999`, dashboard on `:3000`).

## Seed demo data

In another terminal:

```
cd packages/telemetry-dashboard
pnpm seed
```

Generates 40 fake devices with ~30 days of events via the `ingestEvents` API. Requires `TELEMETRY_API_KEY` matching the API's config.
