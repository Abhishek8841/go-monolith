package httpx

import (
	"encoding/json"
	"net/http"
)

type Code string

const (
	CodeInvalidID     Code = "invalid_id"
	CodeInternalError Code = "internal_error"
)

type errorEnvelope struct {
	Error errorPayLoad `json:"error"`
}

type errorPayLoad struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

func Error(w http.ResponseWriter, status int, message string, code Code) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(errorEnvelope{
		Error: errorPayLoad{
			Code:    code,
			Message: message,
		},
	})
}

// we shouldn't depend on codes to tell what happened to the frontend as codes can change
// Message field is also fragile as it can change depending on factors like lang etc
// therefore, code field is used to tell the frontend exactly what happened
