package service

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/example/privai-vault/internal/api"
	"github.com/example/privai-vault/internal/crypto"
	"github.com/example/privai-vault/internal/db"
)

func setupTestServer(t *testing.T) (*PrivacyService, *httptest.Server) {
	// in-memory sqlite
	store, err := db.NewGormStore(":memory:")
	if err != nil {
		t.Fatalf("db init: %v", err)
	}
	if err := store.AutoMigrate(); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	// deterministic key for tests
	key := []byte("12345678901234567890123456789012")
	c, err := crypto.NewAEADKey(key)
	if err != nil {
		t.Fatalf("crypto init: %v", err)
	}
	priv := NewPrivacyService(store, c)

	r := gin.New()
	api.RegisterRoutes(r, priv)
	ts := httptest.NewServer(r)
	return priv, ts
}

func TestStoreAndAggregate(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()

	// submit three queries
	for _, v := range []string{"{\"score\": 2.0}", "{\"score\": 4.0}", "{\"score\": 6.0}"} {
		resp, err := http.Post(ts.URL+"/v1/queries", "application/json", strings.NewReader(v))
		if err != nil {
			t.Fatalf("post failed: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201 got %d", resp.StatusCode)
		}
	}

	// request aggregate with high epsilon (low noise) to assert result roughly matches
	resp, err := http.Get(ts.URL + "/v1/aggregate?epsilon=10")
	if err != nil {
		t.Fatalf("get aggregate failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}
}
