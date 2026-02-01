package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 遅延の制限値
const (
	maxDelayMs       = 10000
	maxRandomDelayMs = 1001
	maxErrorRate     = 100
)

type DelayResponse struct {
	DelayMs int       `json:"delay_ms"`
	Time    time.Time `json:"time"`
}

type ErrorRateResponse struct {
	Success     bool      `json:"success"`
	ErrorRate   int       `json:"error_rate_percent"`
	RandomValue int       `json:"random_value"`
	Time        time.Time `json:"time"`
}

func DelayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/delay/")
	ms, err := strconv.Atoi(path)
	if err != nil || ms < 0 || ms > maxDelayMs {
		writeError(w, http.StatusBadRequest, "invalid delay value (0-10000ms)")
		return
	}

	time.Sleep(time.Duration(ms) * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DelayResponse{
		DelayMs: ms,
		Time:    time.Now(),
	})
}

func RandomDelayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	ms := rand.Intn(maxRandomDelayMs)
	time.Sleep(time.Duration(ms) * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DelayResponse{
		DelayMs: ms,
		Time:    time.Now(),
	})
}

func ErrorRateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/error-rate/")
	percent, err := strconv.Atoi(path)
	if err != nil || percent < 0 || percent > maxErrorRate {
		writeError(w, http.StatusBadRequest, "invalid error rate (0-100)")
		return
	}

	randomValue := rand.Intn(100)
	shouldError := randomValue < percent

	response := ErrorRateResponse{
		Success:     !shouldError,
		ErrorRate:   percent,
		RandomValue: randomValue,
		Time:        time.Now(),
	}

	if shouldError {
		writeJSON(w, http.StatusInternalServerError, response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
