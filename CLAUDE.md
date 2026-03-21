# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make build          # Build binary (go build ./cmd/apollo)
make test           # Run tests (requires DATABASE_URL pointing to postgres)
make test-setup     # Run migrations via golang-migrate
make lint           # Run golangci-lint

# Run a single test
go test ./internal/worker/... -run TestNotifications -v

# Full stack (api + workers + postgres + redis)
docker compose up --build -d

# Apply schema for tests (as CI does)
psql -f docs/schema.sql $DATABASE_URL
```

Tests require a live PostgreSQL database. `.env.test` configures the local test DB and Redis.

## Architecture

Apollo Backend is a **self-hosted push notification server** for Apollo (Reddit client). It polls Reddit inboxes on behalf of users and delivers iOS push notifications via APNs.

### Services

The binary is a single Go binary with multiple cobra subcommands, run as separate processes:

| Command | Role |
|---------|------|
| `apollo api` | HTTP API (port 4000) — device registration, account linking, watcher CRUD |
| `apollo scheduler` | Runs every 5s — enqueues accounts/watchers ready for checking |
| `apollo worker --queue notifications` | Polls Reddit inboxes, sends APNs pushes (16 consumers) |
| `apollo worker --queue stuck-notifications` | Recovers jobs stuck in queue (8 consumers) |
| `apollo worker --queue users` | User watcher jobs |
| `apollo worker --queue subreddits` | Subreddit keyword/flair/author/domain watchers |
| `apollo worker --queue trending` | Trending post watchers |
| `apollo worker --queue live-activities` | iOS 16+ Dynamic Island / Lock Screen updates |

### Data Flow

1. **Registration**: iOS app POSTs device APNS token + Reddit OAuth tokens to `/api/v1/device` and `/api/v1/accounts`
2. **Scheduling**: Scheduler queries accounts where `next_notification_check_at <= now`, acquires Redis locks (Lua script), enqueues batches of 250 to rmq queues
3. **Processing**: Workers dequeue jobs, call Reddit API for new inbox messages/posts, send APNs pushes
4. **Lock Release**: After processing, lock expires (5 min `NotificationCheckTimeout`) preventing duplicate notifications

### Key Packages

- `internal/cmd/` — Cobra command implementations; each service bootstraps its own dependencies via `cmdutil`
- `internal/api/` — HTTP handlers using gorilla/mux; thin layer over repositories
- `internal/domain/` — Domain structs and repository interfaces (Account, Device, Watcher, etc.)
- `internal/repository/` — PostgreSQL implementations of domain interfaces (pgx/v5)
- `internal/worker/` — rmq consumer implementations; `notifications.go` is the core
- `internal/reddit/` — Reddit OAuth client with rate limit awareness and exponential backoff on 429s
- `internal/cmdutil/` — Factory functions for shared infrastructure (DB pool, Redis, StatsD, logger)

### Distributed Locking

Two separate Redis instances:
- **redis-queues**: rmq job queue
- **redis-locks**: Lua script-based distributed lock preventing duplicate work within 5-minute windows

### Database

PostgreSQL with golang-migrate migrations in `migrations/`. Schema is also exported to `docs/schema.sql` for CI test setup. Key tables: `accounts`, `devices`, `devices_accounts` (junction with per-account notification settings), `watchers`, `live_activities`.

## Configuration

Copy `.env.example` to `.env`. Required variables:

- **Apple APNs**: `APPLE_BUNDLE_ID`, `APPLE_KEY_ID`, `APPLE_TEAM_ID`, `APPLE_KEY_PATH` (path to `.p8` file)
- **Reddit OAuth**: `REDDIT_CLIENT_ID`, `REDDIT_CLIENT_SECRET`, `REDDIT_USER_AGENT`
- **Infrastructure**: `DATABASE_CONNECTION_POOL_URL`, `REDIS_QUEUE_URL`, `REDIS_LOCKS_URL`
- **Optional auth**: `APOLLO_SECRET_TOKEN` — enables `X-Apollo-Token` header validation on all API routes
