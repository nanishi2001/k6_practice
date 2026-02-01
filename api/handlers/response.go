package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse はエラーレスポンスの構造体
type ErrorResponse struct {
	Error string `json:"error"`
}

// writeJSON はJSONレスポンスを送信する
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError はエラーレスポンスを送信する
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

// methodNotAllowed は405エラーを返す
func methodNotAllowed(w http.ResponseWriter) {
	http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
}
