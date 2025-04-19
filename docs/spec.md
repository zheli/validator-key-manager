# Ethereum/Gnosis Validator Pubkey Monitoring Tool Specification

## Overview
This tool is designed to monitor Ethereum and Gnosis validator pubkeys across mainnet and testnet environments. It tracks the status of each pubkey (unused, actively validating, slashed, exited), identifies the validator client in use (Lighthouse or Teku), and, for Lido-specific pubkeys, verifies if they are uploaded to the Lido smart contract. 

The tool supports manual input, file upload, and REST API import for pubkeys and integrates with local beacon nodes for status updates. It is containerized and designed for deployment on a HashiCorp Nomad cluster.

---

## Tech Stack
- **Backend:** Golang
- **Frontend:** React
- **Database:** Postgres
- **Deployment:** HashiCorp Nomad Cluster
- **CI/CD:** GitHub Actions
- **Monitoring:** Prometheus (potential integration)
- **Validator Clients Supported:** Lighthouse, Teku
- **Blockchain Networks Supported:** 
  - Ethereum (Mainnet, Holesky)
  - Gnosis (Mainnet, Testnet)

---

## Key Features
### 1. Pubkey Management
- **Input Methods:**
  - Manual input via frontend
  - File upload (CSV, Plain Text, JSON)
  - REST API import with API key authentication
- **Validation:**
  - Duplicate checks
  - Invalid key format checks

### 2. Status Monitoring
- **Pubkey Statuses:**
  - Unused
  - Actively validating
  - Slashed
  - Exited
- **Status Source:** Local hosted beacon nodes
- **Refresh Frequency:**
  - Daily automatic refresh
  - Manual refresh option

### 3. Validator Client Detection
- **Supported Clients:** Lighthouse, Teku
- **Method:**
  - Inspect local configuration
  - Query validator client APIs
- **API Connection:**
  - Manual configuration of API endpoints
  - Authentication using API tokens

### 4. Lido Pubkey Verification
- **Condition:** Only for Lido-related pubkeys
- **Check:** Smart contract query on:
  - Ethereum Mainnet
  - Holesky Testnet
- **Status Display:** Yes/No (No duplicates between blockchain and network)

### 5. Filtering and Display
- **Filter Options:**
  - By blockchain (Ethereum, Gnosis)
  - By network (Mainnet, Testnet)
- **Status Display Format:**
  - Text-based status in the backend
  - Responsive UI with dark mode in the frontend

### 6. Audit Logging
- **Data Logged:**
  - Timestamp
  - Action taken
  - Source IP

### 7. Documentation
- **README** for setup and run instructions
- **User Guides** for using the tool
- **API Documentation** (potentially using Swagger)

---

## Data Storage and Persistence
- **Database:** Postgres
- **Stored Data:**
  - Pubkeys
  - Statuses
  - Configurations

---

## Security and Access Control
- **API Authentication:** API key-based for REST imports
- **Role-based Permissions:** Not required for now, but system is designed to support it in the future
- **Data Encryption:** Not required

---

## Logging and Monitoring
- **Application Logging:** For debugging and audit purposes
- **Metrics and Health Checks:** 
  - Prometheus integration for monitoring tool performance and availability

---

## CI/CD
- **Pipeline:** GitHub Actions
- **Automated Testing:**
  - Unit tests
  - Integration tests
- **Continuous Deployment:** To HashiCorp Nomad cluster

---

## Deployment Requirements
- **Containerization:** Docker
- **Hosting Environment:** HashiCorp Nomad Cluster
- **Configuration Management:** Manual configuration for connecting to Lighthouse and Teku APIs

---

## Error Handling
- **File Upload:**
  - Validation for format (CSV, Plain Text, JSON)
  - Duplicate and invalid key checks
  - Detailed error messages for failed uploads
- **API Import:**
  - API key validation
  - Rate limiting not required
- **Status Querying:**
  - Graceful handling of network timeouts and connection failures
  - Retry mechanism for transient errors

---

## Future Considerations
- Role-based permissions and access control
- History tracking for validator status changes
- Notification and alerting system
- Support for additional validator clients (e.g., Prysm, Nimbus)

---

## Notes and Constraints
- No data export functionality required
- No backup or recovery mechanisms for the database
- No testing beyond automated CI/CD tests required

---

## End of Specification
