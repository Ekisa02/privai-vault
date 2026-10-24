# PrivAI Vault

PrivAI Vault is a modular Go service demonstrating privacy-preserving AI primitives: authenticated symmetric encryption for payload confidentiality and a differential-privacy protected aggregation API. It uses Gin for HTTP, GORM for persistence (SQLite by default), and modular internal packages for crypto, DB, and privacy logic.

Architecture
-----------
- cmd/server: application entrypoint that wires components and starts the Gin HTTP server.
- internal/api: Gin handlers and request/response models.
- internal/crypto: AEAD encryption helpers (XChaCha20-Poly1305 / ChaCha20-Poly1305) to store confidential payloads.
- internal/db: GORM store initialization and migrations helper.
- internal/models: GORM data models (Query).
- internal/service: privacy logic and services (storing encrypted queries, decrypting and computing DP-aggregates).
- configs: sample YAML configuration.
- migrations: SQL migrations for DB schema.

High-level flow
---------------
1. A client submits a measurement via POST /v1/queries (example payload: {"score": 4.2}).
2. The server encrypts the measurement using a symmetric AEAD key and stores ciphertext + nonce.
3. A client requests an aggregate via GET /v1/aggregate?epsilon=0.5. The server decrypts stored payloads, computes an aggregate (average), applies Laplace noise according to epsilon and returns a DP-protected result.

Security & Privacy notes
------------------------
- Stored payloads are encrypted at rest using AEAD; the key is provided via environment variable ENCRYPTION_KEY (32 bytes base64 or raw). In production, keep keys in an HSM or KMS.
- Differential privacy implemented is a Laplace mechanism for numeric aggregation. This is a demonstrative implementation; production DP systems need careful design around sensitivity, composition, auditing, and parameter selection.

Getting started
---------------
Prerequisites:
- Go 1.20+
- make (optional)

Quick run (local SQLite):

1. Copy sample env:

   cp .env.example .env
   # optionally edit .env to set encryption key and DB path

2. Build and run:

   go build ./cmd/server
   ./server

   # or using `go run`
   go run ./cmd/server

3. Example API calls:

- Submit a query (POST)

  curl -X POST http://localhost:8080/v1/queries \
    -H "Content-Type: application/json" \
    -d '{"score": 4.5}'

- Get a DP aggregate (GET)

  curl "http://localhost:8080/v1/aggregate?epsilon=0.5"

Configuration
-------------
See configs/config.yaml for defaults. The server also reads the following environment variables (see .env.example):
- ENCRYPTION_KEY - 32-byte key used for AEAD (hex or base64 recommended).
- DB_PATH - path to SQLite DB file (default: data/privai.db)

Design choices
--------------
- SQLite is used for simplicity; swap with Postgres (gorm.io/driver/postgres) for production.
- AEAD encryption implemented with x/crypto chacha20poly1305 (ChaCha20-Poly1305) for authenticated encryption.
- DP: Laplace mechanism added to average. Sensitivity uses an assumed data range; tune this to your data.

Testing
-------
Run unit tests:

  go test ./...

Files of interest
-----------------
- cmd/server/main.go: app entrypoint
- internal/api/handlers.go: HTTP handlers
- internal/crypto/crypto.go: encryption/decryption
- internal/service/privacysvc.go: DP logic
- internal/db/store.go: GORM store and migrations

License
-------
This example repository is provided as-is for educational/demo purposes.

Contributions welcome — please open issues and PRs with improvements or real-world hardening suggestions.
