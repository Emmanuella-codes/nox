# Workers

Run the recurring backend workers as separate processes from `apps/api`.

Commands:

- `go run ./workers/notifications`
- `go run ./workers/stories`
- `go run ./workers/media`

Shared runtime behavior:

- load config
- optionally run migrations
- connect Postgres
- connect Redis
- initialize repositories
- stop cleanly on `SIGINT` or `SIGTERM`

Current worker-specific config:

- notifications:
  - `PUSH_WORKER_BATCH_SIZE`
  - `PUSH_WORKER_POLL_INTERVAL`
- stories:
  - `STORY_CLEANUP_BATCH_SIZE`
  - `STORY_CLEANUP_INTERVAL`
  - `STORY_EXPIRY_RETENTION`
- media:
  - `MEDIA_CLEANUP_BATCH_SIZE`
  - `MEDIA_CLEANUP_INTERVAL`
  - `MEDIA_PENDING_RETENTION`
  - `MEDIA_FAILED_RETENTION`
