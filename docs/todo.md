# TODO Checklist for Ethereum/Gnosis Validator Pubkey Monitoring Tool

Use this checklist to track progress through each development phase. Check items as you complete them.

## 1. Project Initialization

- [x] Initialize Go module: `go mod init github.com/zheli/validator-key-manager-backend`
- [x] Create folder structure:
  - [x] `cmd/validator-key-manager`
  - [x] `internal/`
  - [x] `pkg/`
- [x] Add `revive.toml` with default linter settings
- [x] Set up pre-commit hook to run `go fmt` and `revive` in`.pre-commit-config.yaml`
- [x] Write basic unit tests to verify directory existence
- [x] Create `Makefile` targets:
  - [x] `format`
  - [x] `lint`
  - [x] `test`

## 2. CI/CD Configuration

- [x] Create `.github/workflows/ci.yml`:
  - [x] Trigger on `push` and `pull_request`
  - [x] Steps: checkout, setup Go 1.24, run `make format`, `make lint`, `make test`
- [x] Add status badges to README for:
  - [x] Build status
  - [x] Test coverage
  - [x] Go version

## 3. Docker Skeleton

- [x] Write `Dockerfile`:
  - [x] Use `golang:1.24` as builder
  - [x] Download dependencies and build binary
  - [x] Use minimal base (scratch/alpine) with non-root user
  - [x] Expose port 8080
- [x] Add `Makefile` target `docker-build`
- [x] Verify `docker build -t validator-key-manager .` succeeds

## 4. Database Migrations

- [x] Add `migrate` as dev dependency
- [x] Create migration files for migration up and down
  - [x] `validators` table with columns:
    - `id SERIAL PRIMARY KEY`
    - `pubkey TEXT UNIQUE NOT NULL`
    - `blockchain TEXT NOT NULL`
    - `blockchian_network TEXT NOT NULL`
    - `status TEXT NOT NULL`
    - `client TEXT`
    - `created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()`
    - `updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()`
- [x] Write integration test using `testcontainers` to test `migrate` command with `postgres` container. Use the new syntax like this one:
  ```
  postgresContainer, err := postgres.Run(ctx,
      "postgres:16-alpine",
      postgres.WithDatabase(dbName),
      postgres.WithUsername(dbUser),
      postgres.WithPassword(dbPassword),
      testcontainers.WithWaitStrategy(
          wait.ForLog("database system is ready to accept connections").
              WithOccurrence(2).
              WithStartupTimeout(5*time.Second)),
  )
  ```

## 5. Database Connection Package

- [x] Implement `internal/db`:
  - [x] Ensure we can either load `DATABASE_URL` from environment variable, or a flag.
  - [x] Use `database/sql` for now to keep dependency small and simple.
  - [x] `func NewDB() (*sql.DB, error)` opens and pings DB
- [x] Write unit tests with `sqlmock` for ping success and failure
- [x] Integrate `NewDB()` into `main.go` for fail-fast startup

## 6. Health Check Endpoint

- [x] Set up `chi` router in `cmd/validator-monitor/main.go`
- [x] Mount `GET /healthz`:
  - [x] Ping database
  - [x] Respond `200 OK` with body `ok` or `500` on failure
- [x] Write HTTP tests using `httptest` (mock DB behavior)

## 7. Validator Model & Repository Interface

- [x] Define `pkg/models/validator.go` with `Validator` struct
- [x] Create `ValidatorRepo` interface:
  - [x] `Create(ctx, v *Validator) error`
  - [x] `GetByPubkey(ctx, pubkey string) (*Validator, error)`
  - [x] `List(ctx, filters map[string]interface{}) ([]Validator, error)`
  - [x] `UpdateStatus(ctx, pubkey, status string) error`
- [x] Generate `mockValidatorRepo` via `mockgen`
- [x] Write tests to confirm mock implements the interface

## 8. Repository Implementation with sql

- [x] Implement `validator_repo.go` in `internal/db/repo`
- [x] SQL for CRUD operations with parameter bindings
- [x] Handle duplicate errors in `Create`
- [x] Return `sql.ErrNoRows` for missing entries in `GetByPubkey`
- [x] Write unit tests with `sqlmock` verifying queries and args
- [x] Wire repo into `pkg/service.NewValidatorService`

## 9. Pubkey Validation Logic

- [x] Implement `ValidatePubkeyFormat(pubkey string) error` in `pkg/validator`
- [x] Write table-driven tests for: invalid length, non-hex chars, duplicates. Use consistent mocking approach (use generated mocks from mockgen)
- [x] Prepare for integration into REST handlers

## 10. REST Handlers: Validator Intake

- [ ] In `internal` folder, add POST `/validators` handler:
  - [ ] Accept JSON array of up to 1000 objects:
    - `pubkey` (string)
    - `blockchain` (e.g. ethereum, gnosis)
    - `network` (e.g. mainnet, holesky, chiado)
    - `is_testnet` (boolean)
    - `note` (optional string)
  - [ ] Validate Ethereum public key format
  - [ ] Validate array length <= 1000
  - [ ] Insert into DB in parallel using goroutines with context timeout (30s max)
  - [ ] Make sure the pubkey is unique in the DB and gracefully handle duplicate pubkey errors
  - [ ] Return summary JSON: success and failed inserts
- [ ] Write HTTP tests for:
  - [ ] Valid input under 1000
  - [ ] More than 1000 pubkeys
  - [ ] Timeout scenarios
  - [ ] Format and DB error handling
- [ ] Ensure handler is wired in `main.go`

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

- [ ] Add GET `/validators` with query params `network`, `status`, `client`
- [ ] Return filtered list from `repo.List`
- [ ] Write HTTP tests covering combinations of filters

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

## 19. Authentication and Authorization

- [ ] Add middleware to authenticate API requests
- [ ] Support two modes:
  - [ ] Google SSO for frontend users
  - [ ] API key-based for internal service calls
- [ ] Secure endpoints requiring identity or service access
- [ ] Write unit and integration tests for both auth flows

## 20. Deployment & Nomad

- [ ] Write `nomad.hcl` job spec:
  - [ ] Docker task using built image
  - [ ] Environment variables for DB, beacon node endpoints, etc.
- [ ] Create smoke test script:
  - [ ] Start local Postgres via testcontainers
  - [ ] Launch Nomad dev agent
  - [ ] Deploy job and verify `/healthz`
- [ ] Integrate smoke test into CI
