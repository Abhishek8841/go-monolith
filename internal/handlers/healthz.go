package handlers

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	// here order matters ie write header flushes the header and we cant add new headers after that for .header.set method which just registers the headers should be written first

	w.Write([]byte(`{"status":"ok"}`))
}
