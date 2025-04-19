Below is a multi‑level breakdown of the Ethereum/Gnosis Validator Pubkey Monitoring Tool, tailored for a Go‑based, test‑driven implementation. We start with a **high‑level blueprint**, then decompose into **iterative chunks**, then into **micro‑steps**, and finally produce **LLM prompts** for each micro‑step. Each prompt builds on prior work, mandates TDD, and ends with wiring the new code into the existing scaffold.

---

## 1. High‑Level Blueprint

1. **Architecture & Planning**  
   - Review specification: data model, integrations, UI, deployment citeturn0file0  
   - Choose libraries: HTTP router (e.g. chi), DB toolkit (sqlx or pgx), migration tool (goose)  
   - Define overall package layout (cmd/, internal/, pkg/, web/)

2. **Environment & Tooling Setup**  
   - Initialize Go module, linting (golangci‑lint), formatting (go fmt)  
   - CI/CD stub with GitHub Actions (unit tests)  
   - Dockerfile skeleton  

3. **Database Layer**  
   - Design Postgres schema: `validators` table (pubkey, network, status, client, timestamps) citeturn0file0  
   - Set up migrations (goose)  
   - DB connection management  

4. **Pubkey Management API**  
   - Input methods: REST endpoints for manual input, file upload, API import  
   - Validation: duplicate and format checks citeturn0file0  
   - Repository pattern  

5. **Status Monitoring**  
   - Service to poll local beacon nodes daily (and manually) for pubkey status  
   - Scheduler abstraction (cron)  

6. **Validator Client Detection**  
   - Query Lighthouse/Teku APIs based on config  
   - Map responses to our model citeturn0file0  

7. **Lido Pubkey Verification**  
   - Ethereum RPC integration  
   - Smart‑contract call for Lido validator registry on mainnet & Holesky citeturn0file0  

8. **Filtering & Frontend**  
   - Build JSON endpoints for UI filters (network, status)  
   - Scaffold React app (dark mode)  

9. **Audit Logging & Monitoring**  
   - Middleware capturing timestamp, action, source IP  
   - Expose Prometheus metrics  

10. **CI/CD, Containerization & Deployment**  
    - Complete GitHub Actions: lint, test, build, Docker image  
    - Nomad job spec  

---

## 2. First‑Round Chunks

| Chunk # | Name                          | Description                                                                        |
|:--------|:------------------------------|:-----------------------------------------------------------------------------------|
| 1       | Project Init & CI             | Go module, folder layout, linting, GitHub Actions stub, Dockerfile skeleton       |
| 2       | DB Schema & Migrations        | Goose setup, initial `validators` table migration                                 |
| 3       | DB Connection & Health        | `internal/db`, config loader, connect, ping, health endpoint                      |
| 4       | Model & Repository Interface  | `Validator` struct, CRUD interface, mock repo, unit tests                         |
| 5       | Pubkey Validation Logic       | Format checks, duplicate detection, unit tests                                     |
| 6       | Pubkey REST Endpoints         | Handlers for POST pubkey, file upload, API import, HTTP tests                      |
| 7       | Status Poller Service         | Cron‑based poller, manual trigger endpoint, stub beacon node client                |
| 8       | Client Detector Module        | Config loader, Lighthouse/Teku clients, detection logic, tests                     |
| 9       | Lido Verification Module      | Ethereum RPC client, smart contract ABI call, tests                                 |
| 10      | Filtering & Frontend API      | GET endpoints with query params, JSON responses                                     |
| 11      | Logging & Metrics             | HTTP middleware for audit logs, Prometheus instrumentation                          |
| 12      | Container & Deployment        | Dockerfile, Nomad job spec, end‑to‑end integration smoke test                       |

---

## 3. Second‑Round: Micro‑Steps

> **Note:** each step must include unit (or integration) tests and “wire in” the new code to the existing scaffold.

1. **Init Repository**  
   - `go mod init github.com/you/validator-monitor`  
   - Create folders: `cmd/validator-monitor`, `internal/`, `pkg/`, `web/`  
   - Add `golangci.yml`, pre‑commit hooks  

2. **CI/CD Stub**  
   - GitHub Actions workflow: on push → run `go test ./...`, `golangci-lint run`  
   - Add badges to README  

3. **Docker Skeleton**  
   - Basic Dockerfile: FROM golang:1.21, workdir, copy, build binary  
   - Ensure `docker build ./...` passes  

4. **Goose Migrations Setup**  
   - Add `goose` as dev dependency  
   - Create `migrations/0001_create_validators.sql` for table schema citeturn0file0  
   - Test `goose up/down` against local Postgres  

5. **DB Connection Package**  
   - `internal/db`: load config (env), open `*sql.DB` with ping  
   - Expose `func NewDB() (*sql.DB, error)`  
   - Unit‑test with `sqlmock`  

6. **Health Endpoint**  
   - In `cmd/validator-monitor/main.go`, mount `/healthz` that pings DB and returns 200  

7. **Validator Model & Repo Interface**  
   - Define `type Validator struct{ ... }` in `pkg/models`  
   - Interface `ValidatorRepo { Create, GetByPubkey, List(filters), UpdateStatus }`  
   - Generate `mockValidatorRepo` with `mockgen`  
   - Unit‑tests validating interface compliance  

8. **Implement Repo with sqlx**  
   - `internal/db/repo/validator_repo.go`: implement interface using `sqlx`  
   - Write SQL queries for CRUD  
   - Tests with `testcontainers` or `sqlmock`  

9. **Pubkey Validation Logic**  
   - `pkg/validator/validate.go`: pubkey hex format, length check  
   - Duplicate check via repo stub  
   - Unit‑tests covering invalid, duplicate, valid cases  

10. **REST Handler: Add Pubkey**  
    - Router setup (`chi.NewRouter()`) in main  
    - POST `/pubkeys` → bind JSON, call validation, repo.Create  
    - HTTP tests with `httptest`  

11. **File Upload Handler**  
    - POST `/upload` accepts CSV/JSON/TXT, parse lines, reuse validation logic  
    - Return summary of successes/failures  
    - Tests uploading sample files with `httptest`  

12. **Manual Refresh Endpoint**  
    - POST `/refresh` triggers status poller once  
    - Stub poller function to record invocation  
    - HTTP test  

13. **Status Poller Core**  
    - `pkg/poller/poller.go`: accept repo + beacon client, iterate validators, update statuses  
    - Unit‑test with fake beacon client  

14. **Cron Scheduler**  
    - Integrate `robfig/cron` in `cmd/validator-monitor` to call poller daily at configurable time  
    - Test scheduling config via unit test  

15. **Client Detector Logic**  
    - `pkg/detector/detector.go`: read client endpoints from config, detect by querying `/node/version` or similar  
    - Map to enum  
    - Unit‑tests with HTTP test server  

16. **Lido Verification Service**  
    - `pkg/lido/client.go`: Ethereum RPC via `go‑ethereum`, call registry contract method `isValidatorRegistered(pubkey)`  
    - Tests with simulated RPC server  

17. **Filtering API**  
    - GET `/pubkeys?network=eth&status=active` → repo.List(filters) → JSON  
    - HTTP tests  

18. **Audit Logging Middleware**  
    - `pkg/middleware/audit.go`: wrap handlers, log to stdout or file with timestamp, IP, action  
    - Unit‑test that middleware logs correct fields  

19. **Prometheus Metrics**  
    - Instrument HTTP endpoints with `promhttp`  
    - Expose `/metrics`  
    - Integration test that metrics endpoint serves counters  

20. **Nomad Job & Integration Smoke Test**  
    - Write `nomad.hcl` with task group, Docker driver, env vars  
    - A simple script that starts a Postgres container, Nomad dev agent, runs job, checks HTTP `/healthz`  
    - Add to CI as smoke test  

---

## 4. LLM Prompts for Test‑Driven Implementation

Below each prompt is a standalone request to a code‑generation LLM. Each must assume the code from previous prompts is already present and wired in.

---

### Prompt 1: Project Initialization

```text
You are setting up a new Go project for the Ethereum/Gnosis Validator Pubkey Monitoring Tool.  
- Initialize the repository at `github.com/you/validator-monitor`  
- Create folders: `cmd/validator-monitor`, `internal/`, `pkg/`, `web/`  
- Add `golangci.yml` for linting (enable all default linters)  
- Configure a pre‑commit hook to run `go fmt` and `golangci‑lint run`  
- Write unit tests to verify that the `internal/`, `pkg/`, and `cmd/validator-monitor` directories exist (using Go’s `os` package).  
Ensure everything passes when running `go test ./...` and `golangci‑lint run`.  
Wire this into a `Makefile` with targets `format`, `lint`, and `test`.  
```

---

### Prompt 2: CI/CD Stub

```text
Extend the project’s GitHub Actions:  
- Create `.github/workflows/ci.yml` that triggers on `push` and `pull_request`  
- Steps: checkout code, set up Go 1.21, run `make format`, `make lint`, and `make test`  
- Add status badges to the README for Go version, build status, and test coverage  
Write a test that asserts the `ci.yml` file exists and contains jobs named `lint` and `test`.  
```

---

### Prompt 3: Docker Skeleton

```text
Write a `Dockerfile` that:  
1. Uses `golang:1.21` as builder  
2. Copies the module files and downloads dependencies  
3. Builds the binary `validator-monitor`  
4. Uses `scratch` or `alpine` as the final image with the binary and a non‑root user  
5. Exposes port 8080  
Ensure `docker build -t validator-monitor .` succeeds.  
Add a `Makefile` target `docker-build` that wraps this.  
```

---

### Prompt 4: Goose Migrations Setup

```text
Integrate `goose` migrations:  
- Add `github.com/pressly/goose/v3` as a dev dependency  
- In `migrations/0001_create_validators.sql`, create table `validators` with columns:  
  - `id SERIAL PRIMARY KEY`  
  - `pubkey TEXT UNIQUE NOT NULL`  
  - `network TEXT NOT NULL`  
  - `status TEXT NOT NULL`  
  - `client TEXT`  
  - `created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()`  
  - `updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()`  
- Write a Go test in `internal/db/migrations_test.go` that runs `goose up` and `goose down` against a Dockerized Postgres using `testcontainers` and asserts no errors.  
```

---

### Prompt 5: DB Connection Package

```text
Implement `internal/db` package:  
- Load Postgres DSN from environment (e.g. `DATABASE_URL`)  
- Expose `func NewDB() (*sql.DB, error)` that opens with `sql.Open("postgres", dsn)` and `db.Ping()`  
- Use `sqlx` (`github.com/jmoiron/sqlx`) for extended features  
- Write unit tests with `sqlmock` to simulate a ping failure and a success.  
Wire `NewDB()` into `main.go` so that the application fails fast if the DB is unreachable.  
```

---

### Prompt 6: Health Endpoint

```text
In `cmd/validator-monitor/main.go`, using `chi` router:  
- Mount `GET /healthz` which calls `db.Ping()` and returns HTTP 200 with body `ok` if successful, otherwise 500  
- Unit-test the handler with `httptest` by mocking the DB to return error and success.  
Ensure the router is started on port from `PORT` env var (default 8080).  
```

---

### Prompt 7: Validator Model & Repo Interface

```text
Define `pkg/models/validator.go`:  
- `type Validator struct { ID int; Pubkey string; Network string; Status string; Client sql.NullString; CreatedAt time.Time; UpdatedAt time.Time }`  
- Create `ValidatorRepo` interface with methods:  
  - `Create(ctx, v *Validator) error`  
  - `GetByPubkey(ctx, pubkey string) (*Validator, error)`  
  - `List(ctx, filters map[string]interface{}) ([]Validator, error)`  
  - `UpdateStatus(ctx, pubkey, status string) error`  
- Generate a mock `mockValidatorRepo` using `mockgen`  
- Write tests that assert `mockValidatorRepo` implements `ValidatorRepo`.  
```

---

### Prompt 8: Implement Repo with sqlx

```text
Implement `internal/db/repo/validator_repo.go`:  
- Use `sqlx.DB` to write SQL for each interface method  
- Ensure `Create` returns error on duplicate pubkey  
- In `GetByPubkey`, return `sql.ErrNoRows` if not found  
- Write tests with `sqlmock` validating SQL queries and parameters for each method  
Wire this repo into a `NewValidatorService(repo ValidatorRepo)` in `pkg/service`.  
```

---

### Prompt 9: Pubkey Validation Logic

```text
In `pkg/validator/validate.go`:  
- Function `ValidatePubkeyFormat(pubkey string) error` ensures hex string of length 96 (48 bytes)  
- Function `CheckDuplicate(ctx, repo ValidatorRepo, pubkey string) error` returns error if `GetByPubkey` finds one  
- Write table‑driven tests for invalid length, non‑hex chars, and duplicates (use mock repo).  
Wire these into the POST handler in the next step.  
```

---

### Prompt 10: REST Handler for Pubkey Intake

```text
In your `web` layer:  
- Using `chi`, add POST `/pubkeys` that:  
  1. Decodes JSON `{ "pubkey": "0x..." }`  
  2. Calls `ValidatePubkeyFormat` and `CheckDuplicate`  
  3. Calls `repo.Create` to persist  
  4. Returns HTTP 201 with full `Validator` JSON  
- Write HTTP tests with `httptest` and a mock repo to cover success, format error, duplicate error, and DB error.  
Ensure the handler is wired into `main.go`.  
```

---

### Prompt 11: Status Poller Core

```text
Create `pkg/poller/poller.go` with:  
- `type BeaconClient interface { GetValidatorStatus(ctx, pubkey string) (string, error) }`  
- `func PollAndUpdate(ctx, repo ValidatorRepo, client BeaconClient) error` that:  
  1. Lists all validators (`repo.List`)  
  2. For each, calls `client.GetValidatorStatus`  
  3. Calls `repo.UpdateStatus` if changed  
- Write unit tests with a fake `BeaconClient` and mock repo verifying correct calls when statuses change and when errors occur.  
```

---

### Prompt 12: Cron Scheduler & Manual Refresh

```text
In `cmd/validator-monitor/main.go`:  
- Import `github.com/robfig/cron/v3`  
- Schedule `PollAndUpdate` to run daily at configurable hour (e.g. `CRON_SCHEDULE="0 3 * * *"`)  
- Expose POST `/refresh` that triggers `PollAndUpdate` immediately and returns 204  
- Test that cron entry is created with correct spec and that POST handler invokes poller (use a stub).  
```

---

Continue similarly for **Client Detection**, **Lido Verification**, **Filtering API**, **Audit Logging**, **Metrics**, and **Deployment**, each with its own **text‑tagged** prompt, ensuring:

- **TDD**: Every prompt ends with “Write tests …”  
- **Wiring**: “Ensure this is wired into `main.go` or the appropriate service.”  
- **No orphaned code**: Each module integrates into the previous scaffold.

This structure guarantees incremental, test‑driven progress, with no “big jumps” in complexity, and full integration at each step.