CREATE TABLE IF NOT EXISTS queries (
    id TEXT PRIMARY KEY,
    encrypted_payload BLOB NOT NULL,
    nonce BLOB NOT NULL,
    created_at DATETIME NOT NULL
);
