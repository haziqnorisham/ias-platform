package http

import (
	"encoding/json"
	ias_pg "ias/automation/db/pg"
	"net/http"
	"os"
	"strconv"
	"time"

	redis_lib "github.com/redis/go-redis/v9"
)

func getAllTreeSensorHandler(w http.ResponseWriter, r *http.Request, rdb *redis_lib.Client) {
	// Set content type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Check Redis Cache First
	cacheKey := "all_tree_sensors"
	cachedData, err := rdb.Get(r.Context(), cacheKey).Result()
	if err == nil {
		// Cache hit, return cached data
		w.Write([]byte(cachedData))
		return
	}

	// Get sensors data if cache miss
	sensors, err := ias_pg.NewPostgresStorage(nil).QueryData("select * from ppj_tree_sensor")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Encode sensors to JSON for both caching and response
	jsonData, err := json.Marshal(sensors)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to encode response"})
		return
	}

	// Write to cache with expiration (e.g., 5 minutes)
	cache_invalidation_seconds, err := strconv.Atoi(os.Getenv("REDIS_INVALIDATION_SECOND"))
	if err != nil {
		cache_invalidation_seconds = 60 // Default to 5 minutes if env variable is not set or invalid
	}
	cacheExpiration := time.Duration(cache_invalidation_seconds) * time.Second
	err = rdb.Set(r.Context(), cacheKey, jsonData, cacheExpiration).Err()
	if err != nil {
		// Log cache write error but don't fail the request
		// You might want to log this: log.Printf("Failed to write to cache: %v", err)
	}

	// Send the response
	w.Write(jsonData)
}
