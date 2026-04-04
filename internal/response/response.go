package response

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func WriteError(w http.ResponseWriter, status int, message string) {
	resp := ErrorResponse{
		Error: ErrorBody{
			Status:  status,
			Message: message,
		},
	}
	WriteJSON(w, status, resp)
}
