# TODO Checklist for Ethereum/Gnosis Validator Pubkey Monitoring Tool

Use this checklist to track progress through each development phase. Check items as you complete them.

## 1. Project Initialization

- [ ] Initialize Go module: `go mod init github.com/zheli/validator-key-manager-backend`
- [ ] Create folder structure:
  - [ ] `cmd/validator-key-manager`
  - [ ] `internal/`
  - [ ] `pkg/`
  - [ ] `web/`
- [ ] Add `golangci.yml` with default linter settings
- [ ] Set up pre-commit hook to run `go fmt` and `golangci-lint run`
- [ ] Write basic unit tests to verify directory existence
- [ ] Create `Makefile` targets:
  - [ ] `format`
  - [ ] `lint`
  - [ ] `test`

## 2. CI/CD Configuration

- [ ] Create `.github/workflows/ci.yml`:
  - [ ] Trigger on `push` and `pull_request`
  - [ ] Steps: checkout, setup Go 1.21, run `make format`, `make lint`, `make test`
- [ ] Add status badges to README for:
  - [ ] Build status
  - [ ] Test coverage
  - [ ] Go version
- [ ] Write a test to assert existence and contents of `ci.yml`

## 3. Docker Skeleton

- [ ] Write `Dockerfile`:
  - [ ] Use `golang:1.21` as builder
  - [ ] Download dependencies and build binary
  - [ ] Use minimal base (scratch/alpine) with non-root user
  - [ ] Expose port 8080
- [ ] Add `Makefile` target `docker-build`
- [ ] Verify `docker build -t validator-monitor .` succeeds

## 4. Database Migrations

- [ ] Add `goose` as dev dependency
- [ ] Create migration `migrations/0001_create_validators.sql`:
  - [ ] `validators` table with columns:
    - `id SERIAL PRIMARY KEY`
    - `pubkey TEXT UNIQUE NOT NULL`
    - `network TEXT NOT NULL`
    - `status TEXT NOT NULL`
    - `client TEXT`
    - `created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()`
    - `updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()`
- [ ] Write integration test using `testcontainers` to run `goose up/down`

## 5. Database Connection Package

- [ ] Implement `internal/db`:
  - [ ] Load `DATABASE_URL` from environment
  - [ ] `func NewDB() (*sqlx.DB, error)` opens and pings DB
- [ ] Write unit tests with `sqlmock` for ping success and failure
- [ ] Integrate `NewDB()` into `main.go` for fail-fast startup

## 6. Health Check Endpoint

- [ ] Set up `chi` router in `cmd/validator-monitor/main.go`
- [ ] Mount `GET /healthz`:
  - [ ] Ping database
  - [ ] Respond `200 OK` with body `ok` or `500` on failure
- [ ] Write HTTP tests using `httptest` (mock DB behavior)

## 7. Validator Model & Repository Interface

- [ ] Define `pkg/models/validator.go` with `Validator` struct
- [ ] Create `ValidatorRepo` interface:
  - [ ] `Create(ctx, v *Validator) error`
  - [ ] `GetByPubkey(ctx, pubkey string) (*Validator, error)`
  - [ ] `List(ctx, filters map[string]interface{}) ([]Validator, error)`
  - [ ] `UpdateStatus(ctx, pubkey, status string) error`
- [ ] Generate `mockValidatorRepo` via `mockgen`
- [ ] Write tests to confirm mock implements the interface

## 8. Repository Implementation with sqlx

- [ ] Implement `validator_repo.go` in `internal/db/repo`
- [ ] SQL for CRUD operations with parameter bindings
- [ ] Handle duplicate errors in `Create`
- [ ] Return `sql.ErrNoRows` for missing entries in `GetByPubkey`
- [ ] Write unit tests with `sqlmock` verifying queries and args
- [ ] Wire repo into `pkg/service.NewValidatorService`

## 9. Pubkey Validation Logic

- [ ] Implement `ValidatePubkeyFormat(pubkey string) error` in `pkg/validator`
- [ ] Implement `CheckDuplicate(ctx, repo, pubkey) error`
- [ ] Write table-driven tests for: invalid length, non-hex chars, duplicates
- [ ] Prepare for integration into REST handlers

## 10. REST Handlers: Pubkey Intake

- [ ] In `web` layer, add POST `/pubkeys` handler:
  - [ ] Parse JSON input `{ "pubkey": "0x..." }`
  - [ ] Validate format and duplicates
  - [ ] Call `repo.Create` and return created record
- [ ] Write HTTP tests for:
  - [ ] Success case
  - [ ] Invalid format
  - [ ] Duplicate error
  - [ ] DB errors
- [ ] Ensure wiring in `main.go`

## 11. File Upload Endpoint

- [ ] Add POST `/upload`:
  - [ ] Accept CSV, JSON, or TXT
  - [ ] Parse and iterate pubkeys
  - [ ] Reuse validation and repo logic
  - [ ] Return summary of successes/failures
- [ ] Write HTTP tests with sample files via `httptest`

## 12. Status Poller Core

- [ ] Define `BeaconClient` interface in `pkg/poller`
- [ ] Implement `PollAndUpdate(ctx, repo, client)`:
  - [ ] List validators
  - [ ] Fetch status from beacon client
  - [ ] Update status in DB if changed
- [ ] Write unit tests with fake `BeaconClient` and mock repo

## 13. Scheduler & Manual Refresh

- [ ] Integrate `robfig/cron` in `main.go`
  - [ ] Schedule poller daily via env `CRON_SCHEDULE`
- [ ] Add POST `/refresh` to trigger poller immediately
- [ ] Write tests ensuring cron entry and HTTP handler invoke the poller stub

## 14. Validator Client Detection

- [ ] Implement `pkg/detector`:
  - [ ] Load API endpoints from config
  - [ ] Query `/node/version` or equivalent
  - [ ] Detect Lighthouse vs Teku
- [ ] Write tests using HTTP test server stubs
- [ ] Save detected client in `Validator.Client`

## 15. Lido Pubkey Verification

- [ ] Use `go-ethereum` RPC client in `pkg/lido`
- [ ] Implement `IsRegistered(pubkey string) (bool, error)` calling Lido registry contract
- [ ] Test with simulated RPC server or `ethclient` mock
- [ ] Wire result into status view

## 16. Filtering & Frontend API

- [ ] Add GET `/pubkeys` with query params `network`, `status`, `client`
- [ ] Return filtered list from `repo.List`
- [ ] Write HTTP tests covering combinations of filters
- [ ] Scaffold React app endpoint and dark mode toggle

## 17. Audit Logging Middleware

- [ ] Create `pkg/middleware/audit.go`:
  - [ ] Capture timestamp, action (HTTP method + path), source IP
  - [ ] Log to stdout or file
- [ ] Write unit tests verifying log output contains required fields
- [ ] Integrate middleware into router chain

## 18. Prometheus Metrics

- [ ] Instrument HTTP handlers with `promhttp` middleware
- [ ] Expose `/metrics` endpoint
- [ ] Write integration test to fetch `/metrics` and verify counters exist

## 19. Deployment & Nomad

- [ ] Write `nomad.hcl` job spec:
  - [ ] Docker task using built image
  - [ ] Environment variables for DB, beacon node endpoints, etc.
- [ ] Create smoke test script:
  - [ ] Start local Postgres via testcontainers
  - [ ] Launch Nomad dev agent
  - [ ] Deploy job and verify `/healthz`
- [ ] Integrate smoke test into CI

---

Once complete, all features will be implemented incrementally, with full test coverage and safe integration at each stage.

