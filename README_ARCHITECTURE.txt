Architecture (textual)

[client] -> HTTP -> [Gin HTTP Server]
                           |
                           v
                    [PrivacyService]
                     /         \
                    v           v
           [internal/crypto]   [internal/db (GORM + SQLite)]
                           |
                           v
                store encrypted payloads in DB

PrivacyService retrieves payloads, decrypts with AEAD, computes aggregate (average) and then applies Laplace noise using provided epsilon.
