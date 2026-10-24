package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/example/privai-vault/internal/api"
	"github.com/example/privai-vault/internal/crypto"
	"github.com/example/privai-vault/internal/db"
	"github.com/example/privai-vault/internal/service"
)

func Run() {
	// read config from env
	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		log.Println("WARNING: ENCRYPTION_KEY not set; using default insecure key. Set ENCRYPTION_KEY in production!")
		encryptionKey = "changeme-32-byte-key-0123456789ab"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/privai.db"
	}

	// initialize store
	store, err := db.NewGormStore(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}

	// ensure migrations exist (creates table if not present)
	if err := store.AutoMigrate(); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// setup crypto
	key := []byte(encryptionKey)
	cryptoSvc, err := crypto.NewAEADKey(key)
	if err != nil {
		log.Fatalf("failed to initialize crypto: %v", err)
	}

	// setup privacy service
	privSvc := service.NewPrivacyService(store, cryptoSvc)

	// setup router
	r := gin.Default()

	api.RegisterRoutes(r, privSvc)

	addr := ":8080"
	log.Printf("starting server on %s (db: %s, api: /v1)", addr, dbPath)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
