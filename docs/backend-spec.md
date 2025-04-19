# Backend Specification (Golang) – Ethereum/Gnosis Validator Pubkey Monitoring Tool

## Overview
The backend service is responsible for managing validator pubkeys, querying their status, identifying validator clients, verifying Lido upload status, and serving APIs for import and retrieval.

---

## Tech Stack
- **Language:** Golang
- **Database:** Postgres
- **Deployment:** Containerized (Docker), deployed to HashiCorp Nomad
- **CI/CD:** GitHub Actions (automated testing and deployment)
- **Monitoring:** Prometheus-compatible metrics and health endpoints

---

## Core Functionalities

### 1. Pubkey Import
- **Methods Supported:**
  - Manual entry (via frontend calling REST API)
  - File upload: support CSV, Plain Text, JSON
  - REST API import with API key auth
- **Validation:**
  - Deduplicate pubkeys
  - Verify key format (BLS public key format)

### 2. Pubkey Status Checking
- **Status Types:**
  - Unused
  - Actively validating
  - Slashed
  - Exited
- **Data Source:** Query local Ethereum/Gnosis beacon nodes (REST APIs)
- **Refresh Options:**
  - Scheduled daily sync
  - Manual trigger (via API)

### 3. Validator Client Discovery
- **Supported Clients:**
  - Lighthouse
  - Teku
- **Discovery Mode:**
  - Manual API endpoint config
  - Token-based API access
- **Query Type:** Check if the pubkey is registered/configured on the client

### 4. Lido Smart Contract Check (ETH only)
- **Smart Contract Query:**
  - Mainnet
  - Holesky
- **Output:** Yes/No if pubkey is in Lido contract registry
- **Assumption:** No duplicate pubkeys between blockchain + network

---

## REST API Endpoints
- `POST /import/pubkeys` – Upload pubkeys via file
- `POST /import/manual` – Submit pubkeys manually
- `POST /import/api` – Import via API with API key
- `GET /status/:pubkey` – Get current status of a pubkey
- `GET /validators` – List all validators (filterable by blockchain and network)
- `POST /refresh` – Trigger manual status refresh
- `GET /audit/logs` – View recent activity logs

---

## Authentication
- **Method:** API key authentication for REST imports
- **No user login or role-based access required** (but design for future extensibility)

---

## Data Model
- `Pubkey`: { pubkey, blockchain, network, imported_at }
- `ValidatorStatus`: { pubkey, status, last_checked }
- `ClientConfig`: { pubkey, client_type, detected_at }
- `LidoStatus`: { pubkey, uploaded: bool, network }
- `AuditLog`: { timestamp, action, source_ip }

---

## Logging and Monitoring
- **Audit Logging:**
  - Actions: import, query, refresh
  - Details: timestamp, action, source IP
- **Application Logs:** Debug/info/error logs
- **Monitoring:**
  - Health check endpoints
  - Prometheus metrics for status query durations, error rates, etc.

---

## Error Handling Strategy
- Upload: reject invalid formats and duplicates with detailed messages
- API: validate API key, return 401 on failure
- Network errors: retry with backoff; return 503 on beacon node failure
- Graceful fallback and logging on Lido contract failures

---

## Testing Plan
- **Unit Tests:** For all services and utility functions
- **Integration Tests:**
  - Beacon node responses
  - Validator client detection
- **CI/CD Pipeline:** GitHub Actions (run on push, test + build + deploy)
