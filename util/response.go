package util

import (
	"encoding/json"
	"net/http"
)

// RespondJSON writes a JSON response with the given status code.
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log := GetLogger()
			log.Error("Failed to encode JSON response: %v", err, true)
		}
	}
}

// RespondError writes a JSON error response with the given status code and message.
func RespondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		log := GetLogger()
		log.Error("Failed to encode error response: %v", err, true)
	}
}
