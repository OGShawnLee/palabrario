package response

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Error string `json:"error"`
}

// JSON escribe una respuesta serializada en application/json con el código de estado indicado.
func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

// Error escribe un error estandarizado {"error": "..."} con Content-Type application/json.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorBody{Error: message})
}
