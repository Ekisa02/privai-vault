package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/example/privai-vault/internal/service"
)

// request/response shapes
type SubmitQueryRequest struct {
	Score float64 `json:"score" binding:"required"`
}

type SubmitQueryResponse struct {
	ID string `json:"id"`
}

type AggregateResponse struct {
	NoisyAverage float64 `json:"noisy_average"`
	Count        int     `json:"count"`
	Epsilon      float64 `json:"epsilon"`
}

// RegisterRoutes wires handlers to router
func RegisterRoutes(r *gin.Engine, svc *service.PrivacyService) {
	v1 := r.Group("/v1")
	{
		v1.POST("/queries", makeSubmitHandler(svc))
		v1.GET("/aggregate", makeAggregateHandler(svc))
	}
}

func makeSubmitHandler(svc *service.PrivacyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SubmitQueryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		id, err := svc.StoreMeasurement(req.Score)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store"})
			return
		}
		c.JSON(http.StatusCreated, SubmitQueryResponse{ID: id})
	}
}

func makeAggregateHandler(svc *service.PrivacyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		epsStr := c.Query("epsilon")
		if epsStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "epsilon query param required"})
			return
		}
		epsilon, err := strconv.ParseFloat(epsStr, 64)
		if err != nil || epsilon <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid epsilon"})
			return
		}

		noisy, count, err := svc.GetDPAverage(epsilon)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "compute failed"})
			return
		}
		c.JSON(http.StatusOK, AggregateResponse{NoisyAverage: noisy, Count: count, Epsilon: epsilon})
	}
}

// Helper: JSON encode for debugging
func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
