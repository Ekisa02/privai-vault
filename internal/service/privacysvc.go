package service

import (
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/example/privai-vault/internal/crypto"
	"github.com/example/privai-vault/internal/db"
	"github.com/example/privai-vault/internal/models"
)

// PrivacyService coordinates storing encrypted measurements and computing DP aggregates.
type PrivacyService struct {
	store *db.GormStore
	crypto *crypto.AEAD
	// assumed data range for sensitivity calculation
	minScore float64
	maxScore float64
}

func NewPrivacyService(store *db.GormStore, c *crypto.AEAD) *PrivacyService {
	rand.Seed(time.Now().UnixNano())
	return &PrivacyService{store: store, crypto: c, minScore: 0.0, maxScore: 10.0}
}

// StoreMeasurement encrypts a single numeric measurement and persists it.
func (p *PrivacyService) StoreMeasurement(score float64) (string, error) {
	// create JSON payload
	payload, _ := json.Marshal(map[string]float64{"score": score})
	ciphertext, nonce, err := p.crypto.Encrypt(payload)
	if err != nil {
		return "", err
	}
	q := models.NewQuery(ciphertext, nonce)
	if err := p.store.SaveQuery(q); err != nil {
		return "", err
	}
	return q.ID, nil
}

// GetDPAverage decrypts all stored measurements, computes an average, and returns a DP-noised average.
func (p *PrivacyService) GetDPAverage(epsilon float64) (float64, int, error) {
	rows, err := p.store.ListQueries()
	if err != nil {
		return 0, 0, err
	}
	if len(rows) == 0 {
		return 0, 0, errors.New("no data")
	}
	var sum float64
	count := 0
	for _, r := range rows {
		pt, err := p.crypto.Decrypt(r.EncryptedPayload, r.Nonce)
		if err != nil {
			// skip malformed rows rather than fail entire batch
			continue
		}
		var m map[string]float64
		if err := json.Unmarshal(pt, &m); err != nil {
			continue
		}
		val, ok := m["score"]
		if !ok {
			continue
		}
		sum += val
		count++
	}
	if count == 0 {
		return 0, 0, errors.New("no decryptable data")
	}
	avg := sum / float64(count)

	// For average, l1 sensitivity = (max-min) / n
	sensitivity := (p.maxScore - p.minScore) / float64(count)
	scale := sensitivity / epsilon
	noise := laplaceNoise(scale)
	noisy := avg + noise
	return noisy, count, nil
}

// laplaceNoise samples from Laplace(0, b) using inverse transform.
func laplaceNoise(b float64) float64 {
	if b <= 0 {
		return 0
	}
	u := rand.Float64() - 0.5
	return -b * math.Signbit(u) * math.Log(1-2*math.Abs(u))
}
