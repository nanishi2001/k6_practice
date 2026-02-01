package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type DelayResponse struct {
	DelayMs int       `json:"delay_ms"`
	Time    time.Time `json:"time"`
}

func DelayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/delay/")
	ms, err := strconv.Atoi(path)
	if err != nil || ms < 0 || ms > 10000 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid delay value (0-10000ms)"})
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
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	ms := rand.Intn(1001) // 0-1000ms
	time.Sleep(time.Duration(ms) * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DelayResponse{
		DelayMs: ms,
		Time:    time.Now(),
	})
}

type ErrorRateResponse struct {
	Success     bool      `json:"success"`
	ErrorRate   int       `json:"error_rate_percent"`
	RandomValue int       `json:"random_value"`
	Time        time.Time `json:"time"`
}

func ErrorRateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/error-rate/")
	percent, err := strconv.Atoi(path)
	if err != nil || percent < 0 || percent > 100 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid error rate (0-100)"})
		return
	}

	randomValue := rand.Intn(100)
	shouldError := randomValue < percent

	w.Header().Set("Content-Type", "application/json")

	if shouldError {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorRateResponse{
			Success:     false,
			ErrorRate:   percent,
			RandomValue: randomValue,
			Time:        time.Now(),
		})
		return
	}

	json.NewEncoder(w).Encode(ErrorRateResponse{
		Success:     true,
		ErrorRate:   percent,
		RandomValue: randomValue,
		Time:        time.Now(),
	})
}
